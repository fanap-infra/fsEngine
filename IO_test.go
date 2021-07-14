package fsEngine

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/fanap-infra/fsEngine/pkg/virtualFile"

	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestIO_OneVirtualFile(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))
	var bytes [][]byte

	MaxID := 1000
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := int(3.5 * blockSizeTest)
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

	vf2, err := fse.OpenVirtualFile(vfID)
	assert.Equal(t, nil, err)

	for _, v := range bytes {
		buf := make([]byte, len(v))
		m, err := vf2.Read(buf)

		assert.Equal(t, nil, err)
		if err != nil {
			break
		}

		assert.Equal(t, len(v), m)
		assert.Equal(t, v, buf)
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
}

func TestIO_MultipleVirtualFileConsecutively(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))

	MaxID := 1000
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := int(3.5 * blockSizeTest)

	virtualFiles := make([]*virtualFile.VirtualFile, 0)
	numberOfVFs := 5
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

	for j, vf := range virtualFiles {
		size := 0
		for {
			token := make([]byte, uint32(rand.Intn(MaxByteArraySize)))
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

	for i, vBytes := range bytes {
		vf2, err := fse.OpenVirtualFile(vfIDs[i])
		assert.Equal(t, nil, err)

		for _, v := range vBytes {
			buf := make([]byte, len(v))
			m, err := vf2.Read(buf)

			assert.Equal(t, nil, err)
			if err != nil {
				break
			}

			assert.Equal(t, len(v), m)
			assert.Equal(t, v, buf)
		}
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
}

// ToDo: make it to write virtual file without any sort
func TestIO_MultipleVirtualFile(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath+fsPathTest, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+fsPathTest))
	assert.Equal(t, true, utils.FileExists(homePath+headerPathTest))

	MaxID := 1000
	MaxByteArraySize := int(blockSizeTest * 0.5)
	VFSize := int(3.5 * blockSizeTest)

	virtualFiles := make([]*virtualFile.VirtualFile, 0)
	numberOfVFs := 5
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

	for j, vf := range virtualFiles {
		size := 0
		for {
			token := make([]byte, uint32(rand.Intn(MaxByteArraySize)))
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

	for i, vBytes := range bytes {
		vf2, err := fse.OpenVirtualFile(vfIDs[i])
		assert.Equal(t, nil, err)

		for _, v := range vBytes {
			buf := make([]byte, len(v))
			m, err := vf2.Read(buf)

			assert.Equal(t, nil, err)
			if err != nil {
				break
			}

			assert.Equal(t, len(v), m)
			assert.Equal(t, v, buf)
		}
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + fsPathTest)
	_ = utils.DeleteFile(homePath + headerPathTest)
}
