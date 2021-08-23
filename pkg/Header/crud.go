package Header_

import "github.com/fanap-infra/fsEngine/pkg/fileIndex"

func (hfs *HFileSystem) CheckIDExist(fileID uint32) bool {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].CheckFileExistWithLock(fileID)
	}
	return hfs.fileIndexes[0].CheckFileExistWithLock(fileID)
}

func (hfs *HFileSystem) AddVirtualFile(fileID uint32, fileName string) error {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].AddFile(fileID, fileName)
	}
	return hfs.fileIndexes[0].AddFile(fileID, fileName)
}

func (hfs *HFileSystem) RemoveVirtualFile(fileID uint32) error {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].RemoveFile(fileID)
	}
	return hfs.fileIndexes[0].RemoveFile(fileID)
}

func (hfs *HFileSystem) GetFileData(fileID uint32) (fileIndex.File, error) {
	if hfs.storeInRedis {
		return hfs.fileIndexes[int(fileID)%len(hfs.fileIndexes)].GetFileInfo(fileID)
	}
	return hfs.fileIndexes[0].GetFileInfo(fileID)
}
