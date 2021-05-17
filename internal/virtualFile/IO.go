package virtualFile

import "errors"

// Write
// returns int bytes of written data.
func (v *VirtualFile) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, errors.New("data cannot be empty")
	}

	v.vfBuf = append(v.vfBuf, data...)
	if uint32(len(v.vfBuf)) > v.blockSize {
		return v.fs.Write(v.vfBuf, v.id)
	}
	return 0, nil
}

// Write
// returns int bytes of written data.
func (v *VirtualFile) WriteAt(data []byte, off uint32) (int, error) {
	// ToDo: implement this
	return 0, nil
}

// Read
func (v *VirtualFile) Read(data []byte) (int, error) {
	// panic("implement me")
	return v.fs.Read(data, v.id)
}

// ReadAt
func (v *VirtualFile) ReadAt(data []byte, off int64) (int, error) {
	return v.fs.ReadAt(data, off, v.id)
}

// Close
func (v *VirtualFile) Close() error {
	//if !v.readOnly {
	//	b, _ := generateFrameChunk(v.frameChunk)
	//	_, _ = v.Write(b)
	//}
	//v.Closed = true
	//if !v.readOnly {
	//	err := v.Fsync(true)
	//	if err != nil {
	//		v.log.Errorv("FSYNC", "err", err.Error())
	//	}
	//}
	//
	//err := v.Arc.UpdateFileIndex()
	//if err != nil {
	//	v.log.Errorv("updateFileIndex", "err", err.Error())
	//}
	//err = v.Arc.SyncBamToDisk()
	//if err != nil {
	//	v.log.Errorv("syncBamToDisk", "err", err.Error())
	//}
	if uint32(len(v.vfBuf)) > 0 {
		_, err := v.fs.Write(v.vfBuf, v.id)
		v.log.Errorv("Can not write to file", "err", err.Error())
	}
	return v.fs.Closed(v.id)
	//return nil
}
