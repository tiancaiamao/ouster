package main

import (
	// "github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster"
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

	aoi := make(chan uint32)
	scene2player := make(chan interface{})
	player2scene := make(chan interface{})

	agent := player.New(playerData, conn, aoi, scene2player, player2scene)

	// get the map that player current in
	m := scene.Query(playerData.Map)
	err = m.Login(agent, ouster.FPoint{
		X: float32(playerData.Pos.X),
		Y: float32(playerData.Pos.Y),
	}, aoi, player2scene, scene2player)
	if err != nil {
		// login to scene error
		return
	}

	// turn into a player agent
	agent.Go()
}
