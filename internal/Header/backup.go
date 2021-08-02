package Header_

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/fanap-infra/fsEngine/pkg/utils"

	"github.com/fanap-infra/fsEngine/internal/constants"
)

func (hfs *HFileSystem) backUp() error {
	backUpPath := hfs.path + "/" + constants.HeaderBackUpPath
	out, err := utils.OpenFile(backUpPath, os.O_CREATE|os.O_RDWR, 0o777)
	if err != nil {
		hfs.log.Errorv("can not create back up file",
			"backUpPath", backUpPath, "err", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, hfs.file)
	if err != nil {
		hfs.log.Errorv("can not copy file to back ups",
			"backUpPath", backUpPath, "err", err.Error())
		return err
	}
	return nil
}

func (hfs *HFileSystem) loadBackUp() error {
	backUpPath := hfs.path + "/" + constants.HeaderBackUpPath
	out, err := utils.OpenFile(backUpPath, os.O_RDWR, 0o777)
	if err != nil {
		hfs.log.Errorv("can not create back up file",
			"backUpPath", backUpPath, "err", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(hfs.file, out)
	if err != nil {
		hfs.log.Errorv("can not copy file to back ups",
			"backUpPath", backUpPath, "err", err.Error())
		return err
	}
	return nil
}

func (hfs *HFileSystem) updateHash() error {
	hashMaker := sha256.New()

	b := make([]byte, HashByteIndex)
	n, err := hfs.file.ReadAt(b, 0)
	if err != nil {
		hfs.log.Errorv("can not read file to hash",
			"n", n, "err", err.Error())
		return err
	}
	if HashByteIndex != n {
		hfs.log.Errorv("can not read file completely to hash",
			"n", n, "HashByteIndex", HashByteIndex)
		return fmt.Errorf("can not read completely to hash ")
	}
	n, err = hashMaker.Write(b)
	if err != nil {
		hfs.log.Errorv("can not write to hash writer",
			"n", n, "err", err.Error())
		return err
	}
	if HashByteIndex != n {
		hfs.log.Errorv("can not write completely to hash writer",
			"n", n, "HashByteIndex", HashByteIndex)
		return fmt.Errorf("can not write completely to hash writer")
	}
	hash := hashMaker.Sum(nil)
	n, err = hfs.file.WriteAt(hash, HashByteIndex)
	if err != nil {
		hfs.log.Errorv("can not write to file",
			"n", n, "err", err.Error())
		return err
	}
	if len(hash) != n {
		hfs.log.Errorv("can not write completely to file",
			"n", n, "HashByteIndex", HashByteIndex)
		return fmt.Errorf("can not write completely to hash file")
	}
	return nil
}

func (hfs *HFileSystem) checkHash() bool {
	hashMaker := sha256.New()

	b := make([]byte, HashByteIndex)
	n, err := hfs.file.ReadAt(b, 0)
	if err != nil {
		hfs.log.Errorv("can not read file for hash checking",
			"n", n, "err", err.Error())
		return false
	}
	if HashByteIndex != n {
		hfs.log.Errorv("can not read file completely for hashing",
			"n", n, "HashByteIndex", HashByteIndex)
		return false
	}
	n, err = hashMaker.Write(b)
	if err != nil {
		hfs.log.Errorv("can not write to hash writer",
			"n", n, "err", err.Error())
		return false
	}
	if HashByteIndex != n {
		hfs.log.Errorv("can not write completely to hash writer",
			"n", n, "HashByteIndex", HashByteIndex)
		return false
	}
	hash := hashMaker.Sum(nil)

	hashValue := make([]byte, HashSize)
	n, err = hfs.file.ReadAt(hashValue, HashByteIndex)
	if err != nil {
		hfs.log.Errorv("can not read hash value",
			"n", n, "err", err.Error())
		return false
	}
	if HashSize != n {
		hfs.log.Errorv("can not read hash value completely",
			"n", n, "HashByteIndex", HashByteIndex)
		return false
	}
	return bytes.Equal(hash, hashValue)
}
