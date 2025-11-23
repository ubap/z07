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

	go func() {
		defer wg.Done()
		srv := proxy.NewServer(
			"Login",
			":7171",
			"world.fibula.app:7171",
			login.HandleLoginLogic,
		)
		log.Fatal(srv.Start())
	}()

	go func() {
		defer wg.Done()
		srv := proxy.NewServer(
			"Game",
			":7172",
			"world.fibula.app:7172",
			game.HandleGameConnection,
		)
		log.Fatal(srv.Start())
	}()

	wg.Wait()
}
