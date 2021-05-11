package fsEngine

// writeInBlock
// just writes plain data in `blockIndex` Block. Just use
// when data-block is structured.
func (fse *FSEngine) writeInBlock(data []byte, blockIndex uint32) (n int, err error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	if blockIndex >= fse.blocks {
		return 0, ErrBlockIndexOutOFRange
	}
	if len(data) != BLOCKSIZE {
		return 0, ErrDataBlockMismatch
	}

	n, err = fse.WriteAt(data, int64(blockIndex)*int64(fse.blockSize))
	if err != nil {
		fse.log.Infov("Error Writing to file", "err", err.Error(), "file", fse.file.Name())
		return n, err
	}
	// arc.setBlockAsAllocated(blockIndex)
	//if blockIndex >= fse.conf.DataStartBlock {
	//	fse.lastWrittenBlock = blockIndex
	//}
	return
}

func (fse *FSEngine) ReadBlock(blockIndex uint32) ([]byte, error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	if blockIndex >= fse.blocks {
		return nil, ErrBlockIndexOutOFRange
	}

	var err error
	buf := make([]byte, fse.blockSize)
	n, err := fse.file.ReadAt(buf, int64(blockIndex*fse.blockSize))
	if err != nil {
		return nil, err
	}
	if n != int(fse.blockSize) {
		return buf, ErrDataBlockMismatch
	}

	return buf, nil
}

func (fse *FSEngine) WriteAt(b []byte, off int64) (n int, err error) {
	n, err = fse.file.WriteAt(b, off)

	//if arc.LastFiletime.IsZero() && off >= int64(arc.conf.DataStartBlock) {
	//	arc.LastFiletime = time.Now()
	//}
	return
}
