package proxy

import (
	"fmt"
	"goTibia/protocol"
	"log"
	"net"
	"time"
)

type ConnectionHandler interface {
	Handle(client *protocol.Connection)
}

type Server struct {
	Name       string
	ListenAddr string
	Handler    ConnectionHandler
}

func NewServer(name, listenAddr string, handler ConnectionHandler) *Server {
	return &Server{
		Name:       name,
		ListenAddr: listenAddr,
		Handler:    handler,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("%s proxy failed to bind %s: %w", s.Name, s.ListenAddr, err)
	}
	defer listener.Close()

	log.Printf("[%s] Proxy listening on %s", s.Name, s.ListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] Accept error: %v", s.Name, err)
			continue
		}

		// Wrap connection immediately
		protoConn := protocol.NewConnection(conn)

		// Hand off to the specific logic in a goroutine
		go func() {
			defer protoConn.Close()
			s.Handler.Handle(protoConn)
		}()
	}
}

func ConnectToBackend(address string) (*protocol.Connection, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("backend unavailable at %s: %w", address, err)
	}
	return protocol.NewConnection(conn), nil
}
