package Header_

import (
	"fmt"

	"github.com/fanap-infra/fsEngine/pkg/blockAllocationMap"
)

func (hfs *HFileSystem) updateBLM() error {
	data, err := blockAllocationMap.Marshal(hfs.blockAllocationMap)
	if err != nil {
		return err
	}

	hfs.blmSize = uint32(len(data))
	if hfs.blmSize > BlockAllocationMaxByteSize {
		return fmt.Errorf("blm size %v is too large, Max valid size: %v", hfs.blmSize, BlockAllocationMaxByteSize)
	}

	if hfs.storeInRedis {
		err := hfs.setRedisKeyValue("arch"+fmt.Sprint(hfs.id)+"_BLM", data)
		if err != nil {
			return err
		}
		return nil
	}

	n, err := hfs.writeAt(data, BlockAllocationMapByteIndex)
	if err != nil {
		return err
	}
	if n != len(data) {
		return fmt.Errorf("blm did not write complete, header size: %v, written size: %v", hfs.blmSize, n)
	}
	n, err = hfs.writeEOPart(int64(BlockAllocationMapByteIndex + n))
	if err != nil {
		return err
	}

	if n != 4 {
		return fmt.Errorf("blm did not write complete, header size: %v, written size: %v", hfs.blmSize, n)
	}

	return nil
}

func (hfs *HFileSystem) parseBLM() error {
	var buf []byte
	var err error
	if hfs.storeInRedis {
		buf, err = hfs.getRedisValue("arch" + fmt.Sprint(hfs.id) + "_BLM")
		if err != nil {
			hfs.log.Errorv("can get value from redis", "key", "arch"+fmt.Sprint(hfs.id)+"_BLM",
				"err", err.Error())
			return err
		}
	} else {
		buf = make([]byte, hfs.blmSize)

		n, err := hfs.readAt(buf, BlockAllocationMapByteIndex)
		if err != nil {
			return err
		}

		if n != int(hfs.blmSize) {
			return ErrDataBlockMismatch
		}
	}

	blm, err := blockAllocationMap.Open(hfs.log, hfs.eventHandler, hfs.maxNumberOfBlocks, hfs.lastWrittenBlock, buf)
	if err != nil {
		return err
	}

	hfs.blockAllocationMap = blm
	return nil
}

func (hfs *HFileSystem) FindNextFreeBlockAndAllocate() uint32 {
	return hfs.blockAllocationMap.FindNextFreeBlockAndAllocate()
}

func (hfs *HFileSystem) SetBlockAsAllocated(blockIndex uint32) error {
	return hfs.blockAllocationMap.SetBlockAsAllocated(blockIndex)
}

func (hfs *HFileSystem) UnsetBlockAsAllocated(blockIndex uint32) {
	hfs.blockAllocationMap.UnsetBlockAsAllocated(blockIndex)
}

func (hfs *HFileSystem) GetBLMArray() []uint32 {
	return hfs.blockAllocationMap.ToArray()
}

func (hfs *HFileSystem) IsBlockAllocated(blockIndex uint32) bool {
	return hfs.blockAllocationMap.IsBlockAllocated(blockIndex)
}
