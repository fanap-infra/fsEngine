package Header_

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

func TestFileIndex(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, err, nil)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eHandler := &EventsHandlerTest{}
	fs, err := CreateHeaderFS(fsID, homePath, fileSizeTest, blockSizeTest, log.GetScope("test"), eHandler, nil)
	assert.Equal(t, err, nil)

	err = fs.Close()
	assert.Equal(t, err, nil)

	fs2, err := ParseHeaderFS(fsID, homePath, log.GetScope("test2"), eHandler, nil)
	if !assert.Equal(t, err, nil) {
		return
	}

	assert.Equal(t, fs2.blmSize, fs.blmSize)

	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}
