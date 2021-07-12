package virtualFile

func (v *VirtualFile) AddBlockID(blockIndex uint32) error {
	if len(v.blockAllocationMap.ToArray()) == 0 {
		v.firstBlockIndex = blockIndex
		err := v.fs.UpdateFileIndexes(v.id, v.firstBlockIndex, v.lastBlock, v.fileSize)
		if err != nil {
			return err
		}
	}
	err := v.blockAllocationMap.SetBlockAsAllocated(blockIndex)
	if err == nil {
		v.lastBlock = blockIndex
	}
	return err
}

func (v *VirtualFile) GetLastBlock() uint32 {
	return v.lastBlock
}
