package Header_

import (
	"behnama/stream/pkg/archiverStorageEngine/internals/virtualFS"
	"behnama/stream/pkg/fsEngine/internal/fileIndex"
	"behnama/stream/pkg/fsEngine/pkg/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/fanap-infra/log"
)

func CreateHeaderFS(path string, size int64, blockSize uint32, log *log.Logger) (*FileSystem, error) {
	// p.log.Infov("CreateArchiver", "path", path, "size", size)
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

	file, err := utils.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		return nil, err
	}

	// make file with size
	token := make([]byte, blockSize)
	_, err = rand.Read(token)
	if err != nil {
		log.Errorv("generate rand token ", "err", err.Error())
	}
	n, err := file.WriteAt(token, size-int64(blockSize))
	if err != nil {
		log.Warnv("write token ", "err", err.Error())
	}
	if uint32(n) != blockSize {
		log.Warnv("Does not write completely ", "err", err.Error(), "n", n)
	}

	fs := &FileSystem{
		file:      file,
		size:      size,
		version:   FileSystemVersion,
		blocks:    uint32(size / int64(blockSize)),
		blockSize: blockSize,
		openFiles: make(map[uint32]*virtualFS.VirtualFile),
		fileIndex: fileIndex.NewFileIndex(),
		// lastWrittenBlock:   DataStartBlock,
		blockAllocationMap: roaring.New(),
		log:                log,
	}

	//n, err = file.WriteAt(token, int64(10*blockSize))
	//if err != nil {
	//	log.Warnv("write token ", "err", err.Error())
	//}
	//if uint32(n) != blockSize {
	//	log.Warnv("Does not write completely ", "err", err.Error(), "n", n)
	//}

	loadConf(fs)

	err = fs.updateFileIndex()
	if err != nil {
		// p.log.Errorv("updateFileIndex ", "err", err.Error())
		return nil, err
	}

	err = fs.updateHeader()
	if err != nil {
		return nil, err
	}

	//// add header to Archiver file
	//_, err = fs.writeInBlock(arc.generateArchiverHeader(), 0)
	//if err != nil {
	//	// p.log.Errorv("writeInBlock Header ", "err", err.Error())
	//	return nil, err
	//}
	//_ = arc.SyncBamToDisk()
	//// it's all redundant, just used for demonstration of what is
	//// intended.
	//arc.blockAllocationMap.AddRange(0, uint64(arc.blocks))
	//arc.blockAllocationMap.Flip(0, uint64(arc.blocks))

	return fs, nil
}

func ParseHeaderFS(path string, log *log.Logger) (*FileSystem, error) {
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

	fs := &FileSystem{
		file:      file,
		size:      size,
		openFiles: make(map[uint32]*virtualFS.VirtualFile),
		log:       log,
	}

	err = fs.parseHeader()
	if err != nil {
		return nil, err
	}

	err = fs.parseFileIndex()
	if err != nil {
		return nil, err
	}

	return fs, nil
}
