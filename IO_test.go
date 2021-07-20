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
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+fsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+headerPath))
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
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}

func TestIO_MultipleVirtualFileConsecutively(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+fsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+headerPath))

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
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}

// ToDo: make it to write virtual file without any sort
func TestIO_MultipleVirtualFile(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+fsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+headerPath))

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
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}

func TestIO_ChangeSeekPointer(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
	eventListener := EventsListener{t: t}
	fse, err := CreateFileSystem(homePath, fileSizeTest, blockSizeTest, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+fsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+headerPath))
	var bytes []byte

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
		bytes = append(bytes, token...)
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
	segmentSize := 10
	readedSize := 0
	for {
		token := make([]byte, segmentSize)
		n, err := vf2.Read(token)

		assert.Equal(t, false, err != nil && n == 0)
		assert.Equal(t, token[:n], bytes[readedSize:readedSize+n])
		readedSize = readedSize + n
		if readedSize == size {
			break
		}
	}

	testCounter := 0

	for {
		token := make([]byte, segmentSize)
		seekTest := uint32(rand.Intn(size - segmentSize))
		err := vf2.ChangeSeekPointer(int64(seekTest))
		assert.Equal(t, nil, err)
		n, err := vf2.Read(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, n, segmentSize)
		assert.Equal(t, int(seekTest)+n, vf2.GetSeek())
		assert.Equal(t, token, bytes[seekTest:int(seekTest)+n])
		testCounter++
		if testCounter == 5 {
			break
		}
	}

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + fsPath)
	_ = utils.DeleteFile(homePath + "/" + headerPath)
}
