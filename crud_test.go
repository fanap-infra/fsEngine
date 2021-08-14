package fsEngine

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/fanap-infra/fsEngine/internal/constants"
	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/fsEngine/pkg/virtualFile"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

//func assertFetal(check bool) {
//	panic("Assert failed")
//}

func TestVirtualFile_Remove(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"), nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))
	var testIDs []uint32
	var testNames []string

	TestSize := 5
	MaxID := 1000
	var vfs []*virtualFile.VirtualFile
	for i := 0; i < TestSize; i++ {
		tmp := uint32(rand.Intn(MaxID))
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

	// cna not remove opened virtual files
	for i := 0; i < TestSize; i++ {
		err := fse.RemoveVirtualFile(testIDs[i])
		assert.NotEqual(t, nil, err)
	}

	for i := 0; i < TestSize/2; i++ {
		err := vfs[i].Close()
		assert.Equal(t, nil, err)
	}

	for i := 0; i < TestSize; i++ {
		if i < TestSize/2 {
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
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}

func TestVirtualFile_Open(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	var eventListener EventsListener
	fse, err := CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest, &eventListener,
		log.GetScope("test"), nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))
	var testIDs []uint32
	var testNames []string

	TestSize := 1
	MaxID := 1000
	var vfs []*virtualFile.VirtualFile
	for i := 0; i < TestSize; i++ {
		tmp := uint32(rand.Intn(MaxID))
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

	for i := 0; i < len(vfs); i++ {
		err := vfs[i].Close()
		assert.Equal(t, nil, err)
	}

	for i := 0; i < len(testIDs); i++ {
		_, err := fse.OpenVirtualFile(testIDs[i])
		assert.Equal(t, nil, err)
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}

func TestVirtualFile_RemoveUnsetBlocks(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(fsID, homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"), nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))
	var testIDs []uint32
	var testNames []string
	blocksIndexes := make([][]uint32, 0)
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := int(3.5 * blockSizeTest)
	TestSize := 5
	MaxID := 1000
	var vfs []*virtualFile.VirtualFile
	for i := 0; i < TestSize; i++ {
		tmp := uint32(rand.Intn(MaxID))
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

	for _, vf := range vfs {
		size := 0
		for {
			token := make([]byte, uint32(rand.Intn(MaxByteArraySize)))
			m, err := rand.Read(token)
			assert.Equal(t, nil, err)
			size = size + m
			n, err := vf.Write(token)
			assert.Equal(t, nil, err)
			assert.Equal(t, m, n)

			if size > VFSize {
				break
			}
		}
		blocksIndexes = append(blocksIndexes, fse.header.GetBLMArray())
		err = vf.Close()
		assert.Equal(t, nil, err)
	}

	for i := 0; i < len(testIDs); i++ {
		err := fse.RemoveVirtualFile(testIDs[i])
		assert.Equal(t, nil, err)
		for j := 0; j < len(blocksIndexes[i]); j++ {
			assert.Equal(t, false, fse.header.IsBlockAllocated(blocksIndexes[i][j]))
		}
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}
