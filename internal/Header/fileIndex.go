package Header_

import (
	"fmt"

	"github.com/fanap-infra/fsEngine/pkg/fileIndex"
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
	if len(fi) > FileIndexMaxByteSize {
		return fmt.Errorf("fileIndex size %v is too large, Max valid size: %v",
			len(fi), FileIndexMaxByteSize)
	}
	hfs.fileIndexSize = uint32(len(fi))
	//checkSum := crc32.ChecksumIEEE(fi)
	//if hfs.fiChecksum == checkSum {
	//	return nil
	//}
	//hfs.fiChecksum = checkSum

	if hfs.fileIndexSize == 0 {
		hfs.log.Warn("file indexes size is zero")
		// return fmt.Errorf("fileIndex size %v is Zero",
		//	hfs.fileIndexSize)
	}

	// n, err := hfs.file.WriteAt(fi, FileIndexByteIndex)
	n, err := hfs.writeAt(fi, HeaderBlockIndex)
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

	if hfs.fileIndexSize == 0 {
		hfs.log.Warnv("file index size is zero", "fileIndexSize", hfs.fileIndexSize)
		return nil
	}
	n, err := hfs.readAt(buf, FileIndexByteIndex)
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

func (hfs *HFileSystem) UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32,
	fileSize uint32, bam []byte, info []byte) error {
	return hfs.fileIndex.UpdateFileIndexes(fileID, firstBlock, lastBlock, fileSize, bam, info)
}

func (hfs *HFileSystem) FindOldestFile() (*fileIndex.File, error) {
	return hfs.fileIndex.FindOldestFile()
}

func (hfs *HFileSystem) UpdateFileOptionalData(fileId uint32, info []byte) error {
	return hfs.fileIndex.UpdateFileOptionalData(fileId, info)
}

func (hfs *HFileSystem) GetFilesList() []*fileIndex.File {
	return hfs.fileIndex.GetFilesList()
}

//func (hfs *HFileSystem) GetFileOptionalData(fileId uint32) ([]byte, error) {
//	return hfs.fileIndex.
//}
