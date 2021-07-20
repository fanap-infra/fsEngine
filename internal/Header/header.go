package Header_

import (
	"encoding/binary"
	"fmt"
)

// generate Header
// 	+------------+---------+-----------+------------+------------+---------------+
//	| Identifier | Version | BLOCKSIZE | BLOCKCOUNT | LastWBlock | FileIndexSize |
//	+============+=========+===========+============+============+===============+
//	|   8 BYTES  | 4 BYTES |  8 BYTES  |   8 BYTES  |   8 BYTES  |     8 BYTES   |
//	+------------+---------+-----------+------------+------------+---------------+
//	|    CONST   |  uint32 |   uint64  |   uint64   |   uint64   |     uint64    |
//	+------------+---------+-----------+------------+------------+---------------+
//
func (hfs *HFileSystem) generateHeader() (header []byte) {
	header = make([]byte, 0)
	tmp32 := make([]byte, 4)

	// byte Identifier
	header = append(header, []byte(FileSystemIdentifier)...)

	// byte Version
	binary.BigEndian.PutUint32(tmp32, hfs.version)
	header = append(header, tmp32...)

	// blocksize, corresponds to BLOCKSIZE
	binary.BigEndian.PutUint32(tmp32, hfs.blockSize)
	header = append(header, tmp32...)

	// max number of blocks
	binary.BigEndian.PutUint32(tmp32, hfs.maxNumberOfBlocks)
	header = append(header, tmp32...)

	// last written block
	binary.BigEndian.PutUint32(tmp32, hfs.lastWrittenBlock)
	header = append(header, tmp32...)

	// file index size
	binary.BigEndian.PutUint32(tmp32, hfs.fileIndexSize)
	header = append(header, tmp32...)

	// blm size
	binary.BigEndian.PutUint32(tmp32, hfs.blmSize)
	header = append(header, tmp32...)

	//// *** why add this line ???
	//dataTmp := make([]byte, fs.blockSize-uint32(len(header)))
	//header = append(header, dataTmp...)
	return
}

func (hfs *HFileSystem) updateHeader() error {
	header := hfs.generateHeader()
	headerSize := len(header)
	// dataTmp := make([]byte, uint32(headerSize))
	// dataTmp = append(header, dataTmp...)

	n, err := hfs.writeAt(header, HeaderBlockIndex)
	if err != nil {
		return err
	}
	if n != headerSize {
		return fmt.Errorf("header did not write complete, header size: %v, written size: %v", headerSize, n)
	}
	//// ToDo:Maybe it does not be necessary
	//err = hfs.file.Sync()
	//if err != nil {
	//	return err
	//}

	//// write header back up
	//n, err = fs.WriteAt(header, fs.size-int64(headerSize))
	//if err != nil {
	//	return err
	//}
	//if n != len(header) {
	//	return fmt.Errorf("Header did not write complete, header size: %v, written size: %v", headerSize, n)
	//}

	return nil
}

func (hfs *HFileSystem) parseHeader() error {
	buf := make([]byte, HeaderByteSize)

	n, err := hfs.file.ReadAt(buf, HeaderBlockIndex)
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

	hfs.version = binary.BigEndian.Uint32(buf[8:12])
	hfs.blockSize = binary.BigEndian.Uint32(buf[12:16])
	hfs.maxNumberOfBlocks = binary.BigEndian.Uint32(buf[16:20])
	hfs.lastWrittenBlock = binary.BigEndian.Uint32(buf[20:24])
	hfs.fileIndexSize = binary.BigEndian.Uint32(buf[24:28])
	hfs.blmSize = binary.BigEndian.Uint32(buf[28:32])

	hfs.size = int64(hfs.maxNumberOfBlocks * hfs.blockSize)
	return nil
}
