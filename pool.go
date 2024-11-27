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
		e := p.head.next
		p.head.next = e.next
		p.head.next.prev = p.head
		p.curSize--
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
		return connFromPool, nil
	}
	p.lock.Unlock()
	conn, err := net.Dial("tcp", p.addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (p *Pool) ReleaseConnection(conn net.Conn) {
	if conn == nil {
		return
	}
	p.lock.Lock()
	if p.curSize < p.maxSize {
		e := &node{val: conn, prev: p.tail.prev, next: p.tail}
		p.tail.prev.next = e
		p.tail.prev = e
		p.curSize++
		p.lock.Unlock()
	} else {
		p.lock.Unlock()
		err := conn.Close()
		if err != nil {
			return
		}
	}
}

func (p *Pool) ClearPool() {
	cur := p.head
	for cur != p.tail {
		conn := cur.val
		if conn != nil {
			err := conn.Close()
			if err != nil {
				continue
			}
		}
		cur = cur.next
	}
	p.head.next = p.tail
	p.tail.prev = p.head
}
