package Header_

import "encoding/binary"

func (fs *FileSystem) writeAt(b []byte, off int64) (n int, err error) {
	n, err = fs.file.WriteAt(b, off)

	return
}

// Close ...
func (fs *FileSystem) Close() error {
	err := fs.updateHeader()
	if err != nil {
		fs.log.Warnv("Can not updateHeader", "err", err.Error())
		// ToDo: remove it
		return err
	}
	// ToDo:update file system
	err = fs.file.Sync()
	if err != nil {
		fs.log.Warnv("Can not sync file", "err", err.Error())
		// ToDo: remove it
		return err
	}
	return fs.file.Close()
}

func (fs *FileSystem) writeEOPart(off int64) (n int, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, EOPart)
	return fs.file.WriteAt(buf, off)
}
