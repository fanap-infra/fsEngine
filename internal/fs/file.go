package fs

import (
	"behnama/stream/pkg/archiverStorageEngine/internals/fileIndex"
	"behnama/stream/pkg/archiverStorageEngine/internals/virtualFS"
	"behnama/stream/pkg/fsEngine/pkg/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fanap-infra/log"

	lru "github.com/hashicorp/golang-lru"
	"github.com/influxdata/roaring"
)

// File
type FileSystem struct {
	file               *os.File // file handle instance
	version            uint32
	size               int64
	CurrentFile        string          // name of the latest file to be created
	LastFiletime       time.Time       // time where the first data of the file has been written
	blocks             uint32          // total number of blocks in Archiver
	blockSize          uint32          // in bytes, size of each block
	lastWrittenBlock   uint32          // the last block that has been written into
	blockAllocationMap *roaring.Bitmap // BAM data in memory coded with roaring, to be synced later on to Disk.
	openFiles          map[uint32]*virtualFS.VirtualFile
	fileIndex          fileIndex.FileIndex
	WMux               sync.Mutex
	RMux               sync.Mutex
	log                *log.Logger
	fiMux              sync.RWMutex
	fiChecksum         uint32
	bamChecksum        uint32
	fsMux              sync.Mutex
	rIBlockMux         sync.Mutex
	Cache              *lru.Cache
	fileIndexIsFlip    bool
	EventsHandler      Events
	conf               configs
	Quit               chan struct{}
}

func CreateFileSystem(path string, size int64, blockSize uint32, log *log.Logger) (*FileSystem, error) {
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
	if size < int64(blockSize*30) {
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

	err = fs.updateHeader()
	if err != nil {
		return nil, err
	}

	err = fs.UpdateFileIndex()
	if err != nil {
		// p.log.Errorv("updateFileIndex ", "err", err.Error())
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

func loadConf(f *FileSystem) {
	f.conf.BLOCKSIZE = f.blockSize // 1 << 19
	f.conf.BLOCKSIZEUSABLE = f.conf.BLOCKSIZE - 20
	f.conf.StorageMaxSize = 1 << 44
	f.conf.BlockAllocationMapSize = uint32(f.conf.StorageMaxSize/uint64(f.conf.BLOCKSIZEUSABLE)) / 8 // Size in blocks.
	f.conf.FileIndexReservedSize = 1 << 22
	f.conf.FileIndexReservedSizeBlocks = (f.conf.FileIndexReservedSize / f.conf.BLOCKSIZEUSABLE) * 2 //(FileIndexReservedSize / BLOCKSIZEUSABLE) * 2
	f.conf.DataStartBlock = (f.conf.BlockAllocationMapSize / f.conf.BLOCKSIZEUSABLE) + f.conf.FileIndexReservedSizeBlocks + 1
	f.lastWrittenBlock = f.conf.DataStartBlock
}

func ParseFileSystem(path string, log *log.Logger) (*FileSystem, error) {
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

// Close ...
func (fs *FileSystem) Close() error {
	err := fs.updateHeader()
	if err != nil {
		fs.log.Warnv("Can not updateHeader", "err", err.Error())
		// ToDo: remove it
		return err
	}
	// ToDo:update file system
	err = fs.file.Sync()
	if err != nil {
		fs.log.Warnv("Can not sync file", "err", err.Error())
		// ToDo: remove it
		return err
	}
	return fs.file.Close()
}
