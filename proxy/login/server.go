package login

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goTibia/protocol"
	"goTibia/protocol/crypto"
	"goTibia/protocol/login"
	"goTibia/proxy"
	"goTibia/proxy/login/handlers"
	"io"
	"log"
	"net"
)

type Server struct {
	ListenAddr      string
	RealServerAddr  string
	HandlerRegistry *protocol.HandlerRegistry
	// You could add other dependencies here, like a specific logger.
}

func NewServer(listenAddr, realServerAddr string) *Server {
	registry := protocol.NewHandlerRegistry()
	registry.Register(login.ServerOpcodeDisconnectClient, &handlers.DisconnectClientHandler{})
	registry.Register(login.ServerOpcodeMOTD, &handlers.MOTDHandler{})

	return &Server{
		ListenAddr:      listenAddr,
		RealServerAddr:  realServerAddr,
		HandlerRegistry: registry,
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

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		return
	}

	decrypted, err := crypto.DecryptXTEA(message, loginPacket.XTEAKey)
	if err != nil {
		return
	}

	dumper := &proxy.HexDumpWriter{Prefix: "SERVER -> CLIENT"}
	dumper.Write(decrypted)
	// 2 byte - message length
	// 1 byte - opcode

	s.processStream(decrypted)

	message, err = crypto.EncryptXTEA(decrypted, loginPacket.XTEAKey)
	if err != nil {
		return
	}
	protoClientConn.WriteMessage(message)

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

/*
*
[ 2-byte Inner Length | Opcode 1 | Data 1 | Opcode 2 | Data 2 | ... | Padding ]

	\___________________/ \_____________________________________/
	        |                             |
	     (Header)                 (Stream of Commands)
*/
func (s *Server) processStream(decryptedPayload []byte) ([]byte, error) {
	// A reader for the incoming decrypted stream.
	streamReader := bytes.NewReader(decryptedPayload)

	// --- 1. Read the single length header at the beginning of the stream. ---
	var streamLength uint16
	if err := binary.Read(streamReader, binary.LittleEndian, &streamLength); err != nil {
		return nil, fmt.Errorf("error reading stream length header: %w", err)
	}

	// Create a new reader that is limited to reading only the command stream.
	commandStream := io.LimitReader(streamReader, int64(streamLength))

	// --- 2. Loop until the command stream is empty. ---
	for {
		// Read the next opcode.
		opcodeBuffer := make([]byte, 1)
		n, err := commandStream.Read(opcodeBuffer)
		if err == io.EOF {
			break // Successfully reached the end of the stream.
		}
		if err != nil {
			return nil, fmt.Errorf("error reading opcode from stream: %w", err)
		}
		if n == 0 {
			break
		}
		opcode := opcodeBuffer[0]
		log.Printf("Login: Processing opcode %#x", opcode)

		handler, err := s.HandlerRegistry.Get(opcode)
		if err != nil {
			log.Printf("Login: Failed to get handler for opcode %#x, short-circuiting", opcode)
			return decryptedPayload, err
		}

		err = handler.Handle(commandStream)
		if err != nil {
			log.Printf("Login: Handler for opcode %d returned error: %v", opcode, err)
			return decryptedPayload, err
		}

	}

	return decryptedPayload, nil
}
