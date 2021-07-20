package Header_

import (
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestFileIndex(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, err, nil)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eHandler := &EventsHandlerTest{}
	fs, err := CreateHeaderFS(homePath+"/"+headerPath, BLOCKSIZE*1000, BLOCKSIZE, log.GetScope("test"), eHandler)
	assert.Equal(t, err, nil)

	err = fs.Close()
	assert.Equal(t, err, nil)

	fs2, err := ParseHeaderFS(homePath+"/"+headerPath, log.GetScope("test2"), eHandler)
	if !assert.Equal(t, err, nil) {
		return
	}

	assert.Equal(t, fs2.blmSize, fs.blmSize)

	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}
