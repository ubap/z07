package login

import (
	"fmt"
	"goTibia/protocol"
	"goTibia/protocol/login"
	"goTibia/proxy"
	"goTibia/proxy/login/handlers"
	"io"
	"log"
	"net"
)

type Server struct {
	ListenAddr     string
	RealServerAddr string
	// You could add other dependencies here, like a specific logger.
}

func NewServer(listenAddr, realServerAddr string) *Server {
	// Create the default handler instance here in the application layer.
	defaultHandler := &handlers.DefaultPassThroughHandler{}

	// Pass it to the registry constructor.
	registry := protocol.NewHandlerRegistry(defaultHandler)
	registry.Register(login.ServerOpcodeMOTD, &handlers.MOTDHandler{})

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
	log.Printf("Login proxy listening on %s, forwarding to %s", s.ListenAddr, s.RealServerAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Login proxy failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	protoClientConn := protocol.NewConnection(clientConn)
	defer protoClientConn.Close()
	log.Printf("Login: Accepted connection from %s", protoClientConn.RemoteAddr())

	// Phase 1: Receive and Decode
	loginPacket, err := s.receiveLoginPacket(protoClientConn)
	if err != nil {
		log.Printf("Login: Failed to process login packet from %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	// Phase 2: Forward and Re-encode
	protoServerConn, err := s.forwardLoginPacket(loginPacket)
	if err != nil {
		log.Printf("Login: Failed to forward login packet for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	// Phase 3: Bridge the connection
	s.bridgeConnections(protoClientConn, protoServerConn)

	log.Printf("Login: Connection for %s finished.", protoClientConn.RemoteAddr())
}

func (s *Server) receiveLoginPacket(client *protocol.Connection) (*login.LoginPacket, error) {
	messageBytes, err := client.ReadMessage()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("client disconnected before sending login packet")
		}
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	packet, err := login.ParseLoginPacket(messageBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing login packet: %w", err)
	}

	log.Printf("Login: Successfully decrypted packet: Account=%d, Password='%s'",
		packet.AccountNumber, packet.Password)
	return packet, nil
}

// forwardLoginPacket connects to the real server and sends the re-encoded packet.
// It returns the established server connection or an error.
func (s *Server) forwardLoginPacket(packet *login.LoginPacket) (*protocol.Connection, error) {
	serverConn, err := net.Dial("tcp", s.RealServerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to real server at %s: %w", s.RealServerAddr, err)
	}
	protoServerConn := protocol.NewConnection(serverConn)
	log.Printf("Login: Successfully connected to real server %s", s.RealServerAddr)

	outgoingMessageBytes, err := packet.Marshal()
	if err != nil {
		protoServerConn.Close() // Clean up connection on failure
		return nil, fmt.Errorf("failed to marshal outgoing packet: %w", err)
	}

	if err := protoServerConn.WriteMessage(outgoingMessageBytes); err != nil {
		protoServerConn.Close() // Clean up connection on failure
		return nil, fmt.Errorf("failed to send login packet to real server: %w", err)
	}

	log.Println("Login: Sent re-encrypted login packet to real server.")
	return protoServerConn, nil
}

// bridgeConnections handles shuttling data between the client and server.
// For the login server, this is a simple one-way copy from server to client.
func (s *Server) bridgeConnections(client *protocol.Connection, server *protocol.Connection) {
	log.Println("Login: Bridging server response to client...")

	// For the login server, we only need to copy the response from the server
	// back to the client. There is no further client->server communication
	// on this specific connection.

	// Create our hex dumper with a clear label.
	dumper := &proxy.HexDumpWriter{Prefix: "SERVER -> CLIENT"}
	teeReader := io.TeeReader(server.RawConn(), dumper)

	// Copy the server's response (e.g., character list) to the client.
	bytesCopied, err := io.Copy(client.RawConn(), teeReader)
	if err != nil {
		log.Printf("Login: Error during bridge copy: %v", err)
	}

	log.Printf("Login: Bridged %d bytes from server to client.", bytesCopied)
}
