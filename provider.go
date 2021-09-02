package fsEngine

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/fanap-infra/fsEngine/pkg/virtualFile"

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

func RecoverHeaderFileSystem(id uint32, path string, blockSize uint32, eventsHandler Events, log *log.Logger, redisDB Header_.RedisDB) (*FSEngine, error) {
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

	fse := &FSEngine{
		id:                id,
		file:              file,
		size:              size,
		blockSize:         blockSize,
		blockSizeUsable:   blockSize - constants.BlockHeaderSize,
		maxNumberOfBlocks: uint32(size / int64(blockSize)),
		openFiles:         make(map[uint32]*VFInfo),
		log:               log,
		eventsHandler:     eventsHandler,
	}

	hfs, err := Header_.CreateHeaderFSForRecovering(id, path, size, blockSize, log, fse, redisDB)
	if err != nil {
		fse.log.Errorv("Can not create header file ", "err", err.Error())
		return nil, err
	}

	fse.header = hfs
	buf := make([]byte, fse.blockSize)
	blockIndex := uint32(0)
	validBlocks := 0
	vfs := make(map[uint32]*virtualFile.VirtualFile)
	for {
		if blockIndex >= fse.maxNumberOfBlocks {
			fse.log.Infov("recovering finished", "validBlocks", validBlocks, "blockIndex", blockIndex)
			break
		}
		blockID := blockIndex
		n, err := fse.file.ReadAt(buf, int64(blockIndex)*int64(fse.blockSize))
		blockIndex++
		if err != nil {
			fse.log.Warnv("can not read block completely", "n", n, "blockSize", fse.blockSize)
			continue
		}
		if n != int(fse.blockSize) {
			continue
		}

		pBlockID := binary.BigEndian.Uint32(buf[0:4])
		pFileID := binary.BigEndian.Uint32(buf[4:8])
		dataSize := binary.BigEndian.Uint32(buf[12:16])
		if dataSize > fse.blockSize-16 {
			log.Warnv("block data size is too large, ", "dataSize", dataSize, "blockID", blockID)
			continue
		}
		if pBlockID != blockID {
			// fse.log.Warnv("blockd id is wrong,", "pBlockID", pBlockID, "blockID", blockID)
			continue
		}

		vf, isExist := vfs[pFileID]
		if !isExist {
			vf, err = fse.NewVirtualFile(pFileID, "test")
			if err != nil {
				fse.log.Warnv("can not add virtual file correctly,", "pFileID", pFileID, "err", err.Error())
				continue
			}
			vfs[pFileID] = vf
		}

		err = vf.AddBlockID(blockID)
		if err != nil {
			fse.log.Warnv("can not add blockID to vf,", "pFileID", pFileID, "blockID", blockID,
				"err", err.Error())
			continue
		}

		vf.AddFileSize(dataSize)

		err = fse.header.SetBlockAsAllocated(blockID)
		if err != nil {
			fse.log.Errorv("can not set block id in header",
				"blockID", blockID, "pFileID", pFileID)
			continue
		}
	}

	vfBlockCounter := uint32(0)
	for i, vf := range vfs {
		vfBlockCounter = vfBlockCounter + uint32(len(vf.GetBLMArray()))
		err := vf.Close()
		if err != nil {
			fse.log.Errorv("can not close virtual file",
				"i", i, "fileID", vf.GetFileID())
		}
	}

	fileIndexes := fse.header.GetFilesList()
	fse.log.Infov("recovery report",
		"len(fileIndexes)", len(fileIndexes), "number of blocks", len(fse.header.GetBLMArray()),
		"sum of vfsBlocks", vfBlockCounter)

	err = hfs.UpdateFSHeader()
	if err != nil {
		log.Errorv("can not update header", "err", err.Error())
	}
	return fse, nil
}
