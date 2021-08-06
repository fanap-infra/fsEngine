package Header_

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/fanap-infra/fsEngine/internal/constants"

	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/internal/fileIndex"
	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/log"
)

func CreateHeaderFS(path string, size int64, blockSize uint32, log *log.Logger,
	eventHandler blockAllocationMap.Events) (*HFileSystem, error) {
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
		log.Warnv("does not write completely ", "err", err.Error(), "n", n)
	}

	fs := &HFileSystem{
		file:               file,
		size:               size,
		version:            FileSystemVersion,
		maxNumberOfBlocks:  uint32(size / int64(blockSize)),
		blockSize:          blockSize,
		fileIndex:          fileIndex.NewFileIndex(),
		blockAllocationMap: blockAllocationMap.New(log, eventHandler, uint32(size/int64(blockSize))),
		log:                log,
		eventHandler:       eventHandler,
		path:               path,
	}

	loadConf(fs)

	err = fs.updateFileIndex()
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

func ParseHeaderFS(path string, log *log.Logger, eventHandler blockAllocationMap.Events) (*HFileSystem, error) {
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
		file:         file,
		size:         size,
		fileIndex:    fileIndex.NewFileIndex(),
		path:         path,
		log:          log,
		eventHandler: eventHandler,
	}

	err = hfs.parseHeader()
	if err != nil {
		return hfs, err
	}

	err = hfs.parseFileIndex()
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
