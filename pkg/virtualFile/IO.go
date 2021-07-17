package virtualFile

import (
	"errors"
	"fmt"

	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"
)

// returns int bytes of written data.
func (v *VirtualFile) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, errors.New("data cannot be empty")
	}

	v.bufTX = append(v.bufTX, data...)
	if uint32(len(v.bufTX)) > v.blockSize {
		m, err := v.fs.Write(v.bufTX[0:len(v.bufTX)-(len(v.bufTX)%int(v.blockSize))], v.id)
		if err != nil {
			return 0, err
		}
		v.fileSize = v.fileSize + uint32(len(v.bufTX[0:len(v.bufTX)-(len(v.bufTX)%int(v.blockSize))]))
		if m != len(v.bufTX)-(len(v.bufTX)%int(v.blockSize)) {
			v.log.Errorv("did not write data completely",
				"data size", len(v.bufTX)-(len(v.bufTX)%int(v.blockSize)), "written size", m)
		}
		v.bufTX = v.bufTX[len(v.bufTX)-(len(v.bufTX)%int(v.blockSize)):]
	}
	return len(data), nil
}

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
	for {
		if v.seekPointer >= v.bufEnd {
			// v.blockAllocationMap.ToArray() we refresh blocks, when simultaneously reading and writing
			blocks := v.blockAllocationMap.ToArray()
			if v.nextBlockIndex >= uint32(len(blocks)) {
				return counter, errors.New("end of file")
			}
			_, err := v.readBlock(blocks[v.nextBlockIndex])
			if err != nil {
				return 0, err
			}
			v.nextBlockIndex = v.nextBlockIndex + 1
		}

		if v.bufEnd-v.seekPointer >= len(data)-counter {
			copy(data[counter:], v.bufRX[v.seekPointer-v.bufStart:v.seekPointer-v.bufStart+len(data)-counter])
			v.seekPointer = v.seekPointer + len(data) - counter
			counter = len(data)
		} else {
			copy(data[counter:counter+v.bufEnd-v.seekPointer], v.bufRX[v.seekPointer-v.bufStart:v.bufEnd-v.bufStart])
			counter = counter + v.bufEnd - v.seekPointer
			v.seekPointer = v.bufEnd
		}

		if counter >= n {
			if counter != n {
				v.log.Warnv("data read more than buffer size", "counter", counter, "n", n)
			}
			return counter, nil
		}
	}
}

func (v *VirtualFile) readBlock(blockIndex uint32) (int, error) {
	buf, err := v.fs.ReadBlock(blockIndex)
	if err != nil {
		return 0, err
	}
	v.bufRX = append(v.bufRX, buf...)
	if len(v.bufRX) > v.bufferSize {
		v.bufRX = v.bufRX[len(v.bufRX)-v.bufferSize:]
	}
	// v.log.Infov("read block", "len(v.bufRX)", len(v.bufRX),
	//	"len(buf)", len(buf), "blockIndex", blockIndex)
	v.bufEnd = v.bufEnd + len(buf)
	v.bufStart = v.bufEnd - len(v.bufRX)

	return len(buf), nil
}

func (v *VirtualFile) ReadAt(data []byte, off int64) (int, error) {
	blocks := v.blockAllocationMap.ToArray()
	maxSize := int64(len(blocks) * int(v.blockSize))
	if off >= maxSize {
		return 0, fmt.Errorf("offset is more than size of file")
	}
	blockIndex := uint32(off * int64(len(blocks)) / maxSize)
	v.bufRX = v.bufRX[:0]
	// v.log.Infov("read at ", "blockIndex",blockIndex, "v.bufStart",v.bufStart,
	//	"bufEnd",v.bufEnd, "len(v.bufRX)",len(v.bufRX), "v.seekPointer", v.seekPointer, "off", off,
	//	"len(data)",len(data))
	v.bufStart = int(blockIndex * v.blockSize)
	v.bufEnd = int(blockIndex * v.blockSize)
	_, err := v.readBlock(blockIndex)
	if err != nil {
		return 0, err
	}
	v.seekPointer = int(off)

	return v.Read(data)
}

// Close
func (v *VirtualFile) Close() error {
	if uint32(len(v.bufTX)) > 0 {
		_, err := v.fs.Write(v.bufTX, v.id)
		if err != nil {
			v.log.Errorv("can not write to file", "err", err.Error())
		}
		v.fileSize = v.fileSize + uint32(len(v.bufTX))
	}
	v.bufTX = v.bufTX[:0]
	v.bufRX = v.bufRX[:0]
	data, err := blockAllocationMap.Marshal(v.blockAllocationMap)
	if err != nil {
		v.log.Errorv("can not marshal bam", "err", err.Error())
	}
	err = v.fs.UpdateFileIndexes(v.id, v.firstBlockIndex, v.lastBlock, v.fileSize)
	if err != nil {
		v.log.Errorv("can not update file indexes", "err", err.Error())
	}
	err = v.fs.BAMUpdated(v.id, data)
	if err != nil {
		v.log.Errorv("can not update bam", "err", err.Error())
	}
	return v.fs.Closed(v.id)
}

func (v *VirtualFile) UpdateFileOptionalData(fileId uint32, info []byte) error {
	return v.fs.UpdateFileOptionalData(fileId, info)
}
