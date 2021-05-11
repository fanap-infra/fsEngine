package blockAllocationMap

import (
	"behnama/stream/pkg/fsEngine/pkg/utils"
	"math/rand"
	"testing"

	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const (
	MaxSize = 100
	Path    = "/test.beh"
)

type EventTest struct {
	count uint32
}

func (eT EventTest) NoSpace() uint32 {
	eT.count = eT.count + 1
	return eT.count - 1
}

func TestAddingAndReadAll(t *testing.T) {
	evetTest := &EventTest{count: 0}
	bAllocationMap := New(log.GetScope("test"), evetTest, MaxSize)
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
	log.Infov("data test", "testBlocks", testBlocks)
	for i := 0; i < MaxSize; i++ {
		if utils.ItemExists(testBlocks, i) {
			log.Infov("Skip this item", "i", i)
			continue
		}
		if !assert.Equal(t, false, bAllocationMap2.IsBlockAllocated(uint32(i))) {
			log.Infov("Block allocated wrong", "i", i)
		}
	}
}