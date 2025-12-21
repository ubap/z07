package sortcon

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortSource(t *testing.T) {

	t.Run("Basic Sorting by Hex Value", func(t *testing.T) {
		input := `package packets
const (
	S2CPing = 0x1E
	S2CLogin = 0x0A
	S2CMap = 0x64
)`
		expected := `package packets

const (
	S2CLogin = 0x0A
	S2CPing  = 0x1E
	S2CMap   = 0x64
)`
		output, err := SortSource([]byte(input))
		assert.NoError(t, err)
		// Trim spaces and newlines to avoid platform-specific failures
		assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(output)))
	})

	t.Run("Multiple Independent Blocks", func(t *testing.T) {
		input := `package packets
const (
	S2CB = 0x02
	S2CA = 0x01
)

const (
	C2SB = 0x02
	C2SA = 0x01
)`
		output, err := SortSource([]byte(input))
		assert.NoError(t, err)

		res := string(output)
		// Verify S2CA comes before S2CB
		assert.True(t, strings.Index(res, "S2CA") < strings.Index(res, "S2CB"), "S2CA should be before S2CB")
		// Verify C2SA comes before C2SB
		assert.True(t, strings.Index(res, "C2SA") < strings.Index(res, "C2SB"), "C2SA should be before C2SB")
	})
}
