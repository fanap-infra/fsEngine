package fsEngine

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/log"

	"github.com/stretchr/testify/assert"
)

const (
	blockSizeTest = 5120
	fileSizeTest  = blockSizeTest * 128
)

func TestCreateFS(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	_, err = CreateFileSystem(homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+fsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+headerPath))
	size, err := utils.FileSize(homePath + "/" + fsPath)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(fileSizeTest), size)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}

func TestParseFS(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	_, err = CreateFileSystem(homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	fs, err := ParseFileSystem(homePath+"/"+fsPath, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, fs.blockSize, uint32(blockSizeTest))
	assert.Equal(t, fs.maxNumberOfBlocks, uint32(fileSizeTest/blockSizeTest))
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}
