package fsEngine

import (
	"fmt"
)

func (fse *FSEngine) writeInBlock(data []byte, blockIndex uint32) (n int, err error) {
	if blockIndex >= fse.maxNumberOfBlocks {
		return 0, ErrBlockIndexOutOFRange
	}

	n, err = fse.file.WriteAt(data, int64(blockIndex)*int64(fse.blockSize))
	if err != nil {
		fse.log.Infov("Error Writing to file", "err", err.Error(), "file", fse.file.Name())
	}

	return
}

func (fse *FSEngine) ReadBlock(blockIndex uint32) ([]byte, error) {
	if blockIndex >= fse.maxNumberOfBlocks {
		return nil, ErrBlockIndexOutOFRange
	}

	var err error
	buf := make([]byte, fse.blockSize)
	n, err := fse.file.ReadAt(buf, int64(blockIndex*fse.blockSize))
	if err != nil {
		return nil, err
	}
	if n != int(fse.blockSize) {
		return buf, ErrDataBlockMismatch
	}
	data, err := fse.parseBlock(buf)
	return data, nil
}

func (fse *FSEngine) ReadAt(data []byte, off int64, fileID uint32) (int, error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	// ToDo: implement it
	return 0, nil
}

func (fse *FSEngine) Read(data []byte, fileID uint32) (int, error) {
	//fse.rIBlockMux.Lock()
	//defer fse.rIBlockMux.Unlock()
	//vf, ok := fse.openFiles[fileID]
	//if !ok {
	//	return 0, fmt.Errorf("this file ID: %v did not opened", fileID)
	//}
	return 0, fmt.Errorf("please impkement me")
}

func (fse *FSEngine) WriteAt(b []byte, off int64, fileID uint32) (n int, err error) {
	// ToDo: complete it
	n, err = fse.file.WriteAt(b, off)

	//if arc.LastFiletime.IsZero() && off >= int64(arc.conf.DataStartBlock) {
	//	arc.LastFiletime = time.Now()
	//}
	return
}

func (fse *FSEngine) Write(data []byte, fileID uint32) (int, error) {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	vf, ok := fse.openFiles[fileID]
	if !ok {
		return 0, fmt.Errorf("this file ID: %v did not opened", fileID)
	}
	n := 0
	blocksNum := uint32(len(data) / BLOCKSIZEUSABLE)
	for i := uint32(0); i < blocksNum; i++ {
		previousBlock := vf.GetLastBlock()
		//blockID := fse.blockAllocationMap.FindNextFreeBlockAndAllocate()
		blockID := fse.header.FindNextFreeBlockAndAllocate()

		d, err := fse.prepareBlock(data, fileID, previousBlock, blockID)
		if err != nil {
			return 0, err
		}
		c, err := fse.writeInBlock(d, blockID)
		if err != nil {
			return 0, err
		}

		err = vf.AddBlockID(blockID)
		if err != nil {
			return 0, err
		}

		err = fse.header.SetBlockAsAllocated(blockID)
		if err != nil {
			return 0, err
		}

		if c != len(d) {
			return 0, fmt.Errorf("block with size: %v did not write correctly, n = %v", c, len(d))
		}
		n = c + n
	}

	if len(data) != int(blocksNum*BLOCKSIZEUSABLE) {
		previousBlock := vf.GetLastBlock()
		blockID := fse.header.FindNextFreeBlockAndAllocate()
		//blockID := fse.blockAllocationMap.FindNextFreeBlockAndAllocate()
		err := vf.AddBlockID(blockID)
		if err != nil {
			return 0, err
		}
		d, err := fse.prepareBlock(data, fileID, previousBlock, blockID)
		if err != nil {
			return 0, err
		}
		c, err := fse.writeInBlock(d, blockID)
		if err != nil {
			return 0, err
		}
		if c != len(d) {
			return 0, fmt.Errorf("block with size: %v did not write correctly, n = %v", c, len(d))
		}
		n = c + n
	}

	return 0, nil
}

func (fse *FSEngine) Closed(fileID uint32) error {
	fse.rIBlockMux.Lock()
	defer fse.rIBlockMux.Unlock()
	// ToDo: check update file index
	delete(fse.openFiles, fileID)
	return nil
}
