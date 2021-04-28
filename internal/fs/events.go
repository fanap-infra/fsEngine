package fs

type Events interface {
	DeleteFile(fileID uint32)
	DeleteFileByArchiver(archiverFile string)
}
