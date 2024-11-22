package bttcp

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer(HandlerFunc(func(b []byte) ([]byte, error) {
		return b, nil
	}))
	s.Run(9999)
}
