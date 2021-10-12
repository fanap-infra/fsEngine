package fsEngine

import (
	"os"
	"sync"
	"time"

	Header_ "github.com/fanap-infra/fsEngine/pkg/Header"

	"github.com/fanap-infra/log"

	lru "github.com/hashicorp/golang-lru"
)

// File
type FSEngine struct {
	id                uint32
	file              *os.File // file handle instance
	header            *Header_.HFileSystem
	version           uint32
	size              int64
	CurrentFile       string    // name of the latest file to be created
	LastFiletime      time.Time // time where the first data of the file has been written
	maxNumberOfBlocks uint32    // total number of blocks in Archiver
	blockSize         uint32    // in bytes, size of each block
	blockSizeUsable   uint32
	openFiles         map[uint32]*VFInfo
	WMux              sync.Mutex
	RMux              sync.Mutex
	log               *log.Logger
	crudMutex         sync.Mutex
	Cache             *lru.Cache
	eventsHandler     Events
	cleaning          uint32
	Quit              chan struct{}
}

// Close ...
func (fse *FSEngine) Close() error {
	for _, vfInfo := range fse.openFiles {
		for _, vf := range vfInfo.vfs {
			err := vf.Close()
			if err != nil {
				fse.log.Warnv("Can not close virtual file", "err", err.Error())
				return err
			}
		}
	}

	err := fse.header.UpdateFSHeader()
	if err != nil {
		fse.log.Warnv("Can not updateHeader", "err", err.Error())
		return err
	}
	// ToDo:update file system
	err = fse.file.Sync()
	if err != nil {
		fse.log.Warnv("Can not sync file", "err", err.Error())
		return err
	}
	return fse.file.Close()
}

func (fse *FSEngine) GetFilePath() string {
	return fse.file.Name()
}

func (fse *FSEngine) GetBlockSize() uint32 {
	return fse.header.GetBlockSize()
}
