package Header_

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/fanap-infra/fsEngine/internal/constants"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/fileIndex"
	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/log"
)

// if redis options is nil, use file for storing data
func CreateHeaderFS(id uint32, path string, size int64, blockSize uint32, log *log.Logger,
	eventHandler blockAllocationMap.Events, redisDB RedisDB) (*HFileSystem, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	headerPath := path + "/" + constants.HeaderPath
	if blockSize < HeaderByteSize {
		return nil, fmt.Errorf("block size must be greater than %v", blockSize)
	}

	if utils.FileExists(headerPath) {
		return nil, errors.New("file already exists")
	}
	if size%int64(blockSize) != 0 {
		return nil, fmt.Errorf("file size must be divisible by %v", blockSize)
	}
	if size < int64(blockSize*60) {
		return nil, fmt.Errorf("file size is too small, Minimum size is %v", blockSize*60)
	}

	file, err := utils.OpenFile(headerPath, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	// make file with size
	token := make([]byte, blockSize)
	_, err = rand.Read(token)
	if err != nil {
		log.Errorv("generate rand token ", "err", err.Error())
	}
	n, err := file.WriteAt(token, HashByteIndex)
	if err != nil {
		log.Warnv("write token ", "err", err.Error())
	}
	if uint32(n) != blockSize {
		log.Warnv("does not write completely ", "n", n)
	}

	fs := &HFileSystem{
		id:                 id,
		file:               file,
		size:               size,
		version:            FileSystemVersion,
		maxNumberOfBlocks:  uint32(size / int64(blockSize)),
		blockSize:          blockSize,
		blockAllocationMap: blockAllocationMap.New(log, eventHandler, uint32(size/int64(blockSize))),
		log:                log,
		eventHandler:       eventHandler,
		path:               path,
		storeInRedis:       redisDB != nil,
	}

	if redisDB != nil {
		fs.redisClient = redisDB
		// ToDo: add numberOfFileIndexes to header
		for i := 0; i < numberOfFileIndexes; i++ {
			fs.fileIndexes = append(fs.fileIndexes, fileIndex.NewFileIndex())
		}
	} else {
		fs.fileIndexes = append(fs.fileIndexes, fileIndex.NewFileIndex())
	}

	loadConf(fs)

	err = fs.updateAllFileIndex()
	if err != nil {
		return nil, err
	}

	err = fs.updateBLM()
	if err != nil {
		return nil, err
	}

	err = fs.updateHeader()
	if err != nil {
		return nil, err
	}

	//err = fs.updateHash()
	//if err != nil {
	//	return nil, err
	//}

	return fs, nil
}

func ParseHeaderFS(id uint32, path string, log *log.Logger, eventHandler blockAllocationMap.Events,
	redisDB RedisDB) (*HFileSystem, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	headerPath := path + "/" + constants.HeaderPath
	size, err := utils.FileSize(headerPath)
	if err != nil {
		return nil, err
	}
	file, err := utils.OpenFile(headerPath, os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	hfs := &HFileSystem{
		id:           id,
		file:         file,
		size:         size,
		path:         path,
		log:          log,
		eventHandler: eventHandler,
		storeInRedis: redisDB != nil,
	}

	if redisDB != nil {
		hfs.redisClient = redisDB
		// ToDo: add numberOfFileIndexes to header
		for i := 0; i < numberOfFileIndexes; i++ {
			hfs.fileIndexes = append(hfs.fileIndexes, fileIndex.NewFileIndex())
		}
	} else {
		hfs.fileIndexes = append(hfs.fileIndexes, fileIndex.NewFileIndex())
	}

	err = hfs.parseHeader()
	if err != nil {
		return hfs, err
	}

	err = hfs.parseAllFileIndexes()
	if err != nil {
		return hfs, err
	}

	err = hfs.parseBLM()
	if err != nil {
		return hfs, err
	}

	//if !hfs.checkHash() {
	//	hfs.log.Warn("hash value of header file is not correct")
	//}
	return hfs, nil
}
