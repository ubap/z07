package main

import (
	"fmt"
	"goTibia/internal/tools/sortcon"
	"os"
)

func main() {
	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		if len(os.Args) < 2 {
			fmt.Println("Usage: Run via 'go generate' or provide a file path.")
			os.Exit(1)
		}
		fileName = os.Args[1]
	}

	src, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	out, err := sortcon.SortSource(src) // Call the library function
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	os.WriteFile(fileName, out, 0644)
}
