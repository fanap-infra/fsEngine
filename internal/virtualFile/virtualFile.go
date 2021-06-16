package virtualFile

import (
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"

	"github.com/fanap-infra/log"
)

// VirtualFile
type VirtualFile struct {
	bufRX []byte
	bufTX []byte

	bufferSize         int
	currentBlock       uint32
	currentBlockIndex  uint32 // current block's index
	name               string
	id                 uint32
	firstBlock         uint32
	Closed             bool
	lastBlock          uint32
	blockAllocationMap *blockAllocationMap.BlockAllocationMap
	blockSize          uint32
	nextBlockIndex     uint32
	size               uint64

	allocatedBlock []uint32
	readOnly       bool

	seekPointer int
	bufStart    int
	bufEnd      int

	fs  FS
	log *log.Logger
}
