package virtualFile

import (
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"

	"github.com/fanap-infra/log"
)

func OpenVirtualFile(name string, fileID uint32, numberOfBlocks uint32, blockSize uint32, fs FS, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		vfBuf:          nil,
		name:           name,
		id:             fileID,
		Closed:         false,
		readOnly:       true,
		numberOfBlocks: numberOfBlocks,
		allocatedBlock: make([]uint32, 0),
		blockSize:      blockSize,
		size:           uint64(numberOfBlocks * blockSize),
		fs:             fs,
		log:            log,
	}
}

func NewVirtualFile(fileName string, fileID uint32, fs FS, eHandler blockAllocationMap.Events, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		vfBuf:              make([]byte, 0),
		name:               fileName,
		id:                 fileID,
		Closed:             false,
		blockAllocationMap: blockAllocationMap.New(log, eHandler, 0),
		allocatedBlock:     make([]uint32, 0),
		numberOfBlocks:     0,
		fs:                 fs,
		log:                log,
	}
}
