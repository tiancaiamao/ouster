package main

import (
	// "github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/login"
	"github.com/tiancaiamao/ouster/player"
	"github.com/tiancaiamao/ouster/scene"
	"log"
	"net"
)

func main() {
	log.Println("Starting the server.")

	scene.Initialize()

	listener, err := net.Listen("tcp", config.ServerPort)
	checkError(err)

	log.Println("Game Server OK.")

	for {
		conn, err := listener.Accept()
		if err == nil {
			go handleClient(conn)
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func handleClient(conn net.Conn) {
	log.Println("accept a connection...")
	// check username/password, load player info, and so on...
	playerData, err := login.Login(conn)
	if err != nil {
		return
	}

	agent := player.New(playerData, conn)

	// get the map that player current in
	m := scene.Query(playerData.Map)
	err = m.Login(agent)
	if err != nil {
		// login to scene error
		return
	}

	// turn into a player agent
	agent.Go()
}
