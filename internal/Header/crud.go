package Header_

import "github.com/fanap-infra/fsEngine/pkg/fileIndex"

func (hfs *HFileSystem) CheckIDExist(id uint32) bool {
	return hfs.fileIndex.CheckFileExistWithLock(id)
}

func (hfs *HFileSystem) AddVirtualFile(id uint32, fileName string) error {
	return hfs.fileIndex.AddFile(id, fileName)
}

func (hfs *HFileSystem) RemoveVirtualFile(id uint32) error {
	return hfs.fileIndex.RemoveFile(id)
}

func (hfs *HFileSystem) GetFileData(id uint32) (fileIndex.File, error) {
	return hfs.fileIndex.GetFileInfo(id)
}
