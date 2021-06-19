package fileIndex

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestSize     = 5
	MaxSize      = 100
	TestByteSize = 100
)

func TestFileIndex_Marshal_UnMarshal(t *testing.T) {
	fi := NewFileIndex()

	testFiles := make([]File, TestSize)
	for i := 0; i < TestSize; i++ {
		testBytes := make([]byte, TestByteSize)
		n, err := rand.Read(testBytes)
		assert.Equal(t, TestByteSize, n)
		assert.Equal(t, nil, err)

		file := File{
			Name:      "test" + strconv.Itoa(i),
			LastBlock: uint32(rand.Intn(MaxSize)), FirstBlock: uint32(rand.Intn(MaxSize)), Id: uint32(rand.Intn(MaxSize)),
			RMapBlocks: testBytes,
		}

		err = fi.AddFile(file.Id, file.Name)
		if err != nil {
			if strings.Contains(err.Error(), "has been added before") {
				i = i - 1
				continue
			}
		}
		assert.Equal(t, nil, err)
		testFiles[i] = file
	}
	for i := 0; i < TestSize; i++ {
		err := fi.UpdateFile(testFiles[i].Id, testFiles[i].FirstBlock, testFiles[i].LastBlock, testFiles[i].Name, testFiles[i].RMapBlocks)
		assert.Equal(t, nil, err)
	}

	buf, err := fi.GenerateBinary()
	assert.Equal(t, nil, err)

	fiParsed := NewFileIndex()
	err = fiParsed.InitFromBinary(buf)
	assert.Equal(t, nil, err)

	for i := 0; i < TestSize; i++ {
		file, ok := fiParsed.table.Files[testFiles[i].Id]
		assert.Equal(t, true, ok)
		assert.Equal(t, testFiles[i].Id, file.GetId())
		assert.Equal(t, testFiles[i].Name, file.GetName())
		assert.Equal(t, testFiles[i].LastBlock, file.GetLastBlock())
		assert.Equal(t, testFiles[i].FirstBlock, file.GetFirstBlock())
		assert.Equal(t, testFiles[i].RMapBlocks, file.GetRMapBlocks())
	}
}
