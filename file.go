package fsEngine

import (
	"github.com/fanap-infra/FSEngine/internal/fileIndex"
	Header_ "github.com/fanap-infra/FSEngine/internal/Header"
	"github.com/fanap-infra/FSEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/FSEngine/internal/virtualFile"
	"os"
	"sync"
	"time"

	"github.com/fanap-infra/log"

	lru "github.com/hashicorp/golang-lru"
)

// File
type FSEngine struct {
	file               *os.File // file handle instance
	header             *Header_.HFileSystem
	version            uint32
	size               int64
	CurrentFile        string                                 // name of the latest file to be created
	LastFiletime       time.Time                              // time where the first data of the file has been written
	maxNumberOfBlocks             uint32                                 // total number of blocks in Archiver
	blockSize          uint32                                 // in bytes, size of each block
	lastWrittenBlock   uint32                                 // the last block that has been written into
	blockAllocationMap *blockAllocationMap.BlockAllocationMap // BAM data in memory coded with roaring, to be synced later on to Disk.
	openFiles          map[uint32]*virtualFile.VirtualFile
	fileIndex          fileIndex.FileIndex
	WMux               sync.Mutex
	RMux               sync.Mutex
	log                *log.Logger
	fiMux              sync.RWMutex
	fiChecksum         uint32
	bamChecksum        uint32
	fsMux              sync.Mutex
	rIBlockMux         sync.Mutex
	crudMutex          sync.Mutex
	Cache              *lru.Cache
	fileIndexIsFlip    bool
	EventsHandler      Events
	Quit               chan struct{}
}

// Close ...
func (fse *FSEngine) Close() error {
	err := fse.header.UpdateFSHeader()
	if err != nil {
		fse.log.Warnv("Can not updateHeader", "err", err.Error())
		// ToDo: remove it
		return err
	}
	// ToDo:update file system
	err = fse.file.Sync()
	if err != nil {
		fse.log.Warnv("Can not sync file", "err", err.Error())
		// ToDo: remove it
		return err
	}
	return fse.file.Close()
}
