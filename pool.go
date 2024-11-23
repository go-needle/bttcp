package bttcp

import (
	"net"
	"sync"
)

type node struct {
	val  net.Conn
	prev *node
	next *node
}

type Pool struct {
	maxSize int
	curSize int
	head    *node
	tail    *node
	addr    string
	lock    sync.Mutex
}

func NewPool(maxSize int, addr string) *Pool {
	headDummy := new(node)
	tailDummy := new(node)
	headDummy.next = tailDummy
	tailDummy.prev = headDummy
	return &Pool{
		maxSize: maxSize,
		head:    headDummy,
		tail:    tailDummy,
		addr:    addr,
	}
}

func (p *Pool) GetConnection() (net.Conn, error) {
	p.lock.Lock()
	if p.curSize > 0 {
		e := p.tail.next
		p.tail.next = e.next
		p.lock.Unlock()
		connFromPool := e.val
		_, err := connFromPool.Write(nil)
		if err != nil {
			connFromPool.Close()
			conn, err := net.Dial("tcp", p.addr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		}
		return e.val, nil
	}
	p.lock.Unlock()
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p *Pool) ReleaseConnection(conn net.Conn) {
	p.lock.Lock()
	if p.curSize < p.maxSize {
		e := &node{val: conn}
		p.tail.prev.next = e
		p.tail.prev = e
		p.lock.Unlock()
	} else {
		p.lock.Unlock()
		err := conn.Close()
		if err != nil {
			return
		}
	}
}
