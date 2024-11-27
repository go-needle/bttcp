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
	c := NewClient("127.0.0.1:9999", 2048, true)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		num := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			resp, err := c.Send([]byte("hello world " + strconv.Itoa(num)))
			if err != nil {
				wg.Done()
				return
			}
			fmt.Println(string(resp))
			wg.Done()
		}()
	}
	wg.Wait()
	c.Close()
}

func TestSend(t *testing.T) {
	c := NewClient("127.0.0.1:9999", 2048, true)
	res, err := c.Send([]byte("hello world 1"))
	if err != nil {
		return
	}
	fmt.Println(string(res))
	time.Sleep(time.Duration(5) * time.Second)
	res, err = c.Send([]byte("hello world 2"))
	if err != nil {
		return
	}
	fmt.Println(string(res))
}
