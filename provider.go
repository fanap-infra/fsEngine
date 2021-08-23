package fsEngine

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/fanap-infra/fsEngine/internal/constants"

	Header_ "github.com/fanap-infra/fsEngine/pkg/Header"
	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
)

func CreateFileSystem(id uint32, path string, size int64, blockSize uint32,
	eventsHandler Events, log *log.Logger, redisDB Header_.RedisDB) (*FSEngine, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	if blockSize < constants.HeaderByteSize {
		return nil, fmt.Errorf("Block size must be greater than %v", blockSize)
	}
	filePath := path + "/" + constants.FsPath
	headerPath := path + "/" + constants.HeaderPath
	if utils.FileExists(path) {
		return nil, errors.New("file already exists")
	}
	if utils.FileExists(headerPath) {
		return nil, errors.New("header file already exists")
	}
	if size%int64(blockSize) != 0 {
		return nil, fmt.Errorf("File size must be divisible by %v", blockSize)
	}
	if size < int64(blockSize*60) {
		return nil, fmt.Errorf("File size is too small, Minimum size is %v", blockSize*60)
	}

	file, err := utils.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	token := make([]byte, blockSize)
	_, err = rand.Read(token)
	if err != nil {
		log.Errorv("generate rand token ", "err", err.Error())
		return nil, err
	}
	n, err := file.WriteAt(token, size-int64(blockSize))
	if err != nil {
		log.Warnv("write token ", "err", err.Error())
		return nil, err
	}
	if uint32(n) != blockSize {
		log.Warnv("Does not write completely ", "n", n)
		return nil, fmt.Errorf("block size is %v, but written size is %v", blockSize, n)
	}

	fs := &FSEngine{
		id:                id,
		file:              file,
		size:              int64(uint32(size/int64(blockSize)) * blockSize),
		version:           constants.FileSystemVersion,
		maxNumberOfBlocks: uint32(size / int64(blockSize)),
		blockSize:         blockSize,
		blockSizeUsable:   blockSize - constants.BlockHeaderSize,
		openFiles:         make(map[uint32]*VFInfo),
		eventsHandler:     eventsHandler,
		log:               log,
	}

	headerFS, err := Header_.CreateHeaderFS(id, path, size, blockSize, log, fs, redisDB)
	if err != nil {
		log.Errorv("Can not create header file ", "err", err.Error())
		return nil, err
	}

	fs.header = headerFS

	return fs, nil
}

func ParseFileSystem(id uint32, path string, eventsHandler Events, log *log.Logger, redisDB Header_.RedisDB) (*FSEngine, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	filePath := path + "/" + constants.FsPath
	size, err := utils.FileSize(filePath)
	if err != nil {
		return nil, err
	}
	file, err := utils.OpenFile(filePath, os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	fs := &FSEngine{
		id:            id,
		file:          file,
		size:          size,
		openFiles:     make(map[uint32]*VFInfo),
		log:           log,
		eventsHandler: eventsHandler,
	}

	// fileName := filepath.Base(path)
	// headerPath := strings.Replace(path, fileName, "Header.Beh", 1)
	hfs, err := Header_.ParseHeaderFS(id, path, log, fs, redisDB)
	if err != nil {
		return nil, err
	}

	fs.header = hfs
	fs.blockSize = hfs.GetBlockSize()
	fs.blockSizeUsable = fs.blockSize - constants.BlockHeaderSize
	fs.maxNumberOfBlocks = hfs.GetBlocksNumber()

	return fs, nil
}
