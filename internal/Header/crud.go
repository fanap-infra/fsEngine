package Header_

func (hfs *HFileSystem) CheckIDExist(id uint32) bool {
	return hfs.fileIndex.CheckFileExistWithLock(id)
}

func (hfs *HFileSystem) AddVirtualFile(id uint32, fileName string) error {
	return hfs.fileIndex.AddFile(id, fileName)
}
