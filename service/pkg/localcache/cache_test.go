package localcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValue struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func TestSetAndGet(t *testing.T) {
	cache, _ := NewDefaultClient()
	t1 := testValue{
		Name: "foo",
		Age:  18,
	}
	cache.Set("test", t1)
	t2 := testValue{}
	cache.Get("test", &t2)
	assert.Equal(t, t1, t2, "The two value should be the same.")
}
