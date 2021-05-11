package fsEngine

import (
	"behnama/stream/pkg/archiverStorageEngine/internals/virtualFS"
	Header_ "behnama/stream/pkg/fsEngine/internal/Header"
	"behnama/stream/pkg/fsEngine/pkg/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanap-infra/log"
)

// name string,
func CreateFileSystem(path string, size int64, blockSize uint32, log *log.Logger) (*FSEngine, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	if blockSize < HeaderByteSize {
		return nil, fmt.Errorf("Block size must be greater than %v", blockSize)
	}

	if utils.FileExists(path) {
		return nil, errors.New("File already exists")
	}
	if size%int64(blockSize) != 0 {
		return nil, fmt.Errorf("File size must be divisible by %v", blockSize)
	}
	if size < int64(blockSize*60) {
		return nil, fmt.Errorf("File size is too small, Minimum size is %v", blockSize*60)
	}
	//+name
	file, err := utils.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o777)
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
		log.Warnv("Does not write completely ", "err", err.Error(), "n", n)
		return nil, fmt.Errorf("block size is %v, but written size is %v", blockSize, n)
	}

	fs := &FSEngine{
		file:      file,
		size:      size,
		version:   FileSystemVersion,
		blocks:    uint32(size / int64(blockSize)),
		blockSize: blockSize,
		openFiles: make(map[uint32]*virtualFS.VirtualFile),
		log:       log,
	}

	fileName := filepath.Base(path)
	headerPath := strings.Replace(path, fileName, "Header.Beh", 1)
	//, "header_"+name
	headerFS, err := Header_.CreateHeaderFS(headerPath, size, blockSize, log, fs)
	if err != nil {
		log.Errorv("Can not create header file ", "err", err.Error())
	}

	fs.header = headerFS

	return fs, nil
}

func ParseFileSystem(path string, log *log.Logger) (*FSEngine, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	size, err := utils.FileSize(path)
	if err != nil {
		return nil, err
	}
	file, err := utils.OpenFile(path, os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	fs := &FSEngine{
		file:      file,
		size:      size,
		openFiles: make(map[uint32]*virtualFS.VirtualFile),
		log:       log,
	}

	fileName := filepath.Base(path)
	headerPath := strings.Replace(path, fileName, "Header.Beh", 1)
	hfs, err := Header_.ParseHeaderFS(headerPath, log)
	if err != nil {
		return nil, err
	}

	fs.header = hfs
	fs.blockSize = hfs.GetBlockSize()
	fs.blocks = hfs.GetBlocksNumber()

	return fs, nil

	//blockCache, _ := lru.New(64)
	//arc.Cache = blockCache
	//
	//d, _, err := arc.ReadBlock(0)
	//if err != nil {
	//	log.Errorv("ReadBlock first block", "err", err.Error())
	//	return nil, err
	//}
	//
	//if string(d[:8]) != ArchiverIdentifier {
	//	if arc.tyRecover() {
	//		return ParseArchiver(file, log)
	//	}
	//	return nil, ErrArchiverIdentifier
	//}
	//if binary.BigEndian.Uint32(d[8:12]) != ArchiverVersion {
	//	if arc.tyRecover() {
	//		return ParseArchiver(file, log)
	//	}
	//	return nil, ErrArchiverVersion
	//}
	//arc.blockSize = uint32(binary.BigEndian.Uint64(d[12:20]))
	//arc.blocks = uint32(binary.BigEndian.Uint64(d[20:28]))
	//arc.lastWrittenBlock = uint32(binary.BigEndian.Uint64(d[28:36]))
	//arc.fileIndex = fileIndex.NewFileIndex()
	//arc.blockAllocationMap = roaring.New()
	//arc.CurrentFile = filepath.Base(file.Name())
	//err = arc.syncBamFromDisk()
	//if err != nil {
	//	log.Errorv("arc syncBamFromDisk", "err", err.Error())
	//}
	//arc.fileIndexIsFlip = arc.blockAllocationMap.Contains(FileIndexStartBlockFlip)
	//if !arc.readFileIndex() {
	//	_ = arc.UpdateFileIndex()
	//}
	//// checkError(err)

	// return arc, nil
}