package virtualFile

type FS interface {
	Write(data []byte, fileID uint32) (int, error)
	WriteAt(data []byte, off int64, fileID uint32) (int, error)
	Read(data []byte, fileID uint32) (int, error)
	ReadAt(data []byte, off int64, fileID uint32) (int, error)
	ReadBlock(blockIndex uint32) ([]byte, error)
	BAMUpdated(fileID uint32, bam []byte) error // update block allocation map byte array
	UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32) error
	Closed(fileID uint32) error
}
