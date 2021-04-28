package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// fileExists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// directoryExists
func DirectoryExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func WalkMatch(root, pattern string) ([]int, error) {
	var matches []int
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := regexp.MatchString(pattern, filepath.Base(filePath)); err != nil {
			return err
		} else if matched {
			val, _ := strconv.Atoi(filepath.Base(filePath))
			matches = append(matches, val)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func MinMax(array []int) (min int, max int) {
	max = array[0]
	min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func LatestFile(path string) (max int) {
	files, err := WalkMatch(path, "^[0-9]*$")
	if err != nil {
		return 0
	}
	if len(files) > 0 {
		_, max = MinMax(files)
	}
	return
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func Normalize(val uint64, percent float64, min uint64, max uint64) uint64 {
	v := uint64(float64(val) * percent)
	if v > max {
		return max
	} else if v < min {
		return min
	} else {
		return v
	}
}

// OpenFile is a modified version of os.OpenFile which sets O_DIRECT
func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	// syscall.O_DIRECT|
	return os.OpenFile(name, flag, perm)
}

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func DeleteFile(filename string) error {
	e := os.Remove(filename)
	return e
}
