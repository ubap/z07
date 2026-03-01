package main

import (
	"fmt"
	"z07/internal/client"
)

func main() {
	err := client.Login("127.0.0.1:7171")
	if err != nil {
		fmt.Println(err)
	}
}
