package Header_

import (
	"fmt"
)

func (hfs *HFileSystem) CheckIDExist(id uint32) bool {
	_, ok := hfs.openFiles[id]
	return ok
}

func (hfs *HFileSystem) AddVirtualFile(id uint32, fileName string) error {
	if hfs.CheckIDExist(id) {
		return fmt.Errorf("this ID: %v, had been taken", id)
	}

	// fs.fileIndex.AddFile()
	// fs.openFiles[id] := virtualFile.NewVirtualFile(fileName,id, )
	return nil
}
