package fileIndex

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileIndex_CRUDFileIndex(t *testing.T) {
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

	for i := 0; i < TestSize/2; i++ {
		err := fi.RemoveFile(testFiles[i].Id)
		assert.Equal(t, nil, err)
	}

	for i := 0; i < TestSize; i++ {
		checked := fi.CheckFileExistWithLock(testFiles[i].Id)
		if i < TestSize/2 {
			assert.Equal(t, false, checked)
		} else {
			assert.Equal(t, true, checked)
		}
	}
}
