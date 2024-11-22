package bttcp

import (
	"container/list"
	"net"
	"sync"
)

type Pool struct {
	maxSize     int
	connections *list.List
	addr        string
	lock        sync.Mutex
}

func NewPool(maxSize int, addr string) *Pool {
	return &Pool{
		maxSize:     maxSize,
		connections: list.New(),
		addr:        addr,
	}
}

func (p *Pool) GetConnection() (net.Conn, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.connections.Len() > 0 {
		e := p.connections.Front()
		p.connections.Remove(e)
		connFromPool := e.Value.(net.Conn)
		_, err := connFromPool.Write(nil)
		if err != nil {
			conn, err := net.Dial("tcp", p.addr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		return e.Value.(net.Conn), nil
	}

	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p *Pool) ReleaseConnection(conn net.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.connections.Len() < p.maxSize {
		p.connections.PushBack(conn)
	} else {
		err := conn.Close()
		if err != nil {
			return
		}
	}
}
