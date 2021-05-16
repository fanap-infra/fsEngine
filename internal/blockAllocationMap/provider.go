package blockAllocationMap

import (
	"bytes"

	"github.com/RoaringBitmap/roaring"
	"github.com/fanap-infra/log"
)

func New(log *log.Logger, listener Events, maxSize uint32) *BlockAllocationMap {
	return &BlockAllocationMap{
		maxSize:          maxSize,
		log:              log,
		trigger:          listener,
		LastWrittenBlock: 0,
		rMap:             roaring.NewBitmap(),
	}
}

// ToDo: add block allocation map parser

func Open(log *log.Logger, listener Events, maxSize uint32, lastWrittenBlock uint32, rMapByte []byte) (*BlockAllocationMap, error) {
	rMap := roaring.NewBitmap()
	b := bytes.NewReader(rMapByte)
	_, err := rMap.ReadFrom(b)
	if err != nil {
		return nil, err
	}
	return &BlockAllocationMap{
		maxSize:          maxSize,
		log:              log,
		trigger:          listener,
		LastWrittenBlock: lastWrittenBlock,
		rMap:             rMap,
	}, nil
}

func Marshal(bAllocationMap *BlockAllocationMap) ([]byte, error) {
	bAllocationMap.rMap.RunOptimize()
	byteArray := bytes.Buffer{}
	_, err := bAllocationMap.rMap.WriteTo(&byteArray)
	if err != nil {
		return nil, err
	}
	return byteArray.Bytes(), nil
}
