/*
*

	go-tibia - A Man-in-the-Middle (MITM) proxy for the Tibia MMO
	Copyright (C) 2025 Jakub Trzebiatowski

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"goTibia/handlers/game"
	"goTibia/handlers/login"
	"goTibia/proxy"
	"goTibia/resources"
	"log"
	"sync"
)

func main() {
	if err := resources.LoadItemsJson("data/772/items.json"); err != nil {
		log.Fatalf("Critical Error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	loginHandler := &login.LoginHandler{
		TargetAddr: "world.fibula.app:7171",
		ProxyMOTD:  "Welcome to GoTibia Proxy!",
	}

	gameHandler := &game.GameHandler{
		TargetAddr: "world.fibula.app:7172",
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
