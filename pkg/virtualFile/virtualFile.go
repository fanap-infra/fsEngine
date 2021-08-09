package virtualFile

import (
	"sync"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"

	"github.com/fanap-infra/log"
)

// VirtualFile
type VirtualFile struct {
	bufRX              []byte
	bufTX              []byte
	optionalData       []byte
	bufferSize         int
	name               string
	id                 uint32
	Closed             bool
	lastBlock          uint32
	firstBlockIndex    uint32
	blockAllocationMap *blockAllocationMap.BlockAllocationMap
	blockSize          uint32
	nextBlockIndex     uint32
	WMux               sync.Mutex
	allocatedBlock     []uint32
	readOnly           bool
	fileSize           uint32
	seekPointer        int
	bufStart           int
	bufEnd             int

	fs  FS
	log *log.Logger
}

func (v *VirtualFile) GetFileName() string {
	return v.name
}

func (v *VirtualFile) GetSeek() int {
	return v.seekPointer
}

func (v *VirtualFile) GetFileSize() uint32 {
	return v.fileSize
}

func (v *VirtualFile) GetOptionalData() []byte {
	return v.optionalData
}
