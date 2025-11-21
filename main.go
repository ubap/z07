package main

import (
	"goTibia/login"
	"log"
)

func main() {
	loginServer := login.NewServer(":7171", "world.fibula.app:7171")

	go func() {
		log.Println("Starting login server...")
		if err := loginServer.Start(); err != nil {
			log.Fatalf("Login server failed: %v", err)
		}
	}()

	// select {} is a common way to block forever.
	select {}
}
