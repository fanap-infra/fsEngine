package fileIndex

import "google.golang.org/protobuf/proto"

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
	return proto.Unmarshal(data, i.table)
}
