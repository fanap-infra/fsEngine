package virtualFile

import (
	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/fileIndex"

	"github.com/fanap-infra/log"
)

// name string, fileID uint32,
func OpenVirtualFile(fileInfo *fileIndex.File, blockSize uint32,
	fs FS, blm *blockAllocationMap.BlockAllocationMap, bufferSize int, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		bufRX:        make([]byte, 0),
		bufTX:        make([]byte, 0),
		optionalData: fileInfo.Optional,
		name:         fileInfo.GetName(),
		id:           fileInfo.GetId(),
		Closed:       false,
		readOnly:     true,
		// numberOfBlocks: numberOfBlocks,
		blockAllocationMap: blm,
		allocatedBlock:     make([]uint32, 0),
		blockSize:          blockSize,
		//size:           uint64(numberOfBlocks * blockSize),
		fs:              fs,
		log:             log,
		seekPointer:     0,
		bufStart:        0,
		bufEnd:          0,
		nextBlockIndex:  0,
		firstBlockIndex: fileInfo.FirstBlock,
		bufferSize:      bufferSize,
		fileSize:        fileInfo.GetFileSize(),
	}
}

func OpenVirtualFileForRecovery(fileInfo *fileIndex.File, blockSize uint32,
	fs FS, blm *blockAllocationMap.BlockAllocationMap, bufferSize int, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		bufRX:        make([]byte, 0),
		bufTX:        make([]byte, 0),
		optionalData: fileInfo.Optional,
		name:         fileInfo.GetName(),
		id:           fileInfo.GetId(),
		Closed:       false,
		readOnly:     false,
		// numberOfBlocks: numberOfBlocks,
		blockAllocationMap: blm,
		allocatedBlock:     make([]uint32, 0),
		blockSize:          blockSize,
		//size:           uint64(numberOfBlocks * blockSize),
		fs:              fs,
		log:             log,
		seekPointer:     0,
		bufStart:        0,
		bufEnd:          0,
		nextBlockIndex:  0,
		firstBlockIndex: fileInfo.FirstBlock,
		bufferSize:      bufferSize,
		fileSize:        fileInfo.GetFileSize(),
	}
}

func NewVirtualFile(fileName string, fileID uint32, blockSize uint32, fs FS,
	blm *blockAllocationMap.BlockAllocationMap, bufferSize int, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		bufRX:              make([]byte, 0),
		bufTX:              make([]byte, 0),
		optionalData:       make([]byte, 0),
		name:               fileName,
		id:                 fileID,
		blockSize:          blockSize,
		Closed:             false,
		blockAllocationMap: blm,
		allocatedBlock:     make([]uint32, 0),
		fs:                 fs,
		log:                log,
		seekPointer:        0,
		bufStart:           0,
		bufEnd:             0,
		nextBlockIndex:     0,
		firstBlockIndex:    0,
		bufferSize:         bufferSize,
		fileSize:           0,
	}
}
