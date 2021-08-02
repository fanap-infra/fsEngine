package constants

import "github.com/fanap-infra/fsEngine/pkg/err"

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
	FsPath               = "fs.beh"
	HeaderPath           = "Header.Beh"
	HeaderBackUpPath     = "BackUpHeader.Beh"
	FileSystemIdentifier = "BehFS;P "
	FileSystemVersion    = 1
	BlockHeaderSize      = 16
	// In this version,we allocate constant space for each part
	HeaderBlockIndex = 0 // BackUp is last block
	HeaderByteSize   = 36

	FileIndexBlockIndex = HeaderBlockIndex + 1 // BackUp Size - 1 - FileIndexBlockSize
	FileIndexBlockSize  = 5

	AllocationMapBlockIndex = FileIndexBlockIndex + FileIndexBlockSize // BackUp Size - 1 - FileIndexBlockSize - AllocationMapBlockSize
	AllocationMapBlockSize  = 5

	// EOFInBlockIndicator    = 0x80000000
	ErrBlockUnallocated     = err.Error("block is unallocated")
	ErrBlockIndexOutOFRange = err.Error("block index is out of range")
	ErrDataBlockMismatch    = err.Error("block data-size mismatch")
	ErrFileExists           = err.Error("file already exist")
	ErrArchiverIdentifier   = err.Error("ArchiverIdentifier is not detected")
	ErrArchiverVersion      = err.Error("Archiver Version is not correct")
	BYTE                    = 1
	KILOBYTE                = BYTE * 1024
	MEGABYTE                = KILOBYTE * 1024

	BLOCKSIZE       = 1 << 19        //
	BLOCKSIZEUSABLE = BLOCKSIZE - 20 // usable size is block size minus two uint32 sized locations reserved for
	StorageMaxSize  = 1 << 44

	VirtualFileBufferBlockNumber = 5
)
