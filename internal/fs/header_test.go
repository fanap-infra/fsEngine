package fs

import (
	"behnama/stream/pkg/fsEngine/pkg/utils"
	"os"
	"testing"

	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const path = "/test.beh"

func TestStoreHeader(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, err, nil)
	_ = utils.DeleteFile(homePath + path)
	fs, err := CreateFileSystem(homePath+path, BLOCKSIZE*1000, BLOCKSIZE, log.GetScope("test"))
	assert.Equal(t, err, nil)
	size := fs.size
	version := fs.version
	blocks := fs.blocks
	blockSize := fs.blockSize
	lastWrittenBlock := fs.lastWrittenBlock
	err = fs.Close()
	assert.Equal(t, err, nil)
	fs2, err := ParseFileSystem(homePath+path, log.GetScope("test2"))
	assert.Equal(t, err, nil)
	assert.Equal(t, size, fs2.size)
	assert.Equal(t, version, fs2.version)
	assert.Equal(t, blocks, fs2.blocks)
	assert.Equal(t, blockSize, fs2.blockSize)
	assert.Equal(t, lastWrittenBlock, fs2.lastWrittenBlock)
}

//func TestHeaderParsing(t *testing.T) {
// homePath, err := os.UserHomeDir()
// assert.Equal(t, err, nil)
// fs, err := ParseFileSystem(homePath+path, log.GetScope("test2"))
// assert.Equal(t, err, nil)
// err = fs.Close()
//}
