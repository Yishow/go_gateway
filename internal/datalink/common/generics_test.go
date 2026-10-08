package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	str := "test"
	strPtr := Ptr(str)
	assert.NotNil(t, strPtr)
	assert.Equal(t, "test", *strPtr)

	num := 42
	numPtr := Ptr(num)
	assert.NotNil(t, numPtr)
	assert.Equal(t, 42, *numPtr)

	b := true
	bPtr := Ptr(b)
	assert.NotNil(t, bPtr)
	assert.True(t, *bPtr)
}

func TestDeref(t *testing.T) {
	str := "hello"
	assert.Equal(t, "hello", Deref(&str, "fallback"))
	assert.Equal(t, "fallback", Deref(nil, "fallback"))

	num := 100
	assert.Equal(t, 100, Deref(&num, -1))
	assert.Equal(t, -1, Deref(nil, -1))
}

func TestDerefOrZero(t *testing.T) {
	str := "world"
	assert.Equal(t, "world", DerefOrZero(&str))
	assert.Equal(t, "", DerefOrZero[string](nil))

	num := 99
	assert.Equal(t, 99, DerefOrZero(&num))
	assert.Equal(t, 0, DerefOrZero[int](nil))
}
