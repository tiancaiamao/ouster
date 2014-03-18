package main

import (	
	"net"
	"log"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/scene"
	"github.com/tiancaiamao/ouster/login"
	"github.com/tiancaiamao/ouster/player"
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

	// get the map that player current in and connect to it
	m := scene.Query(playerData.Map)
	
	ch := make(chan interface{})
	playerId, succ := m.Login(ouster.Point(playerData.Pos), ch)

	// turn into a player agent
	if succ {
		agent := player.New(playerId, playerData, conn, ch)
		agent.Go()
	}	
}
