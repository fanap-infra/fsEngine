package Header_

import (
	"fmt"
	"hash/crc32"

	"github.com/fanap-infra/fsEngine/internal/fileIndex"
)

func (hfs *HFileSystem) generateFileIndex() ([]byte, error) {
	bin, err := hfs.fileIndex.GenerateBinary()
	if err != nil {
		return nil, err
	}
	//pw := bytes.Buffer{}
	//
	//zw, err := gzip.NewWriterLevel(&pw, 2)
	//if err != nil {
	//	return nil, err
	//}
	//_, _ = zw.Write(bin)
	//err = zw.Close()
	//if err != nil {
	//	return nil, err
	//}
	//bin = pw.Bytes()
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

	//in := bytes.NewReader(buf)
	//
	//gz, err := gzip.NewReader(in)
	//if err != nil {
	//	return err
	//}
	//b := new(bytes.Buffer)
	//_, err = io.Copy(b, gz)
	//if err != nil {
	//	return err
	//}
	//
	//binary := b.Bytes()
	//err = gz.Close()
	//
	//if err != nil {
	//	return err
	//}

	err = hfs.fileIndex.InitFromBinary(buf)
	if err != nil {
		return err
	}

	return nil
}

// ToDo:update blm binaries
func (hfs *HFileSystem) UpdateBAM(fileID uint32, data []byte) error {
	return hfs.fileIndex.UpdateBAM(fileID, data)
}

func (hfs *HFileSystem) UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32, fileSize uint32) error {
	return hfs.fileIndex.UpdateFileIndexes(fileID, firstBlock, lastBlock, fileSize)
}

func (hfs *HFileSystem) FindOldestFile() (*fileIndex.File, error) {
	return hfs.fileIndex.FindOldestFile()
}
