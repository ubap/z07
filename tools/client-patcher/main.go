package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Addresses in file (Not in memory)
const (
	tibia772LoginServerAddress = 0x0016D338
	loginServerEntrySize       = 20
	loginServersCount          = 4
	RSAAddress                 = 0x0015B620
	OTPublicRSA                = "109120132967399429278860960508995541528237502902798129123468757937266291492576446330739696001110603907230888610072655818825358503429057592827629436413108566029093628212635953836686562675849720620786279431090218017681061521755056710823876476444260558147179707119674283982419152118103759076030616683978566631413"
)

func main() {
	inputFile := flag.String("binary", "Tibia.exe", "Path to the Tibia binary")
	ip := flag.String("ip", "127.0.0.1", "Proxy IP address")

	flag.Parse()

	if len(*ip) > 20 {
		log.Fatalf("Too long IP address provided ('%s'), max length is 20 characters.", *ip)
	}

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading '%s': %v", *inputFile, err)
	}

	// 2. Patch
	content = changeLoginServers(content, *ip)
	content = changeRSAKey(content)

	// 3. Determine Output Path
	// If input is "C:/Games/Tibia/Tibia.exe", output is "C:/Games/Tibia/Tibia_patched.exe"
	// This ensures the patched exe can find Tibia.spr/dat in that same folder.
	dir := filepath.Dir(*inputFile)
	filename := filepath.Base(*inputFile)
	ext := filepath.Ext(filename)
	rawName := strings.TrimSuffix(filename, ext)

	outputName := fmt.Sprintf("%s_patched%s", rawName, ext)
	outputPath := filepath.Join(dir, outputName)

	// 4. Write Output
	// 0755 makes it executable on Linux/Mac (if running via Wine)
	err = os.WriteFile(outputPath, content, 0755)
	if err != nil {
		log.Fatalf("Error writing '%s': %v", outputPath, err)
	}

	log.Printf("Success! Created patched client at: %s", outputPath)
}

func changeLoginServers(content []byte, newIp string) []byte {
	newServerAddress := []byte(newIp)

	for loginSrvIndex := 0; loginSrvIndex < loginServersCount; loginSrvIndex++ {
		loginSrvAddress := tibia772LoginServerAddress + (loginSrvIndex * loginServerEntrySize)
		destinationSlice := content[loginSrvAddress : loginSrvAddress+loginServerEntrySize]
		copy(destinationSlice, newServerAddress)
		// Pad the rest of the area with null bytes (zeros).
		for i := len(newServerAddress); i < len(destinationSlice); i++ {
			destinationSlice[i] = 0
		}
	}
	return content
}

func changeRSAKey(content []byte) []byte {
	newRSAKey := []byte(OTPublicRSA)
	destinationSlice := content[RSAAddress : RSAAddress+len(newRSAKey)]
	copy(destinationSlice, newRSAKey)
	return content
}
