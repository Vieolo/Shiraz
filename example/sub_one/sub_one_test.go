package subone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFnThree(t *testing.T) {
	assert.Equal(t, FnThee(), "three")
	assert.Equal(t, FnThee(), "threes")
}

func TestFnThreeWithFail(t *testing.T) {
	assert.Equal(t, FnThee(), "three")
	assert.Equal(t, FnThee(), "threes")
}

func TestFnThreeAnother(t *testing.T) {
	assert.Equal(t, FnThee(), "three")
}
