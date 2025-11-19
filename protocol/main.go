package main

import (
	"fmt"
	"net"
)

func main() {
	const rsaModulus = "138358917549655551601135922545920258651079249320630202917602000570926337770168654400102862016157293631277888588897291561865439132767832236947553872456033140205555218536070792283327632773558457562430692973109061064849319454982125688743198270276394129121891795353179249782548271479625552587457164097090236827371"
	const rsaExponent = "65537" // This is almost always the correct value.

	publicKey, err := ParseTibiaRSAPublicKey(rsaModulus, rsaExponent)
	if err != nil {
		panic(fmt.Sprintf("Failed to create RSA public key: %v", err))
	}

	// Establish a TCP connection to the specified host and port.
	conn, err := net.Dial("tcp", "world.fibula.app:7171")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	// Defer the closing of the connection to ensure it's closed when the function exits.
	defer conn.Close()

	// If the connection is successful, the error will be nil.
	fmt.Println("Connected successfully to world.fibula.app:7171")
}
