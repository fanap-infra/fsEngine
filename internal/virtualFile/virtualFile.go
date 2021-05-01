package virtualFile

import (
	"behnama/stream/internal/messages/media"
	"sync"

	"github.com/fanap-infra/log"

	"github.com/RoaringBitmap/roaring"
	"golang.org/x/sync/semaphore"
)

// VirtualFile
type VirtualFile struct {
	vfBuf             []byte
	vfBufMux          sync.Mutex
	sem               *semaphore.Weighted
	maybeSync         *semaphore.Weighted
	currentBlock      uint32
	currentBlockIndex uint32 // current block's index
	name              string
	id                uint32
	lastReadBlockIdx  uint32
	firstBlock        uint32
	Closed            bool
	lastBlock         uint32
	blocks            *roaring.Bitmap
	blockSize         uint32
	size              uint64
	numberOfBlocks    uint32
	allocatedBlock    []uint32
	readOnly          bool
	// Media Structures
	frameChunk           *media.PacketChunk
	curFrameInChunk      uint32
	frameChunkPosInBlock uint32
	fwMUX                sync.Mutex
	// reference to parent Archiver
	fs  FS
	log *log.Logger
}
