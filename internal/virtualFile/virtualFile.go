package virtualFile

import (
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"

	"github.com/fanap-infra/log"
)

// VirtualFile
type VirtualFile struct {
	vfBuf []byte
	// vfBufMux sync.Mutex
	// sem                *semaphore.Weighted
	// maybeSync         *semaphore.Weighted
	currentBlock      uint32
	currentBlockIndex uint32 // current block's index
	name              string
	id                uint32
	// lastReadBlockIdx   uint32
	firstBlock         uint32
	Closed             bool
	lastBlock          uint32
	blockAllocationMap *blockAllocationMap.BlockAllocationMap
	blockSize          uint32
	size               uint64
	numberOfBlocks     uint32
	allocatedBlock     []uint32
	readOnly           bool
	// Media Structures
	// frameChunkPosInBlock uint32
	// fwMUX                sync.Mutex
	// reference to parent Archiver
	fs  FS
	log *log.Logger
}
