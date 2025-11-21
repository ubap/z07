package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ReadString(r io.Reader) (string, error) {
	// 1. Read the 2-byte length prefix.
	var length uint16
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		// If we can't even read the length, return a clear error.
		return "", fmt.Errorf("could not read string length: %w", err)
	}

	// 2. Handle the case of an empty string.
	if length == 0 {
		return "", nil
	}

	// 3. Read the string content.
	buffer := make([]byte, length)
	if _, err := io.ReadFull(r, buffer); err != nil {
		// If we can't read the full content, it's a malformed packet.
		return "", fmt.Errorf("could not read string content (expected %d bytes): %w", length, err)
	}

	// 4. Convert to string and return.
	return string(buffer), nil
}

func ReadByte(r io.Reader) (uint8, error) {
	opcodeBuffer := make([]byte, 1)
	n, err := r.Read(opcodeBuffer)
	if err != nil {
		return 0, fmt.Errorf("error reading opcode from stream: %w", err)
	}
	if n == 0 {
		return 0, fmt.Errorf("0 bytes read when expecting 1 byte")
	}
	return opcodeBuffer[0], nil
}
