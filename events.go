package fsEngine

type Events interface {
	VirtualFileDeleted(fileID uint32, message string)
}
