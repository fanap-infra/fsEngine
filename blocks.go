package fsEngine

import (
	"encoding/binary"
	"fmt"
)

func (fse *FSEngine) NoSpace() uint32 {
	fileIndex, err := fse.header.FindOldestFile()
	if err != nil {
		fse.log.Errorv("can not find oldest file", "err", err.Error())
		return 0
	}
	blockIndex := fileIndex.FirstBlock
	err = fse.RemoveVirtualFile(fileIndex.Id)
	if err != nil {
		fse.log.Errorv("can not remove virtual file", "id", fileIndex.Id,
			"err", err.Error())
		return 0
	}
	fse.eventsHandler.VirtualFileDeleted(fileIndex.Id, "file deleted due to space requirements")
	return blockIndex
}

// BlockStructure
//	+===============+===============+===============+=======+
//	|    				 	  BLOCK 				   		|
//	+--------+------+--------+------+--------+------+-------+
//	|  blockID  |   fileID 	 |	 prevBlock	 | valid Size 	|
//	+--------+------+--------+------+--------+------+-------+
//  |  4 byte   |   4 byte   |    4 byte     |   4 byte     |   16 Byte
//	+--------+------+--------+------+--------+------+-------+
// Warning It is not thread safe
func (fse *FSEngine) prepareBlock(data []byte, fileID uint32, previousBlock uint32, blockID uint32) ([]byte, error) {
	dataTmp := make([]byte, 0)

	headerTmp := make([]byte, 4)
	binary.BigEndian.PutUint32(headerTmp, blockID)
	dataTmp = append(dataTmp, headerTmp...)
	binary.BigEndian.PutUint32(headerTmp, fileID)
	dataTmp = append(dataTmp, headerTmp...)
	binary.BigEndian.PutUint32(headerTmp, previousBlock)
	dataTmp = append(dataTmp, headerTmp...)
	binary.BigEndian.PutUint32(headerTmp, uint32(len(data)))
	dataTmp = append(dataTmp, headerTmp...)
	dataTmp = append(dataTmp, data...)

	return dataTmp, nil
}

func (fse *FSEngine) parseBlock(data []byte) ([]byte, error) {
	dataSize := binary.BigEndian.Uint32(data[12:16])
	if dataSize > fse.blockSize-16 {
		return nil, fmt.Errorf("blockd ata size is too large, dataSize: %v", dataSize)
	}

	return data[16 : dataSize+16], nil
}

func (fse *FSEngine) BAMUpdated(fileID uint32, bam []byte) error {
	// ToDo: because of file index,we use mutex
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	return fse.header.UpdateBAM(fileID, bam)
}

func (fse *FSEngine) UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32) error {
	fse.crudMutex.Lock()
	defer fse.crudMutex.Unlock()
	return fse.header.UpdateFileIndexes(fileID, firstBlock, lastBlock)
}

