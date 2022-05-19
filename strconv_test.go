package sly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS2B(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := "hello"
		assert.Equal(t, []byte("hello"), S2B(s))
	})
	t.Run("empty", func(t *testing.T) {
		s := ""
		assert.Equal(t, []byte(nil), S2B(s))
	})
}

func TestB2S(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		b := []byte("hello")
		assert.Equal(t, "hello", B2S(b))
	})
	t.Run("nil", func(t *testing.T) {
		b := []byte(nil)
		assert.Equal(t, "", B2S(b))
	})
}
