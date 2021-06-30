package fsEngine

import (
	"math/rand"
	"os"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestPrepareAndParseBlock(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))

	numberOfTests := 5
	blockID := 0
	fileIdTest := 5
	previousBlock := 0
	for i := 0; i < numberOfTests; i++ {
		token := make([]byte, uint32(rand.Intn(blockSizeTest)))
		_, err := rand.Read(token)
		assert.Equal(t, nil, err)
		buf, err := fse.prepareBlock(token, uint32(fileIdTest), uint32(previousBlock), uint32(blockID))
		assert.Equal(t, nil, err)
		buf2, err := fse.parseBlock(buf)
		assert.Equal(t, nil, err)
		assert.Equal(t, buf2, token)

		blockID++
	}
}
