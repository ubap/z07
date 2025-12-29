package bot

import (
	"fmt"
	"strings"
)

// FormatForTest converts a byte slice into a string you can copy-paste into a Go test.
func FormatForTest(name string, data []byte) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("// %s\n", name))
	sb.WriteString("rawData := []byte{")
	for i, b := range data {
		sb.WriteString(fmt.Sprintf("0x%02X", b))
		if i < len(data)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
