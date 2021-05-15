package fileIndex

import (
	"fmt"
	"sync"
)

const (
	HashTableSize = 0
)

type FileIndex struct {
	table *Table
	rwMux sync.Mutex
}

func (i *FileIndex) AddFile(fileId uint32, name string) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	if i.checkFileExist(fileId) {
		return fmt.Errorf("file id %v has been added before", fileId)
	}
	i.table.NumberFiles++
	i.table.Files[fileId] = &File{Id: fileId, FirstBlock: 0, LastBlock: 0, Name: name, RMapBlocks: make([]byte, 0)}
	return nil
}

func (i *FileIndex) checkFileExist(fileId uint32) bool {
	_, isExist := i.table.Files[fileId]
	return isExist
}

func (i *FileIndex) CheckFileExistWithLock(fileId uint32) bool {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	_, isExist := i.table.Files[fileId]
	return isExist
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

func (i *FileIndex) UpdateFile(fileId uint32, firstBlock uint32, lastBlock uint32, name string, blocks []byte) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	if !i.checkFileExist(fileId) {
		return fmt.Errorf("file id %v does not exist", fileId)
	}
	i.table.Files[fileId] = &File{Id: fileId, FirstBlock: firstBlock, LastBlock: lastBlock, Name: name, RMapBlocks: blocks}
	return nil
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
