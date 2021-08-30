package fileIndex

import (
	"fmt"

	"github.com/fanap-infra/log"
	"google.golang.org/protobuf/proto"
)

// NewFileIndex
func NewFileIndex() (fi *FileIndex) {
	return &FileIndex{table: &Table{Files: make(map[uint32]*File), NumberFiles: 0}}
}

func (i *FileIndex) GenerateBinary() (data []byte, err error) {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	data, err = proto.Marshal(i.table)
	return
}

func (i *FileIndex) InitFromBinary(data []byte) error {
	i.rwMux.Lock()
	defer i.rwMux.Unlock()
	err := proto.Unmarshal(data, i.table)
	if err == nil && i.table.Files == nil {
		log.Errorv("it does not init fileIndex correctly by Protobuf binary", "len(data)", len(data))
		i.table.Files = make(map[uint32]*File)
		return fmt.Errorf("it does not init fileIndex correctly")
	}
	return err
}
