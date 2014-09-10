package main

import (
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/log"
    "net"
)

func main() {
    log.Infoln("Starting the server.")

    Initialize()

    listener, err := net.Listen("tcp", config.GameServerPort)
    checkError(err)

    log.Infoln("Game Server OK.")

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
    log.Debug("accept a connection...")

    agent := NewAgent(conn)
    go agent.Loop()
}
