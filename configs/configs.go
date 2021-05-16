package configs

import "sync"

var configs = struct {
	MaxSizeVirtualFile uint32
	mutex sync.Mutex
}{
	MaxSizeVirtualFile: 10*1024*1024*1024,  // 10 GB
}

func SetConfigs(maxSizeVirtualFile uint32) {
	configs.mutex.Lock()
	defer configs.mutex.Unlock()
	configs.MaxSizeVirtualFile = maxSizeVirtualFile
}

func GetMaxSizeVirtualFile() uint32 {
	configs.mutex.Lock()
	defer configs.mutex.Unlock()
	return configs.MaxSizeVirtualFile
}


