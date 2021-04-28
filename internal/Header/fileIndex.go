package Header_

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
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
	return bin, nil
}

// previousBlockIndex
// returns block-index of previous blockIndex, -1 if it's the first block.
func (fs *FileSystem) updateFileIndex() error {
	fi, err := fs.generateFileIndex()
	if err != nil {
		return err
	}

	checkSum := crc32.ChecksumIEEE(fi)
	if fs.fiChecksum == crc32.ChecksumIEEE(fi) {
		return nil
	}
	fs.fiChecksum = checkSum
	fs.fileIndexSize = uint32(len(fi))
	n, err := fs.file.WriteAt(fi, FileIndexByteIndex)
	if err != nil {
		return err
	}
	if n != len(fi) {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}
	n, err = fs.writeEOPart(int64(FileIndexByteIndex + n))
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}

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

func (fs *FileSystem) parseFileIndex() error {
	buf := make([]byte, HeaderByteSize)

	n, err := fs.file.ReadAt(buf, HeaderBlockIndex)
	if err != nil {
		return err
	}
	if n != HeaderByteSize {
		return ErrDataBlockMismatch
	}
	// read header
	if string(buf[:len(FileSystemIdentifier)]) != FileSystemIdentifier {
		return ErrArchiverIdentifier

		//// read backup header
		//fs.log.Warnv("First file header is corrupted", "byte array", buf)
		//
		//n, err := fs.file.ReadAt(buf, fs.size-HeaderByteSize)
		//if err != nil {
		//	return err
		//}
		//if n != HeaderByteSize {
		//	return ErrDataBlockMismatch
		//}
		//if string(buf[:len(FileSystemIdentifier)]) != FileSystemIdentifier {
		//	return ErrArchiverIdentifier
		//}
	}

	// ToDO:make compatible for multiple versions
	fs.version = binary.BigEndian.Uint32(buf[8:12])
	fs.blockSize = uint32(binary.BigEndian.Uint64(buf[12:20]))
	fs.blocks = uint32(binary.BigEndian.Uint64(buf[20:28]))
	fs.lastWrittenBlock = uint32(binary.BigEndian.Uint64(buf[28:36]))
	fs.fileIndexSize = uint32(binary.BigEndian.Uint64(buf[36:44]))

	return nil
}
