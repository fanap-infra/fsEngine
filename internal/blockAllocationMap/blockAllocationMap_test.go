package blockAllocationMap

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const (
	MaxSize = 10
	Path    = "/test.beh"
)

type EventTest struct {
	count uint32
}

func (eT *EventTest) NoSpace() uint32 {
	eT.count = eT.count + 1
	return eT.count - 1
}

func TestAddingAndReadAll(t *testing.T) {
	eventTest := &EventTest{count: 0}
	bAllocationMap := New(log.GetScope("test"), eventTest, MaxSize)
	for i := 0; i < MaxSize; i++ {
		err := bAllocationMap.SetBlockAsAllocated(uint32(i))
		assert.Equal(t, nil, err)
	}

	err := bAllocationMap.SetBlockAsAllocated(uint32(rand.Intn(MaxSize)))
	assert.NotEqual(t, nil, err)
	for i := 0; i < MaxSize; i++ {
		assert.Equal(t, bAllocationMap.IsBlockAllocated(uint32(i)), true)
	}
	for i := 0; i < MaxSize; i++ {
		f := bAllocationMap.FindNextFreeBlockAndAllocate()
		assert.Equal(t, f, uint32(0))
	}
}

func TestMarshalAndUnmarshal(t *testing.T) {
	eventTest := &EventTest{count: 0}
	bAllocationMap := New(log.GetScope("test"), eventTest, MaxSize)
	TestSize := 5
	var testBlocks []uint32
	for i := 0; i < TestSize; i++ {
		tmp := uint32(rand.Intn(MaxSize))
		if utils.ItemExists(testBlocks, tmp) {
			i = i - 1
			continue
		}
		testBlocks = append(testBlocks, tmp)
		err := bAllocationMap.SetBlockAsAllocated(tmp)
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, TestSize, len(testBlocks))
	b, err := Marshal(bAllocationMap)
	assert.Equal(t, err, nil)
	bAllocationMap2, err := Open(log.GetScope("test"), eventTest, MaxSize, bAllocationMap.LastWrittenBlock, b)
	assert.Equal(t, err, nil)
	for i := 0; i < len(testBlocks); i++ {
		assert.Equal(t, true, bAllocationMap2.IsBlockAllocated(testBlocks[i]))
	}
	// log.Infov("data test", "testBlocks", testBlocks)
	for i := 0; i < MaxSize; i++ {
		if utils.ItemExists(testBlocks, uint32(i)) {
			// log.Infov("Skip this item", "i", i)
			continue
		}
		if !assert.Equal(t, false, bAllocationMap2.IsBlockAllocated(uint32(i))) {
			log.Infov("Block allocated wrong", "i", i)
		}
	}
}

func TestBlockAllocationMap_ToArray(t *testing.T) {
	eventTest := &EventTest{count: 0}
	bAllocationMap := New(log.GetScope("test"), eventTest, MaxSize)
	TestSize := 5
	var testBlocks []uint32
	for i := 0; i < TestSize; i++ {
		tmp := uint32(rand.Intn(MaxSize))
		if utils.ItemExists(testBlocks, tmp) {
			i = i - 1
			continue
		}
		testBlocks = append(testBlocks, tmp)
		err := bAllocationMap.SetBlockAsAllocated(tmp)
		assert.Equal(t, nil, err)
	}
	assert.Equal(t, TestSize, len(testBlocks))
	blocksIndex := bAllocationMap.ToArray()
	assert.Equal(t, TestSize, len(blocksIndex))
	var testBlocksInt []int
	var blocksInt []int
	for i := 0; i < TestSize; i++ {
		testBlocksInt = append(testBlocksInt, int(testBlocks[i]))
		blocksInt = append(blocksInt, int(blocksIndex[i]))
	}
	sort.Ints(testBlocksInt)
	sort.Ints(blocksInt)
	for i := 0; i < TestSize; i++ {
		assert.Equal(t, testBlocksInt[i], blocksInt[i])
	}
}

func TestFullState(t *testing.T) {
	eventTest := &EventTest{count: 0}
	bAllocationMap := New(log.GetScope("test"), eventTest, MaxSize)
	for i := 0; i < MaxSize; i++ {
		blockIndex := bAllocationMap.FindNextFreeBlockAndAllocate()
		err := bAllocationMap.SetBlockAsAllocated(blockIndex)
		if i == MaxSize-1 {
			// last written block is zero
			assert.Equal(t, uint32(0), blockIndex)
			assert.Equal(t, uint32(0), bAllocationMap.LastWrittenBlock)
		} else {
			assert.Equal(t, uint32(i+1), blockIndex)
			assert.Equal(t, uint32(i+1), bAllocationMap.LastWrittenBlock)
		}
		assert.Equal(t, nil, err)
	}
	eventTest.count = 0
	for i := 0; i < MaxSize; i++ {
		blockIndex := bAllocationMap.FindNextFreeBlockAndAllocate()
		err := bAllocationMap.SetBlockAsAllocated(blockIndex)
		assert.Equal(t, uint32(i), blockIndex)
		assert.Equal(t, uint32(i), bAllocationMap.LastWrittenBlock)
		assert.Equal(t, nil, err)
	}
}
