package seq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUint64(t *testing.T) {
	value, err := NewUint64()
	assert.NoError(t, err)
	assert.NotZero(t, value)
}

func TestNewString(t *testing.T) {
	value, err := NewString()
	assert.NoError(t, err)
	assert.NotZero(t, value)
}

func TestNewUint(t *testing.T) {
	value, err := NewUint()
	assert.NoError(t, err)
	assert.NotZero(t, value)
}
