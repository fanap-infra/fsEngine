package virtualFile

import (
	"errors"
	"fmt"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
	errPackage "github.com/fanap-infra/fsEngine/pkg/errstring"
)

const EndOfFile = errPackage.Error("end of file")

// returns int bytes of written data.
func (v *VirtualFile) Write(data []byte) (int, error) {
	v.WMux.Lock()
	defer v.WMux.Unlock()
	if v.readOnly {
		return 0, errors.New("this virtual file is opened in read mode, so you can not write to it")
	}
	if len(data) == 0 {
		return 0, errors.New("data cannot be empty")
	}

	v.bufTX = append(v.bufTX, data...)
	if uint32(len(v.bufTX)) > v.blockSize {
		m, blocksId, err := v.fs.Write(v.bufTX[0:len(v.bufTX)-(len(v.bufTX)%int(v.blockSize))], v.id, v.GetLastBlock())
		if err != nil {
			return 0, err
		}
		for _, blockId := range blocksId {
			err = v.AddBlockID(blockId)
			if err != nil {
				v.log.Errorv("can not add block to virtual file blocks",
					"blockId", blockId, "err", err.Error())
			}
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
				return counter, EndOfFile
			}
			_, err := v.readBlock(blocks[v.nextBlockIndex])
			if err != nil {
				v.log.Warnv("can not read block", "v.nextBlockIndex", v.nextBlockIndex,
					"blocks[v.nextBlockIndex]", blocks[v.nextBlockIndex], "v.bufStart", v.bufStart,
					"v.seekPointer", v.seekPointer, "v.bufEnd", v.bufEnd, "err", err.Error())
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
	buf, err := v.fs.ReadBlock(blockIndex, v.id)
	if err != nil {
		return 0, err
	}
	v.bufRX = append(v.bufRX, buf...)
	if len(v.bufRX) > v.bufferSize {
		v.bufRX = v.bufRX[len(v.bufRX)-v.bufferSize:]
	}
	// v.log.Infov("VirtualFile read block", "len(v.bufRX)", len(v.bufRX),
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
	if off < 0 {
		return 0, fmt.Errorf("negative offset : %v is not acceptable", off)
	}
	blockIndex := uint32(off * int64(len(blocks)) / maxSize)
	v.bufRX = v.bufRX[:0]
	// v.log.Infov("read at ", "blockIndex",blockIndex, "v.bufStart",v.bufStart,
	//	"bufEnd",v.bufEnd, "len(v.bufRX)",len(v.bufRX), "v.seekPointer", v.seekPointer, "off", off,
	//	"len(data)",len(data))
	v.bufStart = int(blockIndex * v.blockSize)
	v.bufEnd = int(blockIndex * v.blockSize)
	_, err := v.readBlock(blocks[blockIndex])
	if err != nil {
		v.log.Warnv("can not read block", "v.nextBlockIndex", v.nextBlockIndex,
			"blocks[v.nextBlockIndex]", blocks[v.nextBlockIndex], "v.bufStart", v.bufStart,
			"v.seekPointer", v.seekPointer, "v.bufEnd", v.bufEnd, "err", err.Error())
		return 0, err
	}
	v.seekPointer = int(off)
	v.nextBlockIndex = blockIndex + 1
	return v.Read(data)
}

func (v *VirtualFile) ChangeSeekPointer(off int64) error {
	blocks := v.blockAllocationMap.ToArray()
	maxSize := int64(len(blocks) * int(v.blockSize))
	if off >= maxSize {
		return fmt.Errorf("offset is more than size of file")
	}
	if off < 0 {
		return fmt.Errorf("negative offset : %v is not acceptable", off)
	}
	blockIndex := uint32(off * int64(len(blocks)) / maxSize)
	v.bufRX = v.bufRX[:0]
	// v.log.Infov("read at ", "blockIndex",blockIndex, "v.bufStart",v.bufStart,
	//	"bufEnd",v.bufEnd, "len(v.bufRX)",len(v.bufRX), "v.seekPointer", v.seekPointer, "off", off,
	//	"len(data)",len(data))
	v.bufStart = int(blockIndex * v.blockSize)
	v.bufEnd = int(blockIndex * v.blockSize)
	_, err := v.readBlock(blocks[blockIndex])
	if err != nil {
		v.log.Warnv("can not read block", "v.nextBlockIndex", v.nextBlockIndex,
			"blocks[v.nextBlockIndex]", blocks[v.nextBlockIndex], "v.bufStart", v.bufStart,
			"v.seekPointer", v.seekPointer, "v.bufEnd", v.bufEnd, "err", err.Error())
		v.bufStart = 0
		v.bufEnd = 0
		v.seekPointer = 0
		return err
	}

	v.seekPointer = int(off)
	v.nextBlockIndex = blockIndex + 1
	return nil
}

// Close
func (v *VirtualFile) Close() error {
	if !v.readOnly {
		v.WMux.Lock()
		defer v.WMux.Unlock()
		if uint32(len(v.bufTX)) > 0 {
			_, blocksId, err := v.fs.Write(v.bufTX, v.id, v.GetLastBlock())
			if err != nil {
				v.log.Errorv("can not write to file", "err", err.Error())
			}
			for _, blockId := range blocksId {
				err = v.AddBlockID(blockId)
				if err != nil {
					v.log.Errorv("can not add block to virtual file blocks",
						"blockId", blockId, "err", err.Error())
				}
			}

			v.fileSize = v.fileSize + uint32(len(v.bufTX))
		}
		v.bufTX = v.bufTX[:0]
		data, err := blockAllocationMap.Marshal(v.blockAllocationMap)
		if err != nil {
			v.log.Errorv("can not marshal bam", "err", err.Error())
		}
		err = v.fs.UpdateFileIndexes(v.id, v.firstBlockIndex, v.lastBlock, v.fileSize, data, v.optionalData)
		if err != nil {
			v.log.Errorv("can not update file indexes", "err", err.Error())
		}
	}
	v.bufRX = v.bufRX[:0]

	return v.fs.Closed(v.id)
}

func (v *VirtualFile) UpdateFileOptionalData(info []byte) error {
	v.WMux.Lock()
	defer v.WMux.Unlock()
	if v.readOnly {
		return errors.New("this virtual file is opened in read mode, so you can not update any thing")
	}
	v.optionalData = info
	data, err := blockAllocationMap.Marshal(v.blockAllocationMap)
	if err != nil {
		v.log.Errorv("can not marshal bam", "err", err.Error())
		return err
	}
	return v.fs.UpdateFileIndexes(v.id, v.firstBlockIndex, v.lastBlock, v.fileSize, data, v.optionalData)
}

func (v *VirtualFile) AddFileSize(size uint32) {
	v.fileSize = v.fileSize + size
}
