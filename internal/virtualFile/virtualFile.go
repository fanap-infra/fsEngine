package virtualFile

import (
	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"

	"github.com/fanap-infra/log"
)

// VirtualFile
type VirtualFile struct {
	bufRX []byte
	bufTX []byte

	bufferSize         int
	name               string
	id                 uint32
	Closed             bool
	lastBlock          uint32
	firstBlockIndex    uint32
	blockAllocationMap *blockAllocationMap.BlockAllocationMap
	blockSize          uint32
	nextBlockIndex     uint32

	allocatedBlock []uint32
	readOnly       bool

	seekPointer int
	bufStart    int
	bufEnd      int

	fs  FS
	log *log.Logger
}
