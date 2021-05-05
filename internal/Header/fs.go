package Header_

import (
	"behnama/stream/pkg/archiverStorageEngine/internals/virtualFS"
	"behnama/stream/pkg/fsEngine/internal/blockAllocationMap"
	"behnama/stream/pkg/fsEngine/internal/fileIndex"
	"os"
	"time"

	"github.com/fanap-infra/log"
)

type FileSystem struct {
	file               *os.File // file handle instance
	version            uint32
	size               int64
	CurrentFile        string                                 // name of the latest file to be created
	LastFiletime       time.Time                              // time where the first data of the file has been written
	blocks             uint32                                 // total number of blocks in Archiver
	blockSize          uint32                                 // in bytes, size of each block
	lastWrittenBlock   uint32                                 // the last block that has been written into
	blockAllocationMap *blockAllocationMap.BlockAllocationMap // BAM data in memory coded with roaring, to be synced later on to Disk.
	openFiles          map[uint32]*virtualFS.VirtualFile
	fileIndex          *fileIndex.FileIndex
	fileIndexSize      uint32
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

func (fs *FileSystem) UpdateFSHeader() error {
	err := fs.updateFileIndex()
	if err != nil {
		return err
	}

	err = fs.updateHeader()
	if err != nil {
		return err
	}
	return nil
}
