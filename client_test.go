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
	s := NewClient("127.0.0.1:9999", 2048)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		num := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(1)) * time.Second)
			resp, err := s.Send([]byte("hello world " + strconv.Itoa(num)))
			if err != nil {
				wg.Done()
				return
			}
			fmt.Println(string(resp))
			wg.Done()
		}()
	}
	wg.Wait()
}
