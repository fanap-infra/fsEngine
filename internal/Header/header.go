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
	tmp64 := make([]byte, 8)

	// byte Identifier
	header = append(header, []byte(FileSystemIdentifier)...)

	// byte Version
	binary.BigEndian.PutUint32(tmp32, hfs.version)
	header = append(header, tmp32...)

	// blocksize, corresponds to BLOCKSIZE
	binary.BigEndian.PutUint64(tmp64, uint64(hfs.blockSize))
	header = append(header, tmp64...)

	// max number of blocks
	binary.BigEndian.PutUint64(tmp64, uint64(hfs.maxNumberOfBlocks))
	header = append(header, tmp64...)

	// last written block
	binary.BigEndian.PutUint64(tmp64, uint64(hfs.lastWrittenBlock))
	header = append(header, tmp64...)

	// file index size
	binary.BigEndian.PutUint64(tmp64, uint64(hfs.fileIndexSize))
	header = append(header, tmp64...)

	//// *** why add this line ???
	//dataTmp := make([]byte, fs.blockSize-uint32(len(header)))
	//header = append(header, dataTmp...)
	return
}

func (hfs *HFileSystem) updateHeader() error {
	header := hfs.generateHeader()
	headerSize := len(header)
	dataTmp := make([]byte, hfs.blockSize-uint32(headerSize))
	dataTmp = append(header, dataTmp...)

	n, err := hfs.writeAt(dataTmp, HeaderBlockIndex)
	if err != nil {
		return err
	}
	if n != len(dataTmp) {
		return fmt.Errorf("Header did not write complete, header size: %v, written size: %v", len(dataTmp), n)
	}
	// ToDo:Maybe it does not be necessary
	err = hfs.file.Sync()
	if err != nil {
		return err
	}

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
	hfs.blockSize = uint32(binary.BigEndian.Uint64(buf[12:20]))
	hfs.maxNumberOfBlocks = uint32(binary.BigEndian.Uint64(buf[20:28]))
	hfs.lastWrittenBlock = uint32(binary.BigEndian.Uint64(buf[28:36]))
	hfs.fileIndexSize = uint32(binary.BigEndian.Uint64(buf[36:44]))

	return nil
}
