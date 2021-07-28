package fsEngine

import (
	"github.com/fanap-infra/fsEngine/internal/blockAllocationMap"
	"github.com/fanap-infra/fsEngine/pkg/virtualFile"
)

type VFInfo struct {
	id             uint32
	vfs            []*virtualFile.VirtualFile
	blm            *blockAllocationMap.BlockAllocationMap
	numberOfOpened uint32
}
