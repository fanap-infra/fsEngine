package fsEngine

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestFSEngine_GetFilePath(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	fs, err := CreateFileSystem(homePath, fileSizeTest, blockSizeTest,
		&eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, homePath+"/"+fsPath, fs.GetFilePath())
	err = fs.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}
