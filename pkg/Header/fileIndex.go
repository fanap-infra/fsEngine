package Header_

import (
	"fmt"
	"time"

	"github.com/fanap-infra/log"

	"github.com/fanap-infra/fsEngine/pkg/fileIndex"
)

func (hfs *HFileSystem) generateFileIndex(index uint32) ([]byte, error) {
	if len(hfs.fileIndexes) <= int(index) {
		return nil, fmt.Errorf("index :%v is out of range: %v", index, len(hfs.fileIndexes))
	}
	bin, err := hfs.fileIndexes[index].GenerateBinary()
	if err != nil {
		return nil, err
	}

	return bin, nil
}

func (hfs *HFileSystem) updateFileIndex(index uint32) error {
	fi, err := hfs.generateFileIndex(index)
	if err != nil {
		return err
	}
	if len(fi) > FileIndexMaxByteSize {
		return fmt.Errorf("fileIndex size %v is too large, Max valid size: %v",
			len(fi), FileIndexMaxByteSize)
	}
	hfs.fileIndexSize = uint32(len(fi))

	if hfs.fileIndexSize == 0 {
		hfs.log.Warn("file indexes size is zero")
	}

	if hfs.storeInRedis {
		err := hfs.setRedisKeyValue("arch"+fmt.Sprint(hfs.id)+"_fileIndex"+fmt.Sprint(int(index)%len(hfs.fileIndexes)), fi)
		if err != nil {
			return err
		}
		return nil
	}

	n, err := hfs.writeAt(fi, HeaderBlockIndex)
	if err != nil {
		return err
	}
	if n != len(fi) {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}

	n, err = hfs.writeEOPart(int64(FileIndexByteIndex + n))
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}

	return nil
}

func (hfs *HFileSystem) updateAllFileIndex() error {
	for index := range hfs.fileIndexes {
		err := hfs.updateFileIndex(uint32(index))
		if err != nil {
			hfs.log.Errorv("can not update file index", "index", index, "err", err.Error())
			return err
		}
	}
	return nil
}

func (hfs *HFileSystem) parseFileIndex(index uint32) error {
	var buf []byte
	var err error
	if hfs.storeInRedis {
		buf, err = hfs.getRedisValue("arch" + fmt.Sprint(hfs.id) + "_fileIndex" + fmt.Sprint(int(index)%len(hfs.fileIndexes)))
		if err != nil {
			hfs.log.Errorv("can get value from redis", "key", "arch"+fmt.Sprint(hfs.id)+"_fileIndex"+fmt.Sprint(int(index)%len(hfs.fileIndexes)),
				"err", err.Error())
			return err
		}
	} else {
		buf = make([]byte, hfs.fileIndexSize)
		if hfs.fileIndexSize == 0 {
			hfs.log.Warnv("file index size is zero", "fileIndexSize", hfs.fileIndexSize)
			return nil
		}
		n, err := hfs.readAt(buf, FileIndexByteIndex)
		if err != nil {
			return err
		}
		if n != int(hfs.fileIndexSize) {
			return ErrDataBlockMismatch
		}
	}

	if len(buf) == 0 {
		hfs.log.Warnv("parse file index, buf length is zero", "index", index, "id", hfs.id,
			"redisKey", "arch"+fmt.Sprint(hfs.id)+"_fileIndex"+fmt.Sprint(int(index)%len(hfs.fileIndexes)))
		return nil
	}
	err = hfs.fileIndexes[index].InitFromBinary(buf)
	if err != nil {
		hfs.log.Warnv("parse file index, can not init protobuf", "index", index, "id", hfs.id,
			"redisKey", "arch"+fmt.Sprint(hfs.id)+"_fileIndex"+fmt.Sprint(int(index)%len(hfs.fileIndexes)),
			"err", err.Error())
		return err
	}

	return nil
}

func (hfs *HFileSystem) parseAllFileIndexes() error {
	for i := 0; i < numberOfFileIndexes; i++ {
		err := hfs.parseFileIndex(uint32(i))
		if err != nil {
			hfs.log.Errorv("can not parse file index", "i", i, "err", err.Error())
			return err
		}
	}
	return nil
}

// ToDo:update blm binaries
func (hfs *HFileSystem) UpdateBAM(fileID uint32, data []byte) error {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].UpdateBAM(fileID, data)
	}
	return hfs.fileIndexes[0].UpdateBAM(fileID, data)
}

func (hfs *HFileSystem) UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32,
	fileSize uint32, bam []byte, info []byte) error {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].UpdateFileIndexes(fileID, firstBlock, lastBlock, fileSize, bam, info)
	}
	return hfs.fileIndexes[0].UpdateFileIndexes(fileID, firstBlock, lastBlock, fileSize, bam, info)
}

func (hfs *HFileSystem) FindOldestFile() (*fileIndex.File, error) {
	if hfs.storeInRedis {
		return hfs.findOldestBetweenFileIndexes()
	}
	return hfs.fileIndexes[0].FindOldestFile()
}

func (hfs *HFileSystem) findOldestBetweenFileIndexes() (*fileIndex.File, error) {
	oldestTime := time.Now().Local()
	var foundedFile *fileIndex.File
	for i, fIndex := range hfs.fileIndexes {
		oldestFile, err := fIndex.FindOldestFile()
		if err != nil {
			log.Errorv("can not parse file created time", "i", i,
				"err", err.Error())
			continue
		}
		// timestamppb.New(oldestFile.CreatedTime)

		createdTime := oldestFile.CreatedTime.AsTime()
		//if err != nil {
		//	log.Errorv("can not parse file created time", "i", i,
		//		"err", err.Error())
		//	continue
		//}
		if oldestTime.After(createdTime) {
			foundedFile = oldestFile
			oldestTime = createdTime
		}
	}
	return foundedFile, nil
}

func (hfs *HFileSystem) UpdateFileOptionalData(fileID uint32, info []byte) error {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].UpdateFileOptionalData(fileID, info)
	}
	return hfs.fileIndexes[0].UpdateFileOptionalData(fileID, info)
}

func (hfs *HFileSystem) GetFilesList() []*fileIndex.File {
	return hfs.fileIndexes[0].GetFilesList()
}

//func (hfs *HFileSystem) GetFileOptionalData(fileId uint32) ([]byte, error) {
//	return hfs.fileIndex.
//}
