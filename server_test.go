package bttcp

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(HandlerFunc(func(b []byte) []byte {
		return b
	}))
	s.Run(9999)
}
