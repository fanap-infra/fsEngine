package virtualFile

func (v *VirtualFile) AddBlockID(blockIndex uint32) error {
	err := v.blockAllocationMap.SetBlockAsAllocated(blockIndex)
	if err == nil {
		v.lastBlock = blockIndex
	}
	return err
}

func (v *VirtualFile) GetLastBlock() uint32 {
	return v.lastBlock
}
