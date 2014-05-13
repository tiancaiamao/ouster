package main

import (
	// "github.com/tiancaiamao/ouster"
	// "github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	// "github.com/tiancaiamao/ouster/player"
	"github.com/tiancaiamao/ouster/scene"
	"log"
	"net"
)

func main() {
	log.Println("Starting the server.")

	scene.Initialize()

	listener, err := net.Listen("tcp", config.GameServerPort)
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
	defer conn.Close()

	aoi := make(chan scene.ObjectIDType)
	scene2player := make(chan interface{})
	player2scene := make(chan interface{})

	agent := scene.NewPlayer(conn, aoi, scene2player, player2scene)

	// get the map that player current in
	m := scene.Query("limbo_lair_se")
	if m == nil {
		panic("what the fuck??")
	}
	err := m.Login(agent, 40, 50, aoi, player2scene, scene2player)

	if err != nil {
		// login to scene error
		return
	}

	// turn into a player agent
	agent.Go()
}
