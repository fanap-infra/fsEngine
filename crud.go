package fsEngine

import (
	"github.com/fanap-infra/FSEngine/internal/virtualFile"
	"fmt"
)


func (fse *FSEngine) NewVirtualFile(id uint32, fileName string) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	if fse.header.CheckIDExist(id) {
		return nil, fmt.Errorf("this ID: %v, had been taken", id)
	}
	vf := virtualFile.NewVirtualFile(fileName, id, fse, fse, fse.log)
	err := fse.header.AddVirtualFile(id, fileName)
	if err != nil {
		return nil, err
	}
	fse.openFiles[id] = vf
	return vf, nil
}

func (fse *FSEngine) OpenVirtualFile(id uint32) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	_, ok := fse.openFiles[id]
	if ok {
		return nil, fmt.Errorf("this ID: %v is opened before", id)
	}
	fileInfo , err := fse.header.GetFileData(id)
	if err != nil{
		return nil, err
	}
	fileInfo, := virtualFile.OpenVirtualFile()
	err := fse.header.AddVirtualFile(id, fileName)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (fse *FSEngine) DeleteVirtualFile(id uint32) error {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	_, ok := fse.openFiles[id]
	if ok {
		return fmt.Errorf("virtual file id : %d is opened", id)
	}
	//fse.header.
	return nil
}


