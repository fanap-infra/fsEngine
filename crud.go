package fsEngine

import (
	"fmt"

	"github.com/fanap-infra/log"

	"github.com/fanap-infra/fsEngine/internal/constants"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/virtualFile"
)

// Create new virtual file and add opened files
func (fse *FSEngine) NewVirtualFile(id uint32, fileName string) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	if fse.header.CheckIDExist(id) {
		return nil, fmt.Errorf("this ID: %v, had been taken", id)
	}
	blm := blockAllocationMap.New(fse.log, fse, fse.maxNumberOfBlocks)

	vf := virtualFile.NewVirtualFile(fileName, id, fse.blockSize-constants.BlockHeaderSize, fse, blm,
		int(fse.blockSize-constants.BlockHeaderSize)*constants.VirtualFileBufferBlockNumber, fse.log)
	err := fse.header.AddVirtualFile(id, fileName)
	if err != nil {
		return nil, err
	}

	vfInfo := &VFInfo{id: id, blm: blm, numberOfOpened: 1}
	vfInfo.vfs = append(vfInfo.vfs, vf)
	fse.openFiles[id] = vfInfo

	err = fse.header.UpdateFSHeader()
	if err != nil {
		fse.log.Warnv("Can not updateHeader", "err", err.Error())
		return nil, err
	}
	return vf, nil
}

func (fse *FSEngine) OpenVirtualFile(id uint32) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	vfInfo, ok := fse.openFiles[id]
	if ok {
		fileInfo, err := fse.header.GetFileData(id)
		if err != nil {
			return nil, err
		}
		//if len(fileInfo.GetRMapBlocks()) == 0 {
		//	log.Warnv("can not open virtual file, roaring byte array length is zero", "id", id)
		//	return nil, fmt.Errorf("virtual file is empty, id: %v", id)
		//}
		//blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		//	fileInfo.GetRMapBlocks())
		//if err != nil {
		//	return nil, err
		//}
		vf := virtualFile.OpenVirtualFile(&fileInfo, fse.blockSize-constants.BlockHeaderSize, fse, vfInfo.blm,
			int(fse.blockSize-constants.BlockHeaderSize)*constants.VirtualFileBufferBlockNumber, fse.log)
		vfInfo.numberOfOpened = vfInfo.numberOfOpened + 1
		return vf, nil
	}

	fileInfo, err := fse.header.GetFileData(id)
	if err != nil {
		return nil, err
	}
	if len(fileInfo.GetRMapBlocks()) == 0 {
		log.Warnv("can not open virtual file, roaring byte array length is zero", "id", id)
		return nil, fmt.Errorf("virtual file is empty, id: %v", id)
	}
	blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		fileInfo.GetRMapBlocks())
	if err != nil {
		return nil, err
	}
	vf := virtualFile.OpenVirtualFile(&fileInfo, fse.blockSize-constants.BlockHeaderSize, fse, blm,
		int(fse.blockSize-constants.BlockHeaderSize)*constants.VirtualFileBufferBlockNumber, fse.log)

	vfInfo = &VFInfo{id: id, blm: blm, numberOfOpened: 1}
	vfInfo.vfs = append(vfInfo.vfs, vf)
	fse.openFiles[id] = vfInfo

	return vf, nil
}

func (fse *FSEngine) OpenVirtualFileForRecovery(id uint32) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	vfInfo, ok := fse.openFiles[id]
	if ok {
		fileInfo, err := fse.header.GetFileData(id)
		if err != nil {
			return nil, err
		}

		vf := virtualFile.OpenVirtualFileForRecovery(&fileInfo, fse.blockSize-constants.BlockHeaderSize, fse, vfInfo.blm,
			int(fse.blockSize-constants.BlockHeaderSize)*constants.VirtualFileBufferBlockNumber, fse.log)
		vfInfo.numberOfOpened = vfInfo.numberOfOpened + 1
		return vf, nil
	}

	fileInfo, err := fse.header.GetFileData(id)
	if err != nil {
		return nil, err
	}
	if len(fileInfo.GetRMapBlocks()) == 0 {
		log.Warnv("can not open virtual file, roaring byte array length is zero", "id", id)
		return nil, fmt.Errorf("virtual file is empty, id: %v", id)
	}
	blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		fileInfo.GetRMapBlocks())
	if err != nil {
		return nil, err
	}
	vf := virtualFile.OpenVirtualFileForRecovery(&fileInfo, fse.blockSize-constants.BlockHeaderSize, fse, blm,
		int(fse.blockSize-constants.BlockHeaderSize)*constants.VirtualFileBufferBlockNumber, fse.log)

	vfInfo = &VFInfo{id: id, blm: blm, numberOfOpened: 1}
	vfInfo.vfs = append(vfInfo.vfs, vf)
	fse.openFiles[id] = vfInfo

	return vf, nil
}

func (fse *FSEngine) RemoveVirtualFile(id uint32) (int, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	_, ok := fse.openFiles[id]
	if ok {
		return 0, fmt.Errorf("virtual file id : %d is opened", id)
	}

	fileInfo, err := fse.header.GetFileData(id)
	if err != nil {
		log.Warnv("can not remove virtual, because can not get file data", "id", id)
		return 0, err
	}
	if len(fileInfo.GetRMapBlocks()) == 0 {
		log.Warnv("can not remove virtual files block, roaring byte array length is zero", "id", id)
		return 0, fse.header.RemoveVirtualFile(id)
	}
	blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		fileInfo.GetRMapBlocks())
	if err != nil {
		log.Errorv("can not parse block allocation map", "id", id,
			"len(fileInfo.GetRMapBlocks()) ", len(fileInfo.GetRMapBlocks()), "err", err.Error())
		return 0, fse.header.RemoveVirtualFile(id)
	}

	blocks := blm.ToArray()
	log.Infov("remove virtual file", "len(blocks)", len(blocks), "id", id, "blocks", blocks)
	for _, bIndex := range blocks {
		fse.header.UnsetBlockAsAllocated(bIndex)
	}
	// fse.log.Infov("blm length",
	//	"fse blocks length", len(fse.blockAllocationMap.ToArray()))
	return len(blocks), fse.header.RemoveVirtualFile(id)
}
