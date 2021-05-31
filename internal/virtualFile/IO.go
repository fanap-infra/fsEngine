package virtualFile

import (
	"errors"
)

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
	n := len(data)
	if n == 0 {
		return 0, errors.New("data cannot be zero size")
	}
	counter := 0
	for{
		if v.seekPointer >= v.bufEnd {
			//v.blockAllocationMap.ToArray() we refresh blocks, because
			blocks := v.blockAllocationMap.ToArray()
			if v.blockIndex+1 >= uint32(len(blocks)) {
				return 0, errors.New("end of file")
			}
			_, err := v.ReadBlock(v.blockIndex+1)
			if err != nil {
				return 0, err
			}
			v.blockIndex = v.blockIndex+1
		}

		if v.bufEnd - v.seekPointer >= len(data)-counter {
			copy(data[counter:len(data)],v.vfBuf[v.seekPointer: v.seekPointer+len(data)-counter])
			counter = counter + len(data)
		} else {
			copy(data[counter:counter+v.bufEnd-v.seekPointer],v.vfBuf[v.seekPointer: v.seekPointer+len(data)-counter])
			counter = counter + v.bufEnd-v.seekPointer
		}

		if counter >= n {
			return counter, nil
		}
	}
}

func (v *VirtualFile) ReadBlock(blockIndex uint32) (int, error) {
	buf, err := v.fs.ReadBlock(blockIndex)
	if err != nil {
		return 0, err
	}
	v.vfBuf = append(v.vfBuf, buf...)
	if len(v.vfBuf) > v.bufferSize {
		v.vfBuf = v.vfBuf[len(v.vfBuf)-v.bufferSize:]
	}
	v.bufEnd = v.bufEnd + len(buf)
	v.bufStart = v.bufEnd - len(v.vfBuf)
	return len(buf), nil
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
		v.log.Errorv("can not write to file", "err", err.Error())
	}
	return v.fs.Closed(v.id)
	//return nil
}
