package fsEngine

import (
	"fmt"

	"github.com/fanap-infra/fsEngine/internal/constants"
)

func (fse *FSEngine) writeInBlock(data []byte, blockIndex uint32) (int, error) {
	// fse.log.Infov("FSEngine write in block", "blockIndex", blockIndex,
	//	"maxNumberOfBlocks", fse.maxNumberOfBlocks, "len(data)", len(data))
	// fse.WMux.Lock()
	// defer fse.WMux.Unlock()
	if blockIndex >= fse.maxNumberOfBlocks {
		return 0, constants.ErrBlockIndexOutOFRange
	}
	if uint32(len(data)) > fse.blockSize {
		return 0, fmt.Errorf("size of data is larger than block size, size of data: %v,"+
			" but size of block is %v", len(data), fse.blockSize)
	}

	n, err := fse.file.WriteAt(data, int64(blockIndex)*int64(fse.blockSize))
	if err != nil {
		fse.log.Errorv("Error Writing to file", "err", err.Error(),
			"file", fse.file.Name(), "blockIndex", blockIndex)
		return n, err
	}

	return n, nil
}

func (fse *FSEngine) ReadBlock(blockIndex uint32, fileID uint32) ([]byte, error) {
	fse.RMux.Lock()
	defer fse.RMux.Unlock()
	// fse.log.Infov("FSEngine read in block", "blockIndex", blockIndex)
	if blockIndex >= fse.maxNumberOfBlocks {
		return nil, constants.ErrBlockIndexOutOFRange
	}

	var err error
	buf := make([]byte, fse.blockSize)
	n, err := fse.file.ReadAt(buf, int64(blockIndex)*int64(fse.blockSize))
	if err != nil {
		return nil, err
	}
	if n != int(fse.blockSize) {
		return buf, constants.ErrDataBlockMismatch
	}
	data, err := fse.parseBlock(buf, blockIndex, fileID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fse *FSEngine) ReadAt(data []byte, off int64, fileID uint32) (int, error) {
	// fse.rIBlockMux.Lock()
	// defer fse.rIBlockMux.Unlock()
	// ToDo: implement it
	return 0, nil
}

func (fse *FSEngine) Read(data []byte, fileID uint32) (int, error) {
	return 0, fmt.Errorf("please impkement me")
}

//func (fse *FSEngine) WriteAt(b []byte, off int64, fileID uint32) (n int, err error) {
//	// ToDo: complete it
//	n, err = fse.file.WriteAt(b, off)
//
//	return
//}

func (fse *FSEngine) Write(data []byte, fileID uint32, previousBlock uint32) (int, []uint32, error) {
	fse.WMux.Lock()
	defer fse.WMux.Unlock()
	dataSize := len(data)
	if dataSize == 0 {
		return 0, []uint32{}, fmt.Errorf("data siz is zero, file ID: %v ", fileID)
	}
	//vfInfo, ok := fse.openFiles[fileID]
	//if !ok {
	//	return 0, 0, fmt.Errorf("this file ID: %v did not opened", fileID)
	//}
	n := 0
	var err error
	var blocksID []uint32
	for {
		if n >= dataSize {
			if n == dataSize {
				return n, blocksID, nil
			}
			fse.log.Errorv("data is written more than dataSize", "dataSize", dataSize, "n", n)
			return n, []uint32{}, fmt.Errorf("it is wirtten more, dataSize: %v, n = %v", dataSize, n)
		}
		// previousBlock := vfInfo.vfs[0].GetLastBlock()

		blockID := fse.header.FindNextFreeBlockAndAllocate()
		var d []byte
		if dataSize >= n+int(fse.blockSizeUsable) {
			d, err = fse.prepareBlock(data[n:n+int(fse.blockSizeUsable)], fileID, previousBlock, blockID)
			if err != nil {
				return 0, []uint32{}, err
			}
		} else {
			d, err = fse.prepareBlock(data[n:], fileID, previousBlock, blockID)
			if err != nil {
				return 0, []uint32{}, err
			}
		}

		m, err := fse.writeInBlock(d, blockID)
		if err != nil {
			return 0, []uint32{}, err
		}

		//err = vfInfo.vfs[0].AddBlockID(blockID)
		//if err != nil {
		//	fse.log.Errorv("can not add block to virtual file", "fileID", fileID,
		//		"blockID", blockID, "err", err.Error())
		//	return 0, 0, err
		//}

		err = fse.header.SetBlockAsAllocated(blockID)
		if err != nil {
			fse.log.Errorv("can not set block id in header",
				"blockID", blockID, "fileID", fileID)
			return 0, []uint32{}, err
		}
		blocksID = append(blocksID, blockID)
		if m != len(d) {
			return 0, blocksID, fmt.Errorf("block with size: %v did not write correctly, n = %v", m, len(d))
		}
		n = n + m - constants.BlockHeaderSize
	}
}

// It is event handler
func (fse *FSEngine) Closed(fileID uint32) error {
	fse.WMux.Lock()
	defer fse.WMux.Unlock()
	fse.RMux.Lock()
	defer fse.RMux.Unlock()
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	err := fse.header.UpdateFSHeader()
	if err != nil {
		fse.log.Warnv("Can not updateHeader", "err", err.Error())
	}

	vfInfo, ok := fse.openFiles[fileID]
	if !ok {
		return fmt.Errorf("this file ID: %v did not opened", fileID)
	}
	vfInfo.numberOfOpened = vfInfo.numberOfOpened - 1
	fse.log.Infov("file closed", "fileID", fileID, "numberOfOpened", vfInfo.numberOfOpened)
	if vfInfo.numberOfOpened == 0 {
		delete(fse.openFiles, fileID)
	}

	err = fse.file.Sync()
	if err != nil {
		fse.log.Warnv("Can not sync file", "err", err.Error())
	}

	return nil
}
