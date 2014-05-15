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

	agent := scene.NewPlayer(conn)

	// get the map that player current in
	m := scene.Query("limbo_lair_se")
	if m == nil {
		panic("what the fuck??")
	}
	err := m.Login(agent)

	if err != nil {
		// login to scene error
		return
	}

	// turn into a player agent
	agent.Go()
}
