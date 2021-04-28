package fileIndex

var (
	minLoadFactor    = 0.25
	maxLoadFactor    = 0.75
	defaultTableSize = 128 // 32768
)

/*type Record struct {
	Key   int
	Value int		// startBlock
	Value2 int		// endBlock
	Next  *Record
}

type Hash struct {
	Records []*Record
}

type HashTable struct {
	Table    *Hash
	NRecords *int
}*/

//// createHashTable: Called by checkLoadFactorAndUpdate when creating a new hash, for internal use only.
//func createHashTable(tableSize uint32) HashTable {
//	num := uint32(0)
//
//	hash := Hash{Records: make([]*Record, tableSize)}
//	return HashTable{Table: &hash, NRecords: &num}
//}
//
//// CreateHashTable: Called by the user to create a fileIndex.
//func CreateHashTable(initNum uint32) HashTable {
//	hash := Hash{Records: make([]*Record, defaultTableSize)}
//	return HashTable{Table: &hash, NRecords: &initNum}
//}
//
//// hashFunction: Used to calculate the index of record within the slice
//func hashFunction(key *uint32, size uint32) uint32 {
//	return *key % size
//}
//
//// put: inserts a Key into the hash Table, for internal use only
//func (h *HashTable) put(key *uint32, value *uint32, value2 *uint32, blocks []byte) bool {
//	index := hashFunction(key, uint32(len(h.Table.Records)))
//	iterator := h.Table.Records[index]
//	node := Record{Key: key, Value: value, Value2: value2, Blocks: blocks, Next: nil}
//	if iterator == nil || iterator.Key == nil {
//		h.Table.Records[index] = &node
//	} else {
//		tKey := uint32(0)
//		tVal := uint32(0)
//		tVal2 := uint32(0)
//		prev := &Record{Key: &tKey, Value: &tVal, Value2: &tVal2, Blocks: blocks, Next: nil}
//		for iterator != nil && iterator.Key != nil {
//			if *iterator.Key == *key { // Key already exists
//				iterator.Value = value
//				iterator.Value2 = value2
//				iterator.Blocks = blocks
//				return false
//			}
//			prev = iterator
//			iterator = iterator.Next
//		}
//		prev.Next = &node
//	}
//	*h.NRecords += 1
//	return true
//}
//
//// Put: inserts a Key into the hash Table (publicly callable)
//func (h *HashTable) Put(key uint32, value uint32, value2 uint32, blocks []byte) {
//	sizeChanged := h.put(&key, &value, &value2, blocks)
//	if sizeChanged {
//		h.checkLoadFactorAndUpdate()
//	}
//}
//
//// Get: Retrieve a Value for a Key from the hash Table (publicly callable)
//func (h *HashTable) Get(key uint32) (bool, *FileMetadata) {
//	index := hashFunction(&key, uint32(len(h.Table.Records)))
//	iterator := h.Table.Records[index]
//	for iterator != nil {
//		if iterator.Key != nil && *iterator.Key == key { // Key already exists
//			blocks := roaring.New()
//			_, _ = blocks.ReadFrom(bytes.NewReader(iterator.Blocks))
//			return true, &FileMetadata{FirstBlock: *iterator.Value, LastBlock: *iterator.Value2, Blocks: blocks}
//		}
//		iterator = iterator.Next
//	}
//	return false, nil
//}
//
//// del: remove a Key-Value record from the hash Table, for internal use only
//func (h *HashTable) del(key uint32) bool {
//	index := hashFunction(&key, uint32(len(h.Table.Records)))
//	iterator := h.Table.Records[index]
//	if iterator == nil {
//		return false
//	}
//	if iterator.Key == nil {
//		return true
//	}
//	if *iterator.Key == key {
//		h.Table.Records[index] = iterator.Next
//		*h.NRecords--
//		return true
//	} else {
//		prev := iterator
//		iterator = iterator.Next
//		for iterator != nil {
//			if *iterator.Key == key {
//				prev.Next = iterator.Next
//				*h.NRecords--
//				return true
//			}
//			prev = iterator
//			iterator = iterator.Next
//		}
//		return false
//	}
//}
//
//// Del: remove a Key-Value record from the hash Table (publicly available)
//func (h *HashTable) Del(key uint32) bool {
//	sizeChanged := h.del(key)
//	if sizeChanged {
//		h.checkLoadFactorAndUpdate()
//	}
//	return sizeChanged
//}
//
//// getLoadFactor: calculate the loadfactor for the fileIndex
//// Calculated as: number of Records stored / length of underlying slice used
//func (h *HashTable) getLoadFactor() float64 {
//	return float64(*h.NRecords) / float64(len(h.Table.Records))
//}
//
//// checkLoadFactorAndUpdate: if 0.25 > loadfactor or 0.75 < loadfactor,
//// update the underlying slice to have have loadfactor close to 0.5
//func (h *HashTable) checkLoadFactorAndUpdate() {
//	if *h.NRecords == 0 {
//		return
//	} else {
//		loadFactor := h.getLoadFactor()
//		if loadFactor < minLoadFactor {
//			fmt.Println("** Loadfactor below limit, reducing fileIndex size **")
//			hash := createHashTable(uint32(len(h.Table.Records) / 2))
//			for _, record := range h.Table.Records {
//				for record != nil {
//					if record.Key != nil {
//						hash.put(record.Key, record.Value, record.Value2, record.Blocks)
//					}
//					record = record.Next
//				}
//			}
//			h.Table = hash.Table
//		} else if loadFactor > maxLoadFactor {
//			fmt.Println("** Loadfactor above limit, increasing fileIndex size **")
//			hash := createHashTable(*h.NRecords * 2)
//			for _, record := range h.Table.Records {
//				for record != nil {
//					if record.Key == nil {
//						record = record.Next
//						continue
//					}
//					hash.put(record.Key, record.Value, record.Value2, record.Blocks)
//					record = record.Next
//				}
//			}
//			h.Table = hash.Table
//		}
//	}
//}
