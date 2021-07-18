package virtualFile

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/log"
	"github.com/stretchr/testify/assert"
)

const (
	blockSizeTest     = 5120
	maxNumberOfBlocks = 3
	vfID              = 1
)

var byte2D [][]byte

type FSMock struct {
	vBuf        []byte
	seekPointer int
	vBufBlocks  [][]byte
	openFiles   map[uint32]*VirtualFile
	tst         *testing.T
}

func (fsMock *FSMock) Write(data []byte, fileID uint32) (int, error) {
	fsMock.vBuf = append(fsMock.vBuf, data...)
	counter := 0
	for {
		if (len(data) - counter) < blockSizeTest {
			// log.Infov("FSMock Write smaller than blockSizeTest",
			//	"blockID", len(fsMock.vBufBlocks), "data size", len(data), "counter", counter)
			err := fsMock.openFiles[uint32(fileID)].AddBlockID(uint32(len(fsMock.vBufBlocks)))
			if err != nil {
				return 0, err
			}
			fsMock.vBufBlocks = append(fsMock.vBufBlocks, data[counter:])
			counter = len(data)
		} else {
			// log.Infov("FSMock Write greater than blockSizeTest",
			//	"blockID", len(fsMock.vBufBlocks), "data size", len(data), "counter", counter)
			err := fsMock.openFiles[uint32(fileID)].AddBlockID(uint32(len(fsMock.vBufBlocks)))
			if err != nil {
				return 0, err
			}
			fsMock.vBufBlocks = append(fsMock.vBufBlocks, data[counter:blockSizeTest])
			counter = counter + blockSizeTest
		}
		if counter >= len(data) {
			if counter != len(data) {
				log.Warnv("counter greater than data", "counter", counter, "len(data)", len(data))
			}
			return len(data), nil
		}
	}
}

func (fsMock *FSMock) WriteAt(data []byte, off int64, fileID uint32) (int, error) {
	fsMock.vBuf = append(fsMock.vBuf, data...)
	return len(data), nil
}

func (fsMock *FSMock) Read(data []byte, fileID uint32) (int, error) {
	data = fsMock.vBuf[fsMock.seekPointer : fsMock.seekPointer+len(data)]
	fsMock.seekPointer = fsMock.seekPointer + len(data)
	return len(data), nil
}

func (fsMock *FSMock) ReadAt(data []byte, off int64, fileID uint32) (int, error) {
	return len(data), nil
}

func (fsMock *FSMock) ReadBlock(blockIndex uint32) ([]byte, error) {
	return fsMock.vBufBlocks[blockIndex], nil
}

func (fsMock *FSMock) Closed(fileID uint32) error {
	return nil
}

func (fsMock *FSMock) NoSpace() uint32 {
	return 0
}

func (fsMock *FSMock) BAMUpdated(fileID uint32, bam []byte) error {
	return nil
}

func (fsMock *FSMock) UpdateFileIndexes(fileID uint32, firstBlock uint32, lastBlock uint32, fileSize uint32) error {
	return nil
}

func (fsMock *FSMock) UpdateFileOptionalData(fileId uint32, info []byte) error {
	return nil
}

func NewVBufMock(t *testing.T) *FSMock {
	return &FSMock{
		seekPointer: 0,
		openFiles:   make(map[uint32]*VirtualFile),
		tst:         t,
	}
}

func TestIO_WR(t *testing.T) {
	fsMock := NewVBufMock(t)
	blm := blockAllocationMap.New(log.GetScope("test"), fsMock, maxNumberOfBlocks)
	vf := NewVirtualFile("test", vfID, blockSizeTest, fsMock, blm,
		int(blockSizeTest)*2, log.GetScope("test2"))
	fsMock.openFiles[vfID] = vf

	size := 0
	VFSize := int(1.5 * blockSizeTest)
	MaxByteArraySize := int(blockSizeTest * 0.5)
	for {
		token := make([]byte, uint32(rand.Intn(MaxByteArraySize))+1)
		m, err := rand.Read(token)
		assert.Equal(t, nil, err)
		byte2D = append(byte2D, token)
		assert.Equal(t, m, len(token))
		size = size + m
		n, err := vf.Write(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, m, n)

		if size > VFSize {
			break
		}
	}

	err := vf.Close()
	assert.Equal(t, nil, err)
	// log.Infov("writting finished ",
	//	"number of blocks", len(vf.blockAllocationMap.ToArray()), "size", size,
	//	"len(fsMock.vBuf)", len(fsMock.vBuf))
	assert.Equal(t, size, len(fsMock.vBuf))
	counter := 0
	var vBlocks []byte
	for _, v := range fsMock.vBufBlocks {
		vBlocks = append(vBlocks, v...)
	}
	for _, v := range byte2D {
		assert.Equal(t, v, vBlocks[counter:counter+len(v)])
		counter = counter + len(v)
	}

	for _, v := range byte2D {
		buf := make([]byte, len(v))
		for j := range buf {
			buf[j] = 0
		}

		//if i+1 == len(byte2D) {
		//	log.Warn("last packet")
		//}
		_, err := vf.Read(buf)
		assert.Equal(t, nil, err)
		if err != nil {
			log.Warn("Test")
			break
		}
		// assert.Equal(t, v[0], buf[0])
		assert.Equal(t, 0, bytes.Compare(v, buf))

	}

	assert.Equal(t, 0, bytes.Compare(vBlocks, vf.bufRX))
}

func TestIO_ReadAt(t *testing.T) {
	fsMock := NewVBufMock(t)
	blm := blockAllocationMap.New(log.GetScope("test"), fsMock, maxNumberOfBlocks)
	vf := NewVirtualFile("test", vfID, blockSizeTest, fsMock, blm,
		int(blockSizeTest)*2, log.GetScope("test2"))
	fsMock.openFiles[vfID] = vf

	size := 0
	VFSize := int(1.5 * blockSizeTest)
	MaxByteArraySize := int(blockSizeTest * 0.5)
	for {
		token := make([]byte, uint32(rand.Intn(MaxByteArraySize))+1)
		m, err := rand.Read(token)
		assert.Equal(t, nil, err)
		byte2D = append(byte2D, token)
		assert.Equal(t, m, len(token))
		size = size + m
		n, err := vf.Write(token)
		assert.Equal(t, nil, err)
		assert.Equal(t, m, n)

		if size > VFSize {
			break
		}
	}

	err := vf.Close()
	assert.Equal(t, nil, err)

	assert.Equal(t, size, len(fsMock.vBuf))
	counter := 0
	var vBlocks []byte
	for _, v := range fsMock.vBufBlocks {
		vBlocks = append(vBlocks, v...)
	}
	for _, v := range byte2D {
		// ToDo: fix this
		// assert.Equal(t, v, vBlocks[counter:counter+len(v)])
		counter = counter + len(v)
	}

	numberOfTest := 5

	for i := 0; i < numberOfTest; i++ {
		data := make([]byte, uint32(rand.Intn(blockSizeTest*0.3))+1)
		assert.Equal(t, nil, err)
		offset := rand.Intn(int(float32(size) * 0.7))
		n, err := vf.ReadAt(data, int64(offset))

		assert.Equal(t, nil, err)
		assert.Equal(t, len(data), n)
		assert.Equal(t, vBlocks[offset:offset+len(data)], data)
	}
}
