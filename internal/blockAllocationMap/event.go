package blockAllocationMap

type event interface {
	NoSpace() uint32
}
