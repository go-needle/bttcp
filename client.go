package bttcp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-needle/bttcp/proto"
	"net"
	"sync"
)

type Client struct {
	address  string
	poolSize int
	pool     *Pool
	once     sync.Once
}

func NewClient(address string, poolSize int, isTestConn bool) *Client {
	if isTestConn {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			panic(errors.New(fmt.Sprintf("dial bttcp %s: connect: connection refused", address)))
		}
		_, err = conn.Write(nil)
		if err != nil {
			panic(errors.New(fmt.Sprintf("dial bttcp %s: connect: connection refused", address)))
		}
		err = conn.Close()
		if err != nil {
			panic(errors.New(fmt.Sprintf("dial bttcp %s: connect: connection refused", address)))
		}
	}
	return &Client{address: address, poolSize: poolSize}
}

func (c *Client) Send(b []byte) ([]byte, error) {
	c.once.Do(func() {
		c.pool = NewPool(c.poolSize, c.address)
	})
	conn, err := c.pool.GetConnection()
	defer c.pool.ReleaseConnection(conn)
	if err != nil {
		return nil, err
	}
	data, err := proto.Encode(b)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	rb, err := proto.Decode(reader)
	if err != nil {
		return nil, err
	}
	return rb, nil
}

func (c *Client) Close() {
	c.pool.ClearPool()
}
