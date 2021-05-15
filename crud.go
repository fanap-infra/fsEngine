package fsEngine

import (
	"github.com/fanap-infra/FSEngine/internal/virtualFile"
	"fmt"
)

// ToDO: completeCRUD

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
	return vf, nil
	// return nil, nil
}

func (fse *FSEngine) OpenVirtualFile(id uint32) (*virtualFile.VirtualFile, error) {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	if !fse.header.CheckIDExist(id) {
		return nil, fmt.Errorf("this ID: %v does not exist", id)
	}
	//vf := virtualFile.OpenVirtualFile()
	//err := fse.header.AddVirtualFile(id, fileName)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func (fse *FSEngine) DeleteVirtualFile(id uint32) error {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	return nil
}
