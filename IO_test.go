package fsEngine

import (
	"github.com/fanap-infra/FSEngine/internal/virtualFile"
	"github.com/fanap-infra/FSEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
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

	TestSize := 5
	MaxID := 1000
	MaxByteArraySize := int(blockSizeTest*0.5)
	VFSize := 3*blockSizeTest
	
	vf, err := fse.NewVirtualFile(uint32( rand.Intn(MaxID)), "test")
	assert.Equal(t, nil, err)
	size := 0
	for  {
		token := make([]byte, uint32( rand.Intn(MaxByteArraySize)))
		m, err := rand.Read(token)
		assert.Equal(t, nil, err)
		bytes  = append(bytes, token)
		size = size + m
		n,err := vf.Write(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, m, n)

		if size > VFSize {
			break
		}
	}

	vf.Read()


	for i := 0; i < TestSize; i++ {
		_, err := fse.NewVirtualFile(testIDs[i], testNames[i])
		assert.NotEqual(t, nil, err)
	}

	for i := 0; i < TestSize; i++ {
		err := fse.RemoveVirtualFile(testIDs[i])
		assert.NotEqual(t, nil, err)
	}

	for i := 0; i < TestSize/2; i++ {
		err := vfs[i].Close()
		assert.Equal(t, nil, err)
	}

	for i := 0; i < TestSize; i++ {
		if i< TestSize/2 {
			err := fse.RemoveVirtualFile(testIDs[i])
			assert.Equal(t, nil, err)
		} else {
			err := fse.RemoveVirtualFile(testIDs[i])
			assert.NotEqual(t, nil, err)
			err = vfs[i].Close()
			assert.Equal(t, nil, err)
			err = fse.RemoveVirtualFile(testIDs[i])
			assert.Equal(t, nil, err)
		}
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
}


