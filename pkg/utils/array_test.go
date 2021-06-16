package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemExists(t *testing.T) {
	arr := []uint{2, 5, 14, 17}
	assert.Equal(t, false, ItemExists(arr, uint(1)))
	assert.Equal(t, true, ItemExists(arr, uint(14)))

	assert.Equal(t, false, ItemExists(arr, uint(15)))
	assert.Equal(t, true, ItemExists(arr, uint(5)))
}
