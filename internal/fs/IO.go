package fs

// writeInBlock
// just writes plain data in `blockIndex` Block. Just use
// when data-block is structured.
func (fs *FileSystem) writeInBlock(data []byte, blockIndex uint32) (n int, err error) {
	fs.rIBlockMux.Lock()
	defer fs.rIBlockMux.Unlock()
	if blockIndex >= fs.blocks {
		return 0, ErrBlockIndexOutOFRange
	}
	if len(data) != BLOCKSIZE {
		return 0, ErrDataBlockMismatch
	}

	n, err = fs.WriteAt(data, int64(blockIndex)*int64(fs.conf.BLOCKSIZE))
	if err != nil {
		fs.log.Infov("Error Writing to file", "err", err.Error(), "file", fs.file.Name())
		return n, err
	}
	// arc.setBlockAsAllocated(blockIndex)
	if blockIndex >= fs.conf.DataStartBlock {
		fs.lastWrittenBlock = blockIndex
	}
	return
}

func (fs *FileSystem) ReadBlock(blockIndex uint32) ([]byte, error) {
	fs.rIBlockMux.Lock()
	defer fs.rIBlockMux.Unlock()
	if blockIndex >= fs.blocks {
		return nil, ErrBlockIndexOutOFRange
	}

	var err error
	buf := make([]byte, fs.blockSize)
	n, err := fs.file.ReadAt(buf, int64(blockIndex*fs.blockSize))
	if err != nil {
		return nil, err
	}
	if n != int(fs.blockSize) {
		return buf, ErrDataBlockMismatch
	}

	return buf, nil
}

func (fs *FileSystem) WriteAt(b []byte, off int64) (n int, err error) {
	n, err = fs.file.WriteAt(b, off)

	//if arc.LastFiletime.IsZero() && off >= int64(arc.conf.DataStartBlock) {
	//	arc.LastFiletime = time.Now()
	//}
	return
}
