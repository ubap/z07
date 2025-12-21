package packets

import (
	"goTibia/internal/tools/sortcon"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpcodesAreSortedOnDisk(t *testing.T) {
	content, err := os.ReadFile("opcodes.go")
	if err != nil {
		t.Fatalf("could not read opcodes.go: %v", err)
	}

	sortedContent, err := sortcon.SortSource(content)
	assert.NoError(t, err)

	assert.Equal(t, string(content), string(sortedContent),
		"opcodes.go is not sorted. Please run 'go generate ./...' to fix it.")
}
