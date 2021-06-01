package Header_

import (
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/FSEngine/internal/fileIndex"
	"os"
	"time"

	"github.com/fanap-infra/log"
)

type HFileSystem struct {
	file               *os.File // file handle instance
	version            uint32
	size               int64
	CurrentFile        string                                 // name of the latest file to be created
	LastFiletime       time.Time                              // time where the first data of the file has been written
	maxNumberOfBlocks             uint32                                 // total number of blocks in Archiver
	blockSize          uint32                                 // in bytes, size of each block
	lastWrittenBlock   uint32                                 // the last block that has been written into
	blockAllocationMap *blockAllocationMap.BlockAllocationMap // BAM data in memory coded with roaring, to be synced later on to Disk.
	//openFiles          map[uint32]*virtualFile.VirtualFile
	fileIndex          *fileIndex.FileIndex
	fileIndexSize      uint32
	blmSize      uint32
	// WMux               sync.Mutex
	// RMux               sync.Mutex
	log *log.Logger
	// fiMux              sync.RWMutex
	fiChecksum uint32
	// bamChecksum uint32
	// fsMux              sync.Mutex
	// rIBlockMux         sync.Mutex
	// Cache              *lru.Cache
	// fileIndexIsFlip    bool
	conf configs
}

func (hfs *HFileSystem) UpdateFSHeader() error {
	err := hfs.updateFileIndex()
	if err != nil {
		return err
	}

	err = hfs.updateHeader()
	if err != nil {
		return err
	}
	return nil
}
