package Header_

type configs struct {
	BLOCKSIZE                   uint32 // = 1 << 19
	BLOCKSIZEUSABLE             uint32 // BLOCKSIZE - 20 // usable size is block size minus two uint32 sized locations reserved for
	StorageMaxSize              uint64 // 1 << 44
	DataStartBlock              uint32
	FileIndexReservedSize       uint32 // 1 << 22 // 4MB in continuous blocks is reserved for fileIndex
	FileIndexReservedSizeBlocks uint32 //(FileIndexReservedSize / BLOCKSIZEUSABLE) * 2
	// ***
	//	Block Allocation Map is a BitMap of all the blocks in Archiver, hence it is
	//	calculated by dividing maximum storage size by block size divided by 8.
	//	Even when storage size is lower than maximum, BAM is reserved for the maximum
	//	because Archiver size can be increased in code later on.
	// ***
	BlockAllocationMapSize uint32
}

func loadConf(hfs *HFileSystem) {
	hfs.conf.BLOCKSIZE = hfs.blockSize // 1 << 19
	hfs.conf.BLOCKSIZEUSABLE = hfs.conf.BLOCKSIZE - 20
	hfs.conf.StorageMaxSize = 1 << 44
	hfs.conf.BlockAllocationMapSize = uint32(hfs.conf.StorageMaxSize/uint64(hfs.conf.BLOCKSIZEUSABLE)) / 8 // Size in blocks.
	hfs.conf.FileIndexReservedSize = 1 << 22
	hfs.conf.FileIndexReservedSizeBlocks = (hfs.conf.FileIndexReservedSize / hfs.conf.BLOCKSIZEUSABLE) * 2 //(FileIndexReservedSize / BLOCKSIZEUSABLE) * 2
	hfs.conf.DataStartBlock = (hfs.conf.BlockAllocationMapSize / hfs.conf.BLOCKSIZEUSABLE) + hfs.conf.FileIndexReservedSizeBlocks + 1
	hfs.lastWrittenBlock = hfs.conf.DataStartBlock
}

func (hfs *HFileSystem) GetBlockSize() uint32 {
	return hfs.blockSize
}

func (hfs *HFileSystem) GetBlocksNumber() uint32 {
	return hfs.maxNumberOfBlocks
}
