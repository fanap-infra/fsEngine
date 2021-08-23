package Header_

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const (
	fsPath     = "fs.beh"
	headerPath = "Header.Beh"
	fsID       = 11
)

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
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eHandler := &EventsHandlerTest{}
	fs, err := CreateHeaderFS(fsID, homePath, fileSizeTest, blockSizeTest, log.GetScope("test"), eHandler, nil)
	assert.Equal(t, err, nil)
	size := fs.size
	version := fs.version
	maxNumberOfBlocks := fs.maxNumberOfBlocks
	blockSize := fs.blockSize
	lastWrittenBlock := fs.lastWrittenBlock
	assert.Equal(t, blockSize, uint32(blockSizeTest))
	assert.Equal(t, size, int64(fileSizeTest))

	buf := make([]byte, fs.blmSize)

	m, err := fs.readAt(buf, BlockAllocationMapByteIndex)
	assert.Equal(t, err, nil)
	assert.Equal(t, fs.blmSize, uint32(m))

	err = fs.Close()
	assert.Equal(t, err, nil)

	fs2, err := ParseHeaderFS(fsID, homePath, log.GetScope("test2"), eHandler, nil)
	if !assert.Equal(t, err, nil) {
		return
	}
	assert.Equal(t, size, fs2.size)
	assert.Equal(t, version, fs2.version)
	assert.Equal(t, maxNumberOfBlocks, fs2.maxNumberOfBlocks)
	assert.Equal(t, blockSize, fs2.blockSize)
	assert.Equal(t, lastWrittenBlock, fs2.lastWrittenBlock)

	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}

//func TestHeaderParsing(t *testing.T) {
// homePath, err := os.UserHomeDir()
// assert.Equal(t, err, nil)
// fs, err := ParseFileSystem(homePath+path, log.GetScope("test2"))
// assert.Equal(t, err, nil)
// err = fs.Close()
//}
