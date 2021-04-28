package fs

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
