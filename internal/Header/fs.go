package Header_

import (
	"os"
	"sync"
	"time"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/fileIndex"
	"github.com/go-redis/redis/v8"

	"github.com/fanap-infra/log"
)

type HFileSystem struct {
	id                 uint32
	file               *os.File // file handle instance
	wMux               sync.Mutex
	version            uint32
	size               int64
	CurrentFile        string                                 // name of the latest file to be created
	LastFiletime       time.Time                              // time where the first data of the file has been written
	maxNumberOfBlocks  uint32                                 // total number of blocks in Archiver
	blockSize          uint32                                 // in bytes, size of each block
	lastWrittenBlock   uint32                                 // the last block that has been written into
	blockAllocationMap *blockAllocationMap.BlockAllocationMap // BAM data in memory coded with roaring, to be synced later on to Disk.
	fileIndexes        []*fileIndex.FileIndex
	fileIndexSize      uint32
	blmSize            uint32
	path               string
	log                *log.Logger
	mu                 sync.Mutex
	conf               configs
	eventHandler       blockAllocationMap.Events
	storeInRedis       bool
	redisClient        *redis.Client
}

func (hfs *HFileSystem) UpdateFSHeader() error {
	hfs.mu.Lock()
	defer hfs.mu.Unlock()

	err := hfs.updateAllFileIndex()
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

	//err = hfs.file.Sync()
	//if err != nil {
	//	hfs.log.Warnv("Can not sync file", "err", err.Error())
	//}

	//err = hfs.updateHash()
	//if err != nil {
	//	return err
	//}

	//err = hfs.backUp()
	//if err != nil {
	//	return err
	//}
	return nil
}
