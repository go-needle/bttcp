package bttcp

import (
	"bufio"
	"github.com/go-needle/bttcp/proto"
	"sync"
)

type Client struct {
	address  string
	poolSize int
	pool     *Pool
	once     sync.Once
}

func NewClient(address string, poolSize int) *Client {
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
