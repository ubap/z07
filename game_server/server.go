package game_server

import (
	"goTibia/packets/game"
	"goTibia/protocol"
	"log"
	"net"
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

}
