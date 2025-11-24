package main

import (
	"fmt"
	"log"
	"os" // Using the 'os' package, as 'ioutil' is deprecated
	"regexp"
)

// define all options here
const (
	// where to save RSA key
	outputFile = "RSA.txt"

	// where is your Tibia client
	inputFile = "/Users/jtrzebiatowski/programowanie_go/goTibia/resources/772/Tibia.exe"

	// this is most important, it assumes that RSA key is a string of 245 digits long or longer
	// and no other string is so long like this one.
	rsaPattern = `\d{245,}`
)

func main() {
	// Read the entire content of the input file.
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file '%s': %v", inputFile, err)
	}

	// Compile the regular expression.
	re, err := regexp.Compile(rsaPattern)
	if err != nil {
		log.Fatalf("Invalid regular expression pattern: %v", err)
	}

	// Find the indexes of the first match in the content.
	// FindIndex returns a slice with [startIndex, endIndex].
	matchIndexes := re.FindIndex(content)

	// Check if a match was found.
	if matchIndexes == nil {
		fmt.Println("No RSA key found matching the specified pattern.")
		return
	}

	// The starting index is the address (offset) of the key.
	rsaAddress := matchIndexes[0]

	// Extract the key itself by slicing the content using the found indexes.
	rsaKey := content[matchIndexes[0]:matchIndexes[1]]

	// Write the found RSA key to the output file.
	err = os.WriteFile(outputFile, rsaKey, 0644)
	if err != nil {
		log.Fatalf("Error writing to output file '%s': %v", outputFile, err)
	}

	// --- MODIFICATION ---
	// Confirm by outputting the address and the result to the console.
	// We print the address in both decimal and hexadecimal format.
	fmt.Printf("Found RSA key at address (offset): %d (0x%X)\n", rsaAddress, rsaAddress)
	fmt.Printf("Successfully extracted RSA key: %s\n", rsaKey)
	fmt.Printf("RSA key saved to %s\n", outputFile)
}
