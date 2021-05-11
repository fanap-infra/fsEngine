package Header_

import "encoding/binary"

func (hfs *HFileSystem) writeAt(b []byte, off int64) (n int, err error) {
	n, err = hfs.file.WriteAt(b, off)

	return
}

// Close ...
func (hfs *HFileSystem) Close() error {
	err := hfs.updateHeader()
	if err != nil {
		hfs.log.Warnv("Can not updateHeader", "err", err.Error())
		// ToDo: remove it
		return err
	}
	// ToDo:update file system
	err = hfs.file.Sync()
	if err != nil {
		hfs.log.Warnv("Can not sync file", "err", err.Error())
		// ToDo: remove it
		return err
	}
	return hfs.file.Close()
}

func (hfs *HFileSystem) writeEOPart(off int64) (n int, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, EOPart)
	return hfs.file.WriteAt(buf, off)
}
