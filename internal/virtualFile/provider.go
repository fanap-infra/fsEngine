package virtualFile

import (
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/FSEngine/internal/fileIndex"

	"github.com/fanap-infra/log"
)
// name string, fileID uint32,
func OpenVirtualFile( fileInfo fileIndex.File, blockSize uint32,
	 fs FS, blm *blockAllocationMap.BlockAllocationMap, log *log.Logger) *VirtualFile {

	return &VirtualFile{
		vfBuf:              make([]byte, 0),
		name:           fileInfo.GetName(),
		id:             fileInfo.GetId(),
		Closed:         false,
		readOnly:       true,
		//numberOfBlocks: numberOfBlocks,
		blockAllocationMap: blm,
		allocatedBlock: make([]uint32, 0),
		blockSize:      blockSize,
		//size:           uint64(numberOfBlocks * blockSize),
		fs:             fs,
		log:            log,
	}
}

func NewVirtualFile(fileName string, fileID uint32, blockSize uint32, fs FS, blm *blockAllocationMap.BlockAllocationMap, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		vfBuf:              make([]byte, 0),
		name:               fileName,
		id:                 fileID,
		blockSize: blockSize,
		Closed:             false,
		blockAllocationMap: blm,
		allocatedBlock:     make([]uint32, 0),
		//numberOfBlocks:     0,
		fs:                 fs,
		log:                log,
	}
}
