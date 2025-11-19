package main

import (
	"flag"
	"log"
	"os"
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
	inputFile := flag.String("binary", "Tibia.exe", "The binary to patch (default: Tibia.exe)")
	ip := flag.String("ip", "127.0.0.1", "New login servers IP address (default: 127.0.0.1)")

	flag.Parse()

	if len(*ip) > 20 {
		log.Fatalf("Too long IP address provided ('%s'), max length is 20 characters.", *ip)
	}

	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file '%s': %v", *inputFile, err)
	}

	content = changeLoginServers(content, *ip)
	content = changeRSAKey(content)

	err = os.WriteFile("Tibia_patched.exe", content, 0755)
	if err != nil {
		log.Fatalf("Error writing output file '%s': %v", inputFile, err)
	}

	log.Println("Successfully wrote output to Tibia_patched.exe")
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
