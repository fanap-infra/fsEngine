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

	v.bufTX = append(v.bufTX, data...)
	if uint32(len(v.bufTX)) > v.blockSize {
		_, err := v.fs.Write(v.bufTX[0: len(v.bufTX) - (len(v.bufTX)%int(v.blockSize))], v.id)
		if err != nil {
			return 0, err
		}
		v.bufTX = v.bufTX[len(v.bufTX) - (len(v.bufTX)%int(v.blockSize)):]
		//if err != nil {
		//	return 0, err
		//}
	}
	return len(data), nil
}

// Write
// returns int bytes of written data.
func (v *VirtualFile) WriteAt(data []byte, off uint32) (int, error) {
	// ToDo: implement this
	v.log.Warn("virtual file WriteAt method does not implemented")
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
			if v.blockIndex >= uint32(len(blocks)) {
				return 0, errors.New("end of file")
			}
			_, err := v.ReadBlock(blocks[v.blockIndex])
			if err != nil {
				return 0, err
			}
			v.blockIndex = v.blockIndex + 1

		}

		if v.bufEnd - v.seekPointer >= len(data)-counter {
			copy(data[counter:len(data)],v.bufRX[v.seekPointer-v.bufStart: v.seekPointer-v.bufStart+len(data)-counter])
			v.seekPointer = v.seekPointer+len(data)-counter
			counter =  len(data)
		} else {
			copy(data[counter:counter+v.bufEnd-v.seekPointer],v.bufRX[v.seekPointer-v.bufStart: v.bufEnd-v.bufStart])
			counter = counter + v.bufEnd-v.seekPointer
			v.seekPointer = v.bufEnd

		}

		if counter >= n {
			if counter != n {
				v.log.Warnv("data read more than buffer size", "counter", counter, "n",n)
			}
			return counter, nil
		}
	}
}

func (v *VirtualFile) ReadBlock(blockIndex uint32) (int, error) {
	buf, err := v.fs.ReadBlock(blockIndex)
	if err != nil {
		return 0, err
	}
	v.bufRX = append(v.bufRX, buf...)
	if len(v.bufRX) > v.bufferSize {
		v.bufRX = v.bufRX[len(v.bufRX)-v.bufferSize:]
	}
	v.bufEnd = v.bufEnd + len(buf)
	v.bufStart = v.bufEnd - len(v.bufRX)

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
	if uint32(len(v.bufTX)) > 0 {
		_, err := v.fs.Write(v.bufTX, v.id)
		if err != nil {
			v.log.Errorv("can not write to file", "err", err.Error())
		}
	}

	v.bufTX = v.bufTX[:]
	v.bufRX = v.bufRX[:]

	return v.fs.Closed(v.id)
}
