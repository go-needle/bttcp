package bttcp

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestClient(t *testing.T) {
	s := NewClient("127.0.0.1:9999", 100)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
			resp, err := s.Send([]byte("hello" + strconv.Itoa(i)))
			if err != nil {
				return
			}
			fmt.Println(string(resp))
			wg.Done()
		}()
	}
	wg.Wait()
}
