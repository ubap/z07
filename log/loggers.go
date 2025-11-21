package log

import (
	"encoding/hex"
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

type HexDumpWriter struct {
	// Prefix allows us to label the output, e.g., "SERVER->" or "CLIENT->"
	Prefix string
}

// Write is the only method needed to satisfy the io.Writer interface.
func (w *HexDumpWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("\n--- Data Dump (%s) ---\n", w.Prefix)
	fmt.Printf("%s", hex.Dump(p))
	fmt.Println("--- End of Dump ---")
	// We return the number of bytes processed and no error.
	return len(p), nil
}
