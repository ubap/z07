package assets

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadItemsJson reads the JSON file and populates the global 'Things' slice.
func LoadItemsJson(path string) error {
	// 1. Read the file from disk
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open items file: %w", err)
	}

	// 2. Parse JSON into a temporary slice
	var loadedItems []ItemType
	if err := json.Unmarshal(fileBytes, &loadedItems); err != nil {
		return fmt.Errorf("failed to parse items json: %w", err)
	}

	// 3. Find the Maximum ID to determine slice size
	// We need to know how big 'Things' should be to avoid index out of range.
	maxID := 0
	for _, item := range loadedItems {
		if int(item.ID) > maxID {
			maxID = int(item.ID)
		}
	}

	// 4. Initialize the Global Registry (Allocate memory)
	// We add +1 because IDs are 0-indexed in the slice.
	Initialize(maxID + 1)

	// 5. Populate the Global Registry
	count := 0
	for _, item := range loadedItems {
		// Verify bounds just in case
		if int(item.ID) < len(Things) {
			Things[item.ID] = item
			count++
		}
	}

	fmt.Printf("[Data] Loaded %d items from %s. Max ID: %d\n", count, path, maxID)
	return nil
}
