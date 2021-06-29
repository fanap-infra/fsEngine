package fsEngine

import (
	"fmt"

	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/FSEngine/internal/virtualFile"
)

// Create new virtual file and add opened files
func (fse *FSEngine) NewVirtualFile(id uint32, fileName string) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	if fse.header.CheckIDExist(id) {
		return nil, fmt.Errorf("this ID: %v, had been taken", id)
	}
	blm := blockAllocationMap.New(fse.log, fse, fse.maxNumberOfBlocks)

	vf := virtualFile.NewVirtualFile(fileName, id, fse.blockSize-BlockHeaderSize, fse, blm,
		int(fse.blockSize-BlockHeaderSize)*VirtualFileBufferBlockNumber, fse.log)
	err := fse.header.AddVirtualFile(id, fileName)
	if err != nil {
		return nil, err
	}
	fse.openFiles[id] = vf
	return vf, nil
}

func (fse *FSEngine) OpenVirtualFile(id uint32) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	_, ok := fse.openFiles[id]
	if ok {
		return nil, fmt.Errorf("this ID: %v is opened before", id)
	}
	fileInfo, err := fse.header.GetFileData(id)
	if err != nil {
		return nil, err
	}
	blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		fileInfo.GetRMapBlocks())
	if err != nil {
		return nil, err
	}
	vf := virtualFile.OpenVirtualFile(&fileInfo, fse.blockSize-BlockHeaderSize, fse, blm,
		int(fse.blockSize-BlockHeaderSize)*VirtualFileBufferBlockNumber, fse.log)
	//err = fse.header.AddVirtualFile(id, fileInfo.GetName())
	//if err != nil {
	//	return nil, err
	//}
	return vf, nil
}

func (fse *FSEngine) RemoveVirtualFile(id uint32) error {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	_, ok := fse.openFiles[id]
	if ok {
		return fmt.Errorf("virtual file id : %d is opened", id)
	}
	return nil
}
