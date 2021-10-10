package fsEngine

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/fanap-infra/fsEngine/internal/constants"

	"github.com/fanap-infra/fsEngine/pkg/virtualFile"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

type EventsListener struct {
	t      *testing.T
	fileID uint32
}

func (el *EventsListener) VirtualFileDeleted(fileID uint32, message string) {
	// log.Infov("VirtualFileDeleted event", "fileID", fileID, "message", message,
	//	"first fileID", el.fileID)
	assert.Equal(el.t, el.fileID, fileID)
}

func TestPrepareAndParseBlock(t *testing.T) {
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
		buf2, err := fse.parseBlock(buf, uint32(blockID), uint32(fileIdTest))
		assert.Equal(t, nil, err)
		assert.Equal(t, buf2, token)

		blockID++
	}
}

func TestPrepareAndParseBlockManually(t *testing.T) {
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

	blockID := 8201
	fileIdTest := 2100
	previousBlock := 0
	// for i := 0; i < numberOfTests; i++ {
	token := make([]byte, uint32(rand.Intn(blockSizeTest)))
	_, err = rand.Read(token)
	assert.Equal(t, nil, err)
	buf, err := fse.prepareBlock(token, uint32(fileIdTest), uint32(previousBlock), uint32(blockID))
	assert.Equal(t, nil, err)
	buf2, err := fse.parseBlock(buf, uint32(blockID), uint32(fileIdTest))
	assert.Equal(t, nil, err)
	assert.Equal(t, buf2, token)

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}

func TestFSEngine_NoSpace(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(fsID, homePath, fileSizeTest*2, blockSizeTest, &eventListener, log.GetScope("test"), nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))

	MaxID := 1000
	numberOfVFs := 5
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := (fileSizeTest / numberOfVFs) + blockSizeTest

	virtualFiles := make([]*virtualFile.VirtualFile, 0)

	bytes := make([][][]byte, numberOfVFs)
	vfIDs := make([]uint32, 0)
	for i := 0; i < numberOfVFs; i++ {
		vfID := uint32(rand.Intn(MaxID))
		if utils.ItemExists(vfIDs, vfID) {
			i = i - 1
			continue
		}
		vfIDs = append(vfIDs, vfID)
		vf, err := fse.NewVirtualFile(vfID, "test"+strconv.Itoa(i))
		if assert.Equal(t, nil, err) {
			virtualFiles = append(virtualFiles, vf)
		}
	}

	if len(virtualFiles) != numberOfVFs {
		return
	}

	eventListener.fileID = vfIDs[0]

	for j, vf := range virtualFiles {
		size := 0
		for {
			token := make([]byte, uint32(1+rand.Intn(MaxByteArraySize)))
			m, err := rand.Read(token)
			assert.Equal(t, nil, err)
			bytes[j] = append(bytes[j], token)
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
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}
