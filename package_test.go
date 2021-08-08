package fsEngine

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/fanap-infra/fsEngine/internal/constants"
	"github.com/fanap-infra/fsEngine/pkg/utils"
	"github.com/fanap-infra/fsEngine/pkg/virtualFile"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

func TestIO_ConcurrentMultipleVirtualFile(t *testing.T) {
	homePath, err := os.UserHomeDir()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
	eventListener := EventsListener{t: t}
	blockSize := uint32(1 << 19)
	fileSize := int64(10000 * blockSize)
	fse, err := CreateFileSystem(homePath, fileSize, blockSize, &eventListener, log.GetScope("test"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.FsPath))
	assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderPath))
	// assert.Equal(t, true, utils.FileExists(homePath+"/"+constants.HeaderBackUpPath))

	MaxID := 1000
	MaxByteArraySize := int(float32(blockSize) * 1.7)
	VFSize := int(11.3 * float32(blockSize))

	virtualFiles := make([]*virtualFile.VirtualFile, 0)
	numberOfVFs := 25
	bytes := make([][][]byte, numberOfVFs)
	vfIDs := make([]uint32, 0)
	var wg sync.WaitGroup
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
	var mu sync.Mutex
	for j, vf := range virtualFiles {
		wg.Add(1)
		go func(j int, vf *virtualFile.VirtualFile) {
			defer wg.Done()
			size := 0
			for {
				token := make([]byte, uint32(rand.Intn(MaxByteArraySize)))
				m, err := rand.Read(token)
				assert.Equal(t, nil, err)
				mu.Lock()
				bytes[j] = append(bytes[j], token)
				mu.Unlock()
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
		}(j, vf)

	}
	wg.Wait()

	for i, vBytes := range bytes {
		vf2, err := fse.OpenVirtualFile(vfIDs[i])
		assert.Equal(t, nil, err)
		wg.Add(1)
		go func(vBytes [][]byte, vf *virtualFile.VirtualFile) {
			defer wg.Done()
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
		}(vBytes, vf2)
	}
	wg.Wait()

	err = fse.Close()
	assert.Equal(t, nil, err)
	_ = utils.DeleteFile(homePath + "/" + constants.FsPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderPath)
	_ = utils.DeleteFile(homePath + "/" + constants.HeaderBackUpPath)
}
