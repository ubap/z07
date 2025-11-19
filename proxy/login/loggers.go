package main

import (
	"fmt"
	"strings"
)

func PrintAsGoSlice(packetBody []byte) {
	fmt.Println("\n--- COPY THE VARIABLE BELOW INTO YOUR TEST FILE ---")
	fmt.Println(formatAsGoSlice(packetBody))
	fmt.Println("--- END OF GO SLICE ---")
}

func formatAsGoSlice(data []byte) string {
	var builder strings.Builder
	builder.WriteString("capturedPacket := []byte{")

	for i, b := range data {
		// Add a newline every 12 bytes to keep it readable.
		if i%12 == 0 {
			builder.WriteString("\n\t")
		}
		// Write the byte in 0xXX format.
		builder.WriteString(fmt.Sprintf("0x%02x, ", b))
	}

	builder.WriteString("\n}")
	return builder.String()
}
