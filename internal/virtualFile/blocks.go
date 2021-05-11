package virtualFile

func (v *VirtualFile) AddBlock(blockIndex uint32) error {
	return v.blockAllocationMap.SetBlockAsAllocated(blockIndex)
}
