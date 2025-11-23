package main

import (
	"goTibia/game_server"
	"goTibia/login_server"
	"log"
)

func main() {
	loginServer := login_server.NewServer(":7171", "world.fibula.app:7171")
	gameServer := game_server.NewServer(":7172", "world.fibula.app:7172")

	go func() {
		log.Println("Starting login server...")
		if err := loginServer.Start(); err != nil {
			log.Fatalf("Login server failed: %v", err)
		}
	}()

	go func() {
		log.Println("Starting game server...")
		if err := gameServer.Start(); err != nil {
			log.Fatalf("Game server failed: %v", err)
		}
	}()

	// select {} is a common way to block forever.
	select {}
}
