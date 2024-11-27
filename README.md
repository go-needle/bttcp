<!-- markdownlint-disable MD033 MD041 -->
<div align="center">

# 🪡bttcp

<!-- prettier-ignore-start -->
<!-- markdownlint-disable-next-line MD036 -->
Binary Interaction Protocol Based on TCP
<!-- prettier-ignore-end -->

<img src="https://img.shields.io/badge/golang-1.11+-blue" alt="golang">
</div>

## introduction
bttcp is a binary transmission protocol based on TCP, where the response header only records the length of the binary stream and data exchange is limited to byte exchange. bttcp is designed to improve transmission performance.

## installing
Select the version to install

`go get github.com/go-needle/bttcp@version`

If you have already get , you may need to update to the latest version

`go get -u github.com/go-needle/bttcp`


## quickly start

### server code
```golang
package main

import "github.com/go-needle/bttcp"

func main() {
	s := bttcp.NewServer(bttcp.HandlerFunc(func(b []byte) []byte{
		return b
	}))
	s.Run(9999)
}
```

### client code
```golang
package main

import (
	"fmt"
	"github.com/go-needle/bttcp"
	"math/rand"
	"strconv"
	"sync"
	"time"
)
var wg sync.WaitGroup

func main() {
	s := bttcp.NewClient("127.0.0.1:9999", 100, true)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		num := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
			resp, err := s.Send([]byte("hello" + strconv.Itoa(num)))
			if err != nil {
				return
			}
			fmt.Println(string(resp))
			wg.Done()
		}()
	}
	wg.Wait()
}
```