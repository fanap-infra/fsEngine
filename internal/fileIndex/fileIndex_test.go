package fileIndex

import (
	"fmt"
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/go-playground/assert/v2"
)

func TestFileIndex_GetFileInfo(t *testing.T) {
	fi := NewFileIndex()
	fi.AddFile(1)
	fi.AddFile(30)
	bin, err := fi.GenerateBinary()
	assert.Equal(t, nil, err)
	fi2 := NewFileIndex()
	err = fi2.InitFromBinary(bin)
	assert.Equal(t, nil, err)
	fi2.AddFile(2)
	fi2.AddFile(34)
	fi2.EditFileMeta(2, FileMetadata{FirstBlock: 2, LastBlock: 200, Blocks: roaring.New()})
	fi2.EditFileMeta(34, FileMetadata{FirstBlock: 34, LastBlock: 2000, Blocks: roaring.New()})
	thirty, _ := fi2.GetFileInfo(2)
	thirtyFour, _ := fi2.GetFileInfo(34)
	assert.NotEqual(t, nil, thirty)
	assert.NotEqual(t, nil, thirtyFour)

	fmt.Println(fi2.GetFileInfo(30))
}
