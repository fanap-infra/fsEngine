package virtualFile

type FS interface {
	Write(data []byte, fileID uint32) (int, error)
	WriteAt(data []byte, off int64, fileID uint32) (int, error)
	Read(data []byte, fileID uint32) (int, error)
	ReadAt(data []byte, off int64, fileID uint32) (int, error)
	Closed(fileID uint32) error
}
