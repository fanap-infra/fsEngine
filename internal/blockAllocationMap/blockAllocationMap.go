package blockAllocationMap

import (
	"fmt"

	"github.com/fanap-infra/log"

	"github.com/RoaringBitmap/roaring"
)

type BlockAllocationMap struct {
	rMap             *roaring.Bitmap // BAM data in memory coded with roaring, to be synced later on to Disk.
	LastWrittenBlock uint32
	size             uint32
	// numberOfAllocated uint32
	trigger event
	log     *log.Logger
}

func (blm *BlockAllocationMap) SetBlockAsAllocated(blockIndex uint32) error {
	if blm.IsBlockAllocated(blockIndex) {
		return fmt.Errorf("Block number %v is allocated before", blockIndex)
	}
	blm.LastWrittenBlock = blockIndex
	blm.rMap.Add(blockIndex)
	return nil
}

func (blm *BlockAllocationMap) UnsetBlockAsAllocated(blockIndex uint32) {
	blm.rMap.Remove(blockIndex)
}

func (blm *BlockAllocationMap) IsBlockAllocated(blockIndex uint32) bool {
	return blm.rMap.Contains(blockIndex)
}

func (blm *BlockAllocationMap) FindNextFreeBlockAndAllocate() uint32 {
	//arc.WMux.Lock()
	//defer arc.WMux.Unlock()
	//if arc.lastWrittenBlock < DataStartBlock {
	//	arc.lastWrittenBlock = DataStartBlock - 1
	//}
	alloc := true
	freeIndex := blm.LastWrittenBlock
	for alloc {
		// iterate
		freeIndex += 1
		// ToDo: check full space by size
		if freeIndex == blm.LastWrittenBlock {
			blm.log.Warnv("There is no space", "freeIndex", freeIndex, "LastWrittenBlock", blm.LastWrittenBlock)
			freeIndex = blm.trigger.NoSpace()
			blm.UnsetBlockAsAllocated(freeIndex)
			return freeIndex
		}
		if freeIndex == blm.size {
			freeIndex = 0
			continue
		}
		alloc = blm.IsBlockAllocated(freeIndex)
	}

	return freeIndex
}
