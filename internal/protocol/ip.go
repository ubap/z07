package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

// IPToString converts a Little Endian uint32 to "127.0.0.1"
func IPToString(ipInt uint32) string {
	// Create a 4-byte slice
	ipBytes := make([]byte, 4)

	// Write the uint32 into bytes using LittleEndian
	binary.LittleEndian.PutUint32(ipBytes, ipInt)

	// net.IP works with the byte slice [127, 0, 0, 1]
	return net.IP(ipBytes).String()
}

// StringToIP converts "127.0.0.1" to Little Endian uint32
func StringToIP(ipString string) (uint32, error) {
	// Parse the string
	ip := net.ParseIP(ipString)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipString)
	}

	// Ensure it is a 4-byte representation (IPv4)
	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipString)
	}

	// Convert bytes to uint32 using LittleEndian
	return binary.LittleEndian.Uint32(ip), nil
}
