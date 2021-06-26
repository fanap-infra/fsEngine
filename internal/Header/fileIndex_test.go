package Header_

import (
	"os"
	"testing"

	"github.com/fanap-infra/FSEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestFileIndex(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, err, nil)
	_ = utils.DeleteFile(homePath + path)
	eHandler := &EventsHandlerTest{}
	fs, err := CreateHeaderFS(homePath+path, BLOCKSIZE*1000, BLOCKSIZE, log.GetScope("test"), eHandler)
	assert.Equal(t, err, nil)

	err = fs.Close()
	assert.Equal(t, err, nil)

	fs2, err := ParseHeaderFS(homePath+path, log.GetScope("test2"), eHandler)
	if !assert.Equal(t, err, nil) {
		return
	}

	assert.Equal(t, fs2.blmSize, fs.blmSize)

	_ = utils.DeleteFile(homePath + path)
}
