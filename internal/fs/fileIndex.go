package fs

import (
	"bytes"
	"compress/gzip"
	"hash/crc32"
)

func (fs *FileSystem) generateFileIndex() ([]byte, error) {
	bin, err := fs.fileIndex.GenerateBinary()
	if err != nil {
		return nil, err
	}
	pw := bytes.Buffer{}

	zw, err := gzip.NewWriterLevel(&pw, 2)
	if err != nil {
		return nil, err
	}

	_, _ = zw.Write(bin)

	err = zw.Close()
	if err != nil {
		return nil, err
	}
	bin = pw.Bytes()
}

// previousBlockIndex
// returns block-index of previous blockIndex, -1 if it's the first block.
func (fs *FileSystem) UpdateFileIndex() error {
	//if !sem.TryAcquire(1) {
	//	return nil
	//}
	//defer sem.Release(1)
	//fs.fiMux.Lock()
	//defer fs.fiMux.Unlock()
	fi, err := fs.generateFileIndex()
	if err != nil {
		return err
	}

	checkSum := crc32.ChecksumIEEE(fi)
	if fs.fiChecksum == crc32.ChecksumIEEE(fi) {
		return nil
	}
	fs.fiChecksum = checkSum

	//// for mirroring
	//FileIndexStartBlock := FileIndexStartBlockFlip
	//if fs.fileIndexIsFlip {
	//	FileIndexStartBlock = FileIndexStartBlockFlop
	//}

	//blockLengthToRead := BLOCKSIZEUSABLE
	//blockNum := int(math.Ceil(float64(len(bin)) / float64(BLOCKSIZEUSABLE)))
	//for i := 0; i < blockNum; i++ {
	//	if len(bin[uint64(i)*BLOCKSIZEUSABLE:]) < BLOCKSIZEUSABLE {
	//		blockLengthToRead = len(bin[uint64(i)*BLOCKSIZEUSABLE:])
	//	}
	//	prevBlock := uint32(0)
	//	if i > FileIndexStartBlock {
	//		prevBlock = uint32(FileIndexStartBlock + i)
	//	}
	//	var binP []byte
	//	if i == blockNum-1 {
	//		binP = bin[uint64(i)*BLOCKSIZEUSABLE:]
	//	} else {
	//		binP = bin[uint64(i)*BLOCKSIZEUSABLE : (uint64(i)*BLOCKSIZEUSABLE)+uint64(blockLengthToRead)]
	//	}
	//	d := arc.PrepareBlock(binP,
	//		uint32(FileIndexStartBlock+i),
	//		prevBlock,
	//		0,
	//		i == blockNum-1)
	//	_, err = arc.WriteBlockAt(d, uint32(FileIndexStartBlock+i))
	//	if err != nil {
	//		arc.log.Error(err)
	//		return err
	//	}
	//}
	//fs.fiChecksum = crc32.ChecksumIEEE(bin)
	//if fs.fileIndexIsFlip {
	//	fs.blockAllocationMap.Remove(FileIndexStartBlockFlip)
	//} else {
	//	fs.blockAllocationMap.Remove(FileIndexStartBlockFlop)
	//}
	//fs.fileIndexIsFlip = !arc.fileIndexIsFlip

	return nil
}
