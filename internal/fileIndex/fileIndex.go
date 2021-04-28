package fileIndex

import (
	"fmt"
	"sync"

	"github.com/RoaringBitmap/roaring"
)

const (
	HashTableSize = 0
)

type FileIndex struct {
	table *Table
	rwMux sync.Mutex
}

type FileMetadata struct {
	FirstBlock uint32
	LastBlock  uint32
	Blocks     *roaring.Bitmap
}

/*type fileIndex interface {
	AddFile(fileId uint32) bool
	RemoveFile(fileId uint32) bool
	EditFileMeta(fileId uint32, meta FileMetadata) bool
	GenerateBinary() (data []byte, err error)
}*/

//// InitFileIndex
//func InitFileIndex(data []byte) (fi FileIndex, err error) {
//	hash := CreateHashTable(HashTableSize)
//	err = proto.Unmarshal(data, &hash)
//
//	fi.hash = &hash
//	return
//}

// AddFile
func (i *FileIndex) AddFile(fileId uint32) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	_, isExist := i.table.Files[fileId]
	if isExist {
		return fmt.Errorf("file id %v has been added before", fileId)
	}
	i.table.NumberFiles++
	i.table.Files[fileId] = &File{Id: fileId, FirstBlock: 0, LastBlock: 0}
	return nil
}

// RemoveFile
func (i *FileIndex) RemoveFile(fileId uint32) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	_, isExist := i.table.Files[fileId]
	if isExist {
		delete(i.table.Files, fileId)
		return nil
	}

	return fmt.Errorf("file id %v does not exist", fileId)
}

//// EditFileMeta
//func (i *FileIndex) EditFileMeta(fileId uint32, meta FileMetadata) bool {
//	i.rwMux.Lock()
//	defer i.rwMux.Unlock()
//
//	block := bytes.Buffer{}
//	meta.Blocks.RunOptimize()
//	_, _ = meta.Blocks.WriteTo(&block)
//	i.hash.Put(fileId, meta.FirstBlock, meta.LastBlock, block.Bytes())
//	return true
//}
//
//// GetFileInfo
//func (i *FileIndex) GetFileInfo(fileId uint32) (meta *FileMetadata, err error) {
//	ok, v := i.hash.Get(fileId)
//	if !ok {
//		return nil, errors.New("file info cannot be retrieved")
//	}
//	meta = &FileMetadata{v.FirstBlock, v.LastBlock, v.Blocks}
//
//	return meta, nil
//}
