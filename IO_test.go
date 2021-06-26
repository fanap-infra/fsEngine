package fsEngine

import (
	"math/rand"
	"os"
	"testing"

	"github.com/fanap-infra/FSEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestIO_OneVirtualFile(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))
	var bytes [][]byte

	MaxID := 1000
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := int(1.5 * blockSizeTest)
	vfID := uint32(rand.Intn(MaxID))
	vf, err := fse.NewVirtualFile(vfID, "test")
	assert.Equal(t, nil, err)
	size := 0

	for {
		token := make([]byte, uint32(rand.Intn(MaxByteArraySize)))
		m, err := rand.Read(token)
		assert.Equal(t, nil, err)
		bytes = append(bytes, token)
		size = size + m
		n, err := vf.Write(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, m, n)

		if size > VFSize {
			break
		}
	}

	err = vf.Close()
	assert.Equal(t, nil, err)

	vf, err = fse.OpenVirtualFile(vfID)
	assert.Equal(t, nil, err)

	for _, v := range bytes {
		buf := make([]byte, len(v))
		_, err := vf.Read(buf)

		assert.Equal(t, nil, err)
		if err != nil {
			break
		}

		assert.Equal(t, v, buf)
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
}
