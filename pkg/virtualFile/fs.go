package virtualFile

type FS interface {
	Write(data []byte, fileID uint32, previousBlock uint32) (int, []uint32, error)
	// WriteAt(data []byte, off int64, fileID uint32) (int, error)
	Read(data []byte, fileID uint32) (int, error)
	ReadAt(data []byte, off int64, fileID uint32) (int, error)
	ReadBlock(blockIndex uint32, fileID uint32) ([]byte, error)
	BAMUpdated(fileID uint32, bam []byte) error // update block allocation map byte array
	UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32, fileSize uint32, bam []byte, info []byte) error
	UpdateFileOptionalData(fileId uint32, info []byte) error
	Closed(fileID uint32) error
}
