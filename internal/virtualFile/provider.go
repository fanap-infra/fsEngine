package virtualFile

import (
	"behnama/stream/internal/messages/media"

	"github.com/fanap-infra/log"
)

func OpenVirtualFile(path string, fileID uint32, numberOfBlocks uint32, blockSize uint32, fs FS, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		vfBuf:          nil,
		name:           path,
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

func NewVirtualFile(fileName string, fileID uint32, fs FS, log *log.Logger) *VirtualFile {
	return &VirtualFile{
		vfBuf:  make([]byte, 0),
		name:   fileName,
		id:     fileID,
		Closed: false,
		//blocks:         roaring.NewBitmap(),
		allocatedBlock: make([]uint32, 0),
		frameChunk:     &media.PacketChunk{},
		numberOfBlocks: 0,
		//sem:            semaphore.NewWeighted(1),
		//maybeSync:      semaphore.NewWeighted(1),
		fs:  fs,
		log: log,
	}
}
