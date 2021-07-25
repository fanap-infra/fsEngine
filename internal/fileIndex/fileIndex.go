package fileIndex

import (
	"fmt"
	"sync"
	"time"

	"github.com/fanap-infra/log"

	"github.com/golang/protobuf/ptypes"
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
	createdTime, err := ptypes.TimestampProto(time.Now().Local())
	if err != nil {
		return err
	}
	i.table.NumberFiles++
	i.table.Files[fileId] = &File{
		Id: fileId, FirstBlock: 0, LastBlock: 0,
		Name: name, RMapBlocks: make([]byte, 0),
		CreatedTime: createdTime,
	}
	return nil
}

func (i *FileIndex) FindOldestFile() (*File, error) {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	oldestTime := time.Now().Local()
	var foundedFile *File
	foundedFile = nil
	if len(i.table.Files) == 0 {
		return nil, fmt.Errorf("there is no file")
	}
	for _, file := range i.table.Files {
		createdTime, err := ptypes.Timestamp(file.CreatedTime)
		if err != nil {
			log.Errorv("can not parse file created time", "id", file.Id,
				"err", err.Error())
			continue
		}
		if oldestTime.After(createdTime) {
			foundedFile = file
			oldestTime = createdTime
		}
	}

	return foundedFile, nil
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

func (i *FileIndex) GetFileInfo(fileId uint32) (File, error) {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	fileInfo, isExist := i.table.Files[fileId]
	if !isExist {
		return File{}, fmt.Errorf("file id %v does not exist", fileId)
	}

	return File{
		Id: fileInfo.Id, RMapBlocks: fileInfo.RMapBlocks, FirstBlock: fileInfo.FirstBlock,
		LastBlock: fileInfo.LastBlock, Name: fileInfo.Name, Optional: fileInfo.Optional,
		FileSize: fileInfo.GetFileSize(),
	}, nil
}

func (i *FileIndex) UpdateBAM(fileId uint32, bam []byte) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	if !i.checkFileExist(fileId) {
		return fmt.Errorf("file id %v does not exist", fileId)
	}
	i.table.Files[fileId].RMapBlocks = bam
	return nil
}

func (i *FileIndex) UpdateFileIndexes(fileId uint32, firstBlock uint32, lastBlock uint32, fileSize uint32) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	if !i.checkFileExist(fileId) {
		return fmt.Errorf("file id %v does not exist", fileId)
	}
	i.table.Files[fileId].FirstBlock = firstBlock
	i.table.Files[fileId].LastBlock = lastBlock
	i.table.Files[fileId].FileSize = fileSize
	return nil
}

func (i *FileIndex) UpdateFileOptionalData(fileId uint32, info []byte) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	if !i.checkFileExist(fileId) {
		return fmt.Errorf("file id %v does not exist", fileId)
	}
	i.table.Files[fileId].Optional = info

	return nil
}

//func (i *FileIndex) GetFileOptionalData(fileId uint32) ([]byte, error) {
//	i.rwMux.Lock()
//	defer i.rwMux.Unlock()
//	if !i.checkFileExist(fileId) {
//		return nil, fmt.Errorf("file id %v does not exist", fileId)
//	}
//
//	return nil
//}
