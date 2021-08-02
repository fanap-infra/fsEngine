package Header_

import (
	"os"
	"time"

	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/internal/fileIndex"

	"github.com/fanap-infra/log"
)

type HFileSystem struct {
	file               *os.File // file handle instance
	version            uint32
	size               int64
	CurrentFile        string                                 // name of the latest file to be created
	LastFiletime       time.Time                              // time where the first data of the file has been written
	maxNumberOfBlocks  uint32                                 // total number of blocks in Archiver
	blockSize          uint32                                 // in bytes, size of each block
	lastWrittenBlock   uint32                                 // the last block that has been written into
	blockAllocationMap *blockAllocationMap.BlockAllocationMap // BAM data in memory coded with roaring, to be synced later on to Disk.
	// openFiles          map[uint32]*virtualFile.VirtualFile
	fileIndex     *fileIndex.FileIndex
	fileIndexSize uint32
	blmSize       uint32
	path          string
	log           *log.Logger
	fiChecksum    uint32

	conf         configs
	eventHandler blockAllocationMap.Events
}

func (hfs *HFileSystem) UpdateFSHeader() error {
	err := hfs.updateFileIndex()
	if err != nil {
		return err
	}

	err = hfs.updateBLM()
	if err != nil {
		return err
	}

	err = hfs.updateHeader()
	if err != nil {
		return err
	}

	err = hfs.file.Sync()
	if err != nil {
		hfs.log.Warnv("Can not sync file", "err", err.Error())
	}

	err = hfs.updateHash()
	if err != nil {
		return err
	}

	err = hfs.backUp()
	if err != nil {
		return err
	}
	return nil
}
