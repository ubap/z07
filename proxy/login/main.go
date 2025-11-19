package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"goTibia/protocol"
	"io"
	"log"
	"net"
)

// The address and port for our dummy server to listen on.
const listenAddr = ":7171"

func main() {
	// Start listening for incoming TCP connections on the specified address.
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Ensure the listener is closed when the main function exits.
	defer listener.Close()

	log.Printf("Dummy server listening on %s. Waiting for Tibia client to connect...", listenAddr)

	// Loop forever, accepting new connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	// Tibia packets are prefixed with a 2-byte little-endian length.
	// This length specifies the size of the payload that follows.
	var packetLength uint16
	err := binary.Read(conn, binary.LittleEndian, &packetLength)
	if err != nil {
		if err == io.EOF {
			log.Println("Client disconnected before sending data.")
		} else {
			log.Printf("Error reading packet length: %v", err)
		}
		return
	}

	log.Printf("Packet length header received: %d bytes", packetLength)

	// Create a buffer to hold the rest of the packet.
	packetBody := make([]byte, packetLength)

	// Read the exact number of bytes specified by the length header.
	// io.ReadFull is used to ensure we get all the data or an error.
	_, err = io.ReadFull(conn, packetBody)
	if err != nil {
		log.Printf("Error reading packet body: %v", err)
		return
	}

	packet, err := protocol.ParseLoginPacket(packetBody)
	if err != nil {
		log.Printf("Error parsing login packet: %v", err)
		return
	}

	log.Printf("Packet received: %v", packet)

	// --- The most important part for reverse engineering ---
	// Print a detailed hex dump of the packet's body.
	fmt.Printf("\n--- Packet Received from %s ---\n", conn.RemoteAddr())
	fmt.Printf("%s", hex.Dump(packetBody)) // hex.Dump provides a beautiful, formatted output.
	fmt.Println("--- End of Packet ---")

	PrintAsGoSlice(packetBody)
}
