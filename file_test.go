package fsEngine

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/internal/constants"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestFSEngine_GetFilePath(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	fs, err := CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"), nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, homePath+"/"+constants.FsPath, fs.GetFilePath())
	err = fs.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}