/*
var byteOrder = binary.BigEndian

var errCorruptBlock = errors.New("corrupt block")

// readBlock takes the raw bytes of a block, reads and parses the footer
// and returns a block with the footer indexes.
func readBlock(raw []byte) (block, error) {
	if len(raw) < 4 {
		return block{}, errCorruptBlock
	}

	// the trailing 4 bytes are a uint32 specifying the offset
	// of where the footer begins.
	footerOff := byteOrder.Uint32(raw[len(raw)-4:])
	if int(footerOff) >= len(raw) {
		return block{}, errCorruptBlock
	}

	footer := bytes.NewReader(raw[footerOff : len(raw)-4])
	restarts := make([]uint32, 0)
	for footer.Len() > 0 {
		off, err := binary.ReadUvarint(footer)
		if err != nil {
			return block{}, err
		}
		if int(off) >= len(raw) {
			return block{}, errCorruptBlock
		}
		restarts = append(restarts, uint32(off))
	}

	return block{
		data:     raw[:footerOff],
		restarts: restarts,
	}, nil
}

type block struct {
	data     []byte   // raw block, without footer
	restarts []uint32 // parsed footer
}

// iter iterates over the block calling fn on each key-value pair.
// The key bytes are owned by iter and are only valid until fn
// returns. The value is a slice of the block's bytes. Keeping
// a reference to it will prevent GC-ing of the block.
func (b block) iter(fn func(key, val []byte)) error {
	it := blockIterator{block: b}
	for it.hasNext() {
		err := it.next()
		if err != nil {
			return err
		}

		fn(it.key, it.value)
	}
	return nil
}

// iterAt returns a block iterator starting at the provided
// key, or where the provided key would be if it were contained
// in the block. It binary searches among the block's restart
// points and linear scans from there.
func (b block) iterAt(key []byte) (blockIterator, error) {
	bi := blockIterator{block: b}

	var decodeErr error
	i := sort.Search(len(b.restarts), func(r int) bool {
		bi.off = int(b.restarts[r])
		err := bi.next()
		if err != nil {
			decodeErr = err
			return false
		}

		return bytes.Compare(bi.key, key) >= 0
	})
	if decodeErr != nil {
		return bi, decodeErr
	}

	// i is now the index of the smallest restart point >= key.
	// We now linear scan from the i-1 restart point.
	if i == 0 {
		// Scan from the beginning of the block, not the first
		// restart point which is already several keys into the
		// block.
		bi.off = 0
	} else {
		bi.off = int(b.restarts[i-1])
	}

	for bi.hasNext() {
		err := bi.next()
		if err != nil {
			return bi, err
		}
		if bytes.Compare(bi.key, key) >= 0 {
			break
		}
	}
	return bi, nil
}

type blockIterator struct {
	block block
	off   int

	entryOff   int
	key, value []byte
}

func (bi blockIterator) hasNext() bool {
	return bi.off < len(bi.block.data)
}

func (bi *blockIterator) ReadByte() (byte, error) {
	if bi.off >= len(bi.block.data) {
		return 0, io.ErrUnexpectedEOF
	}
	b := bi.block.data[bi.off]
	bi.off++
	return b, nil
}

func (bi *blockIterator) next() error {
	startOff := bi.off
	shared, err := binary.ReadUvarint(bi)
	if err != nil {
		return err
	}
	nonshared, err := binary.ReadUvarint(bi)
	if err != nil {
		return err
	}
	valueLen, err := binary.ReadUvarint(bi)
	if err != nil {
		return err
	}

	rest := bi.block.data[bi.off:]

	// set bi.key, bi.value for the next value
	keyDelta := rest[:nonshared]
	bi.key = bi.key[:shared]
	bi.key = append(bi.key, keyDelta...)
	bi.value = rest[nonshared : nonshared+valueLen]
	bi.off += int(nonshared + valueLen)
	bi.entryOff = startOff
	return nil
}

// blockBuilder generates blocks with prefix-compressed keys.
//
// Every restartInterval keys blockBuilder will restart the
// prefix compression. It saves the restart offset to restarts.
// Upon finishing the block, the restarts are added to the end.
// A read may binary search through the restart points to find
// where to begin searching.
type blockBuilder struct {
	buf             bytes.Buffer
	lastKey         []byte
	counter         int
	restartInterval int
	restarts        []int
	tmp             [binary.MaxVarintLen64]byte // varint scratch space
}

// Reset resets the builder to be empty,
// but it retains the underlying storage for use by future writes.
func (bb *blockBuilder) reset() {
	bb.buf.Reset()
	bb.lastKey = nil
	bb.counter = 0
	bb.restarts = nil
}

// size estimates the size of the finished block
func (bb *blockBuilder) size() int {
	return bb.buf.Len() + 4*len(bb.restarts)
}

func (bb *blockBuilder) finish() []byte {
	// Write the restart offsets too as a footer. All integers
	// are written as uvarints except for the final offset, which
	// is written as a uint32.
	// [r_0 | uvarint] [r_1 | uvarint] ... [r_n | uvarint] [offset of r_0 | uint32]
	restartsOff := bb.buf.Len()
	for _, off := range bb.restarts {
		bb.putUvarint(off)
	}
	var tmp [4]byte
	byteOrder.PutUint32(tmp[:], uint32(restartsOff))
	bb.buf.Write(tmp[:])
	return bb.buf.Bytes()
}

func (bb *blockBuilder) add(k, v []byte) {
	shared := 0
	if bb.counter >= bb.restartInterval {
		bb.restarts = append(bb.restarts, bb.buf.Len())
		bb.counter = 0
	} else {
		// Count how many characters are shared between k and lastKey
		minLen := len(bb.lastKey)
		if len(k) < minLen {
			minLen = len(k)
		}
		for shared < minLen && bb.lastKey[shared] == k[shared] {
			shared++
		}
	}
	nonshared := len(k) - shared

	// Write the lengths: shared key, unshared key, value
	bb.putUvarint(shared)
	bb.putUvarint(nonshared)
	bb.putUvarint(len(v))

	// Write the unshared key bytes and the value
	bb.buf.Write(k[shared:])
	bb.buf.Write(v)

	bb.lastKey = k
	bb.counter++
}

func (bb *blockBuilder) putUvarint(v int) {
	b := binary.PutUvarint(bb.tmp[:], uint64(v))
	bb.buf.Write(bb.tmp[:b])
}
*/
