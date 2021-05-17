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

func assertFetal(check bool)  {
	panic("Assert failed")
}

func TestVirtualFile_CRUD(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))
	var testIDs []uint32
	var testNames []string

	TestSize := 5
	MaxID := 1000
	var vfs []*virtualFile.VirtualFile
	for i := 0; i < TestSize; i++ {
		tmp := uint32( rand.Intn(MaxID))
		if utils.ItemExists(testIDs, tmp) {
			i = i - 1
			continue
		}
		testIDs = append(testIDs, tmp)
		testNames = append(testNames, "test"+strconv.Itoa(i))
		vf, err := fse.NewVirtualFile(testIDs[i], testNames[i])
		assert.Equal(t, nil, err)
		vfs = append(vfs, vf)
	}


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


