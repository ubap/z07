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

func WriteString(w io.Writer, str string) error {
	// 1. Write the 2-byte length prefix.
	length := uint16(len(str))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("could not write string length: %w", err)
	}

	// 2. Write the string content.
	if length > 0 {
		if n, err := w.Write([]byte(str)); err != nil {
			return fmt.Errorf("could not write string content: %w", err)
		} else if n != int(length) {
			return fmt.Errorf("incomplete write: expected %d bytes, wrote %d bytes", length, n)
		}
	}
	return nil
}

func ReadByte(r io.Reader) (uint8, error) {
	opcodeBuffer := make([]byte, 1)
	n, err := r.Read(opcodeBuffer)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmt.Errorf("0 bytes read when expecting 1 byte")
	}
	return opcodeBuffer[0], nil
}
