package Header_

import (
	"os"
	"testing"

	"github.com/fanap-infra/FSEngine/pkg/utils"

	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const path = "/test.beh"

type EventsHandlerTest struct {
	count uint32
}

func (eT *EventsHandlerTest) NoSpace() uint32 {
	eT.count = eT.count + 1
	return eT.count - 1
}

func TestStoreHeader(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, err, nil)
	_ = utils.DeleteFile(homePath + path)
	eHandler := &EventsHandlerTest{}
	fs, err := CreateHeaderFS(homePath+path, BLOCKSIZE*1000, BLOCKSIZE, log.GetScope("test"), eHandler)
	assert.Equal(t, err, nil)
	size := fs.size
	version := fs.version
	maxNumberOfBlocks := fs.maxNumberOfBlocks
	blockSize := fs.blockSize
	lastWrittenBlock := fs.lastWrittenBlock

	buf := make([]byte, fs.blmSize)

	m, err := fs.file.ReadAt(buf, BlockAllocationMapByteIndex)
	assert.Equal(t, err, nil)
	assert.Equal(t, fs.blmSize, uint32(m))

	err = fs.Close()
	assert.Equal(t, err, nil)

	fs2, err := ParseHeaderFS(homePath+path, log.GetScope("test2"), eHandler)
	if !assert.Equal(t, err, nil) {
		return
	}
	assert.Equal(t, size, fs2.size)
	assert.Equal(t, version, fs2.version)
	assert.Equal(t, maxNumberOfBlocks, fs2.maxNumberOfBlocks)
	assert.Equal(t, blockSize, fs2.blockSize)
	assert.Equal(t, lastWrittenBlock, fs2.lastWrittenBlock)

	_ = utils.DeleteFile(homePath + path)
}

//func TestHeaderParsing(t *testing.T) {
// homePath, err := os.UserHomeDir()
// assert.Equal(t, err, nil)
// fs, err := ParseFileSystem(homePath+path, log.GetScope("test2"))
// assert.Equal(t, err, nil)
// err = fs.Close()
//}
