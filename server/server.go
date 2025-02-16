package server

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
)

type Server struct {
	backends    []string
	port        int
	connections map[string]int
	mutex       sync.Mutex
}

func NewServer(backends []string, port int) *Server {
	return &Server{
		backends: backends,
		port:     port,
	}
}

func (s *Server) Listen() {
	if s.connections == nil {
		s.connections = make(map[string]int)
	}
	listener, err := net.Listen("tcp", ":8080")
	fmt.Println("Listening on :" + fmt.Sprint(s.port))
	if err != nil {
		panic(err)
	}
	for {
		connection, err := listener.Accept()
		fmt.Println("Accepted connection: ", connection.RemoteAddr())
		if err != nil {
			panic(err)
		}
		go s.handle(connection)
	}
}

func (s *Server) decreaseConnection(backend string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connections[backend]--
	fmt.Println("Decreased connection count for backend: ", backend)

}

func (s *Server) handle(connection net.Conn) {
	backend := s.selectBackend() // or s.randomBackend()
	defer s.decreaseConnection(backend)
	backendConnection, err := net.Dial("tcp", backend)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to backend: ", backend)
	go pipe(connection, backendConnection)
	go pipe(backendConnection, connection)
}

func (s *Server) randomBackend() string {
	index := rand.Int() % len(s.backends)
	return s.backends[index]
}

func (s *Server) selectBackend() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.connections) == 0 {
		s.connections[s.backends[0]] = 1
		return s.backends[0]
	}
	min := s.connections[s.backends[0]]
	minBackend := s.backends[0]
	for _, backend := range s.backends {
		if s.connections[backend] < min {
			min = s.connections[backend]
			minBackend = backend
		}
	}

	s.connections[minBackend]++

	return minBackend
}

func pipe(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			return
		}
		_, err = dst.Write(buf[:n])
		if err != nil {
			return
		}
	}
}
