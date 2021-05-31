package fsEngine

import "github.com/fanap-infra/FSEngine/pkg/err"

// Storage Shape
//	+---------+---------------+---------+---------------------+------+----------+---------+------------+------+--------+
//	|  Header  | FileIndex | BlockAllocationMap  |         Blocks        |  BlockAllocationMap | FileIndex |   Header  |
//	+=========+===============+=========+=====================+======+==========+=========+
//	| 1 Block |   5 Blocks    | 4 BYTES | (BLOCKSIZE-20)BYTES | 1BIT |   31BIT  | 4 BYTES |
//	+---------+---------------+---------+---------------------+------+----------+---------+
//	|                                  <=  BLOCKSIZE  =>                                  |
//	+-------------------------------------------------------------------------------------+
//

const (
	FileSystemIdentifier = "BehFS;P "
	FileSystemVersion    = 1
	// In this version,we allocate constant space for each part
	HeaderBlockIndex = 0 // BackUp is last block
	HeaderByteSize   = 36

	FileIndexBlockIndex = HeaderBlockIndex + 1 // BackUp Size - 1 - FileIndexBlockSize
	FileIndexBlockSize  = 5

	AllocationMapBlockIndex = FileIndexBlockIndex + FileIndexBlockSize // BackUp Size - 1 - FileIndexBlockSize - AllocationMapBlockSize
	AllocationMapBlockSize  = 5

	//MaxArchiverTime     = 3600 * time.Second
	//ArchiverCheckerIter = 60 * time.Second
	//// ArchiverIdentifier is added to the beginning of archiver to identify it as such.

	// EOFInBlockIndicator    = 0x80000000
	ErrBlockUnallocated     = err.Error("block is unallocated")
	ErrBlockIndexOutOFRange = err.Error("block index is out of range")
	ErrDataBlockMismatch    = err.Error("block data-size mismatch")
	ErrFileExists           = err.Error("file already exist")
	ErrArchiverIdentifier   = err.Error("ArchiverIdentifier is not detected")
	ErrArchiverVersion      = err.Error("Archiver Version is not correct")
	BYTE                   = 1
	KILOBYTE               = BYTE * 1024
	MEGABYTE               = KILOBYTE * 1024
	GIGABYTE               = MEGABYTE * 1024
	//DefaultArchiverSize    = 20 * GIGABYTE
	//ArchiverMinimumPadding = 250 * MEGABYTE
	//ArchiverMaximumPadding = 5 * GIGABYTE
	//ArchiverPaddingPercent = 0.02 // 2 percent
	//
	//// Default value is 512KB, must be divisible by 4096B (disk block size).
	BLOCKSIZE       = 1 << 19        //
	BLOCKSIZEUSABLE = BLOCKSIZE - 20 // usable size is block size minus two uint32 sized locations reserved for
	StorageMaxSize  = 1 << 44
	//FileHeaderCapacity = 64 // make header size be power of 2
	//// BAMStorageBlockSize is 17TB as it stands currently.
	//BAMStorageBlockSize     = ((1 << 26) / BLOCKSIZEUSABLE) / 8
	//FileIndexReservedSize   = 1 << 22 // 4MB in continuous blocks is reserved for fileIndex
	//DEFAULTSTARTBLOCK       = 1
	//BlockAllocationMapStart = 2 // Starting Block of BAM(Block Allocation Map)
	///*
	//	Block Allocation Map is a BitMap of all the blocks in Archiver, hence it is
	//	calculated by dividing maximum storage size by block size divided by 8.
	//	Even when storage size is lower than maximum, BAM is reserved for the maximum
	//	because Archiver size can be increased in code later on.
	//*/
	//BlockAllocationMapSize = (StorageMaxSize / BLOCKSIZEUSABLE) / 8 // Size in blocks.
	//
	//FileIndexStartBlockFlip     = (BlockAllocationMapSize / BLOCKSIZEUSABLE) + 1
	//FileIndexStartBlockFlop     = ((BlockAllocationMapSize / BLOCKSIZEUSABLE) + 1) + FileIndexReservedSizeBlocks/2
	//FileIndexReservedSizeBlocks = (FileIndexReservedSize / BLOCKSIZEUSABLE) * 2

	// DataStartBlock = (BlockAllocationMapSize / BLOCKSIZEUSABLE) + FileIndexReservedSizeBlocks + 1

	VirtualFileBufferBlockNumber = 5

)
