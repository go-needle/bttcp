package bttcp

import (
	"bufio"
	"fmt"
	"github.com/go-needle/bttcp/proto"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

// Handler defines the request handler
type Handler interface {
	Handle(b []byte) []byte
}

// HandlerFunc realizes the Handler
type HandlerFunc func(b []byte) []byte

func (f HandlerFunc) Handle(b []byte) []byte {
	return f(b)
}

type Server struct {
	handler Handler
}

func NewServer(handler Handler) *Server {
	return &Server{handler: handler}
}

func getInternalIP() (string, error) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range adders {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no internal IP address found, check for multiple interfaces")
}

func welcome() {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("🪡 Welcome to use go-needle-bttcp")
	ip, err := getInternalIP()
	if err == nil {
		fmt.Println("🪡 IP address: " + ip)
	}
}

func (s *Server) Run(port int) {
	welcome()
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	fmt.Println("🪡 The http server is listening at port " + strconv.Itoa(port))
	defer func(listen net.Listener) {
		err = listen.Close()
		if err != nil {
			panic(err)
		}
	}(listen)
	for {
		conn, err := listen.Accept()
		log.Printf("a new connection links from %s -> %s", conn.RemoteAddr().String(), conn.LocalAddr().String())
		if err != nil {
			log.Println("accept failed, err: ", err)
			continue
		}
		go s.process(conn)
	}
}

func (s *Server) process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	reader := bufio.NewReader(conn)
	for {
		b, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("decode msg failed, err: ", err)
			return
		}
		wb := s.handler.Handle(b)
		data, err := proto.Encode(wb)
		if err != nil {
			log.Println("encode msg failed, err: ", err)
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Println("write msg failed, err: ", err)
			return
		}
	}
}
