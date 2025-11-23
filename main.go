package main

import (
	"goTibia/handlers/game"
	"goTibia/handlers/login"
	"goTibia/proxy"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	loginHandler := &login.LoginHandler{
		TargetAddr: "world.fibula.app:7171",
		ProxyMOTD:  "Welcome to GoTibia Proxy!",
	}

	gameHandler := &game.GameHandler{
		TargetAddr: "world.fibula.app:7171",
	}

	go func() {
		defer wg.Done()
		srv := proxy.NewServer(
			"Login",
			":7171",
			loginHandler,
		)
		log.Fatal(srv.Start())
	}()

	go func() {
		defer wg.Done()
		srv := proxy.NewServer(
			"Game",
			":7172",
			gameHandler,
		)
		log.Fatal(srv.Start())
	}()

	wg.Wait()
}
