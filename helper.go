package fsEngine

import (
	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/fileIndex"
)

func (fse *FSEngine) GetFileList() []*fileIndex.File {
	return fse.header.GetFilesList()
}

func (fse *FSEngine) GetFileBLM(id uint32) (*blockAllocationMap.BlockAllocationMap, error) {
	fileInfo, err := fse.header.GetFileData(id)
	if err != nil {
		return nil, err
	}
	blm, err := blockAllocationMap.Open(fse.log, fse, fse.maxNumberOfBlocks, fileInfo.GetLastBlock(),
		fileInfo.GetRMapBlocks())
	if err != nil {
		return nil, err
	}
	return blm, nil
}
