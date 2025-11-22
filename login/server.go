package login

import (
	"fmt"
	"goTibia/protocol"
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

	packetReader, err := protoClientConn.ReadMessage()
	if err != nil {
		log.Printf("error reading message from %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	loginPacket, err := ParseCredentialsPacket(packetReader)
	if err != nil {
		log.Printf("Login: Failed to parse login packet: %v", err)
		return
	}

	// TODO Rename, struct more meaningfully
	protoServerConn, err := s.forwardLoginPacket(loginPacket)
	if err != nil {
		log.Printf("Login: Failed to forward credentials packet for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	protoServerConn.EnableXTEA(loginPacket.XTEAKey)
	protoClientConn.EnableXTEA(loginPacket.XTEAKey)

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("Login: Failed to read server response for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	resultMessage, err := s.receiveLoginResultMessage(message)
	if err != nil {
		log.Printf("Login: Failed to receive login result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	marshal, err := resultMessage.Marshal()
	if err != nil {
		log.Printf("Login: Failed to marshal result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	protoClientConn.WriteMessage(marshal)

	log.Printf("Login: Connection for %s finished.", protoClientConn.RemoteAddr())
}

// forwardLoginPacket connects to the real server and sends the re-encoded packet.
// It returns the established server connection or an error.
func (s *Server) forwardLoginPacket(packet *ClientCredentialPacket) (*protocol.Connection, error) {
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
func (s *Server) receiveLoginResultMessage(packetReader *protocol.PacketReader) (*LoginResultMessage, error) {
	message := LoginResultMessage{}

	// --- 2. Loop until the command stream is empty. ---
	for {
		// Read the next opcode.
		opcode := packetReader.ReadByte()
		err := packetReader.Err()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		log.Printf("Login: Processing opcode %#x", opcode)

		switch opcode {
		case S2COpcodeDisconnectClient:
			disconnectedReason := packetReader.ReadString()
			log.Print("DisconnectClientHandler: " + disconnectedReason)
			message.ClientDisconnected = true
			message.ClientDisconnectedReason = disconnectedReason
		case S2COpcodeMOTD:
			motd, err := ReadMotd(packetReader)
			if err != nil {
				return nil, err
			}
			message.Motd = motd
		case S2COpcodeCharacterList:
			charList, err := ReadCharacterList(packetReader)
			if err != nil {
				return nil, err
			}
			message.CharacterList = charList
		default:
			panic("unknown opcode " + fmt.Sprintf("%#x", opcode))
		}

	}
	return &message, nil
}
