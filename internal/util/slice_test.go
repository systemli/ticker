package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/systemli/ticker/internal/util"
)

func TestContains(t *testing.T) {
	slice := []int{1, 2}

	assert.True(t, Contains(slice, 1))
	assert.False(t, Contains(slice, 3))
}

func TestAppend(t *testing.T) {
	slice := []int{1, 2}

	assert.Equal(t, slice, Append(slice, 2))
	assert.Equal(t, []int{1, 2, 3}, Append(slice, 3))
}

func TestRemove(t *testing.T) {
	slice := []int{1, 2}

	assert.Equal(t, slice, Remove(slice, 3))
	assert.Equal(t, []int{1}, Remove(slice, 2))
	assert.Equal(t, []int{}, Remove([]int{}, 2))
}

func TestContainsString(t *testing.T) {
	slice := []string{"a", "b"}

	assert.True(t, ContainsString(slice, "a"))
	assert.False(t, ContainsString(slice, "c"))
}
