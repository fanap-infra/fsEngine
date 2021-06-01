package Header_

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"hash/crc32"
)

func (hfs *HFileSystem) generateFileIndex() ([]byte, error) {
	bin, err := hfs.fileIndex.GenerateBinary()
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

func (hfs *HFileSystem) updateFileIndex() error {
	fi, err := hfs.generateFileIndex()
	if err != nil {
		return err
	}

	checkSum := crc32.ChecksumIEEE(fi)
	if hfs.fiChecksum == crc32.ChecksumIEEE(fi) {
		return nil
	}
	hfs.fiChecksum = checkSum
	hfs.fileIndexSize = uint32(len(fi))
	if hfs.fileIndexSize > FileIndexMaxByteSize {
		return fmt.Errorf("fileIndex size %v is too large, Max valid size: %v",
			hfs.fileIndexSize, FileIndexMaxByteSize)
	}
	n, err := hfs.file.WriteAt(fi, FileIndexByteIndex)
	if err != nil {
		return err
	}
	if n != len(fi) {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}
	n, err = hfs.writeEOPart(int64(FileIndexByteIndex + n))
	if err != nil {
		return err
	}
	if n != 4 {
		return fmt.Errorf("fileIndex did not write complete, header size: %v, written size: %v", len(fi), n)
	}

	return nil
}

func (hfs *HFileSystem) parseFileIndex() error {
	buf := make([]byte, hfs.fileIndexSize)

	n, err := hfs.file.ReadAt(buf, FileIndexByteIndex)
	if err != nil {
		return err
	}
	if n != int(hfs.fileIndexSize) {
		return ErrDataBlockMismatch
	}

	err = hfs.fileIndex.InitFromBinary(buf)
	if err != nil {
		return err
	}

	return nil
}
