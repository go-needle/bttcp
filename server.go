package bttcp

import (
	"bufio"
	"fmt"
	"github.com/go-needle/bttcp/proto"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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
	handler   Handler
	whiteList map[string]struct{}
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

// Run will listen to TCP ports and block
func (s *Server) Run(port int) {
	welcome()
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	fmt.Println("🪡 The bttcp server is listening at port " + strconv.Itoa(port))
	defer func(listen net.Listener) {
		err = listen.Close()
		if err != nil {
			panic(err)
		}
	}(listen)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("accept failed, err: ", err)
			continue
		}
		if s.whiteList != nil {
			sp := strings.Split(conn.RemoteAddr().String(), ":")
			if len(sp) > 1 {
				ip := sp[0]
				if _, has := s.whiteList[ip]; !has {
					err := conn.Close()
					if err != nil {
						continue
					}
					continue
				}
			}
		}
		log.Printf("a new bttcp connection links from %s -> %s", conn.RemoteAddr().String(), conn.LocalAddr().String())
		go s.process(conn)
	}
}

// SetWhiteList is to set an IP whitelist. Setting an IP whitelist can reject IP connections that are not on the whitelist. If it is not set, it will not be enabled
func (s *Server) SetWhiteList(remoteIPs ...string) {
	if s.whiteList == nil {
		s.whiteList = make(map[string]struct{})
	}
	for _, rd := range remoteIPs {
		s.whiteList[rd] = struct{}{}
	}
}

func (s *Server) process(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recover panic : ", err)
		}
	}()
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
