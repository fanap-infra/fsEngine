package fsEngine

import (
	"fmt"

	"github.com/fanap-infra/fsEngine/internal/constants"
)

func (fse *FSEngine) writeInBlock(data []byte, blockIndex uint32) (n int, err error) {
	// fse.log.Infov("FSEngine write in block", "blockIndex", blockIndex,
	//	"maxNumberOfBlocks", fse.maxNumberOfBlocks, "len(data)", len(data))
	if blockIndex >= fse.maxNumberOfBlocks {
		return 0, constants.ErrBlockIndexOutOFRange
	}

	n, err = fse.file.WriteAt(data, int64(blockIndex)*int64(fse.blockSize))
	if err != nil {
		fse.log.Errorv("Error Writing to file", "err", err.Error(),
			"file", fse.file.Name(), "blockIndex", blockIndex)
	}

	return
}

func (fse *FSEngine) ReadBlock(blockIndex uint32) ([]byte, error) {
	// fse.log.Infov("FSEngine read in block", "blockIndex", blockIndex)
	if blockIndex >= fse.maxNumberOfBlocks {
		return nil, constants.ErrBlockIndexOutOFRange
	}

	var err error
	buf := make([]byte, fse.blockSize)
	n, err := fse.file.ReadAt(buf, int64(blockIndex*fse.blockSize))
	if err != nil {
		return nil, err
	}
	if n != int(fse.blockSize) {
		return buf, constants.ErrDataBlockMismatch
	}
	data, err := fse.parseBlock(buf)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fse *FSEngine) ReadAt(data []byte, off int64, fileID uint32) (int, error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	// ToDo: implement it
	return 0, nil
}

func (fse *FSEngine) Read(data []byte, fileID uint32) (int, error) {
	return 0, fmt.Errorf("please impkement me")
}

func (fse *FSEngine) WriteAt(b []byte, off int64, fileID uint32) (n int, err error) {
	// ToDo: complete it
	n, err = fse.file.WriteAt(b, off)

	return
}

func (fse *FSEngine) Write(data []byte, fileID uint32) (int, error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	dataSize := len(data)
	if dataSize == 0 {
		return 0, fmt.Errorf("data siz is zero, file ID: %v ", fileID)
	}
	vfInfo, ok := fse.openFiles[fileID]
	if !ok {
		return 0, fmt.Errorf("this file ID: %v did not opened", fileID)
	}
	n := 0
	var err error
	for {
		if n >= dataSize {
			return n, nil
		}
		previousBlock := vfInfo.vfs[0].GetLastBlock()
		blockID := fse.header.FindNextFreeBlockAndAllocate()
		var d []byte
		if dataSize >= n+int(fse.blockSizeUsable) {
			d, err = fse.prepareBlock(data[n:n+int(fse.blockSizeUsable)], fileID, previousBlock, blockID)
			if err != nil {
				return 0, err
			}
		} else {
			d, err = fse.prepareBlock(data[n:], fileID, previousBlock, blockID)
			if err != nil {
				return 0, err
			}
		}

		m, err := fse.writeInBlock(d, blockID)
		if err != nil {
			return 0, err
		}

		err = vfInfo.vfs[0].AddBlockID(blockID)
		if err != nil {
			return 0, err
		}

		err = fse.header.SetBlockAsAllocated(blockID)
		if err != nil {
			return 0, err
		}

		if m != len(d) {
			return 0, fmt.Errorf("block with size: %v did not write correctly, n = %v", m, len(d))
		}
		n = m - constants.BlockHeaderSize + n
	}
}

// It is event handler
func (fse *FSEngine) Closed(fileID uint32) error {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	err := fse.header.UpdateFSHeader()
	if err != nil {
		fse.log.Warnv("Can not updateHeader", "err", err.Error())
	}

	vfInfo, ok := fse.openFiles[fileID]
	if !ok {
		return fmt.Errorf("this file ID: %v did not opened", fileID)
	}
	vfInfo.numberOfOpened = vfInfo.numberOfOpened - 1
	if vfInfo.numberOfOpened == 0 {
		delete(fse.openFiles, fileID)
	}

	err = fse.file.Sync()
	if err != nil {
		fse.log.Warnv("Can not sync file", "err", err.Error())
	}

	return nil
}
