package fsEngine

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/mocks"

	"github.com/fanap-infra/fsEngine/internal/constants"
	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestCreateFS_Redis(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	redisMock := mocks.NewRedisMock()
	_, err = CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"), &redisMock)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))
	size, err := utils.FileSize(homePath + "/" + constants.FsPath)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(fileSizeTest), size)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}

func TestParseFS_Redis(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	redisMock := mocks.NewRedisMock()
	_, err = CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"), &redisMock)
	assert.Equal(t, nil, err)
	fs, err := ParseFileSystem(fsID, homePath, &eventListener, log.GetScope("test"), &redisMock)
	assert.Equal(t, nil, err)
	assert.Equal(t, fs.blockSize, uint32(blockSizeTest))
	assert.Equal(t, fs.maxNumberOfBlocks, uint32(fileSizeTest/blockSizeTest))
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}
