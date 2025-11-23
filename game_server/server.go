package game_server

import (
	"fmt"
	"goTibia/packets/game"
	"goTibia/protocol"
	"log"
	"net"
	"time"
)

type Server struct {
	ListenAddr     string
	RealServerAddr string
}

func NewServer(listenAddr, realServerAddr string) *Server {
	return &Server{
		ListenAddr:     listenAddr,
		RealServerAddr: realServerAddr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Game proxy listening on %s, forwarding to %s", s.ListenAddr, s.RealServerAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Game proxy failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	protoClientConn := protocol.NewConnection(clientConn)
	defer protoClientConn.Close()
	log.Printf("Game: Accepted connection from %s", protoClientConn.RemoteAddr())

	packetReader, err := protoClientConn.ReadMessage()
	if err != nil {
		log.Printf("Game: error reading message from %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	_, err = game.ParseLoginRequest(packetReader)
	if err != nil {
		log.Printf("Login: Failed to parse login packet: %v", err)
		return
	}

	protoServerConn, err := s.connectToServer()
	if err != nil {
		log.Printf("Login: Failed to connect to %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

}

func (s *Server) connectToServer() (*protocol.Connection, error) {
	conn, err := net.DialTimeout("tcp", s.RealServerAddr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("backend unavailable at %s: %w", s.RealServerAddr, err)
	}

	return protocol.NewConnection(conn), nil
}
