package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	// rsaPattern looks for a sequence of digits at least 245 chars long.
	// Standard Tibia RSA keys (1024-bit) are typically ~309 digits in decimal.
	rsaPattern = `\d{245,}`
)

func main() {
	// 1. Define Flags
	inputFile := flag.String("binary", "Tibia.exe", "Path to the Tibia binary to scan")
	outputFile := flag.String("output", "rsa_key.txt", "File path to save the extracted key")
	flag.Parse()

	// 2. Read Binary
	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file '%s': %v", *inputFile, err)
	}

	// 3. Compile Regex
	re, err := regexp.Compile(rsaPattern)
	if err != nil {
		log.Fatalf("Invalid regex pattern: %v", err)
	}

	// 4. Find Match
	// FindIndex returns [start, end]
	matchIndexes := re.FindIndex(content)

	if matchIndexes == nil {
		log.Fatal("No RSA key found in the binary matching the pattern.")
	}

	rsaOffset := matchIndexes[0]
	rsaKey := content[matchIndexes[0]:matchIndexes[1]]

	// 5. Output Results
	fmt.Println("---------------------------------------------------")
	fmt.Printf("Found RSA Key!\n")
	fmt.Printf("File Offset (Decimal): %d\n", rsaOffset)
	fmt.Printf("File Offset (Hex):     0x%X\n", rsaOffset)
	fmt.Printf("Key Length:            %d digits\n", len(rsaKey))
	fmt.Println("---------------------------------------------------")

	// 6. Save to File
	if err := os.WriteFile(*outputFile, rsaKey, 0644); err != nil {
		log.Fatalf("Error writing to output file '%s': %v", *outputFile, err)
	}

	fmt.Printf("Key successfully saved to: %s\n", *outputFile)
}
