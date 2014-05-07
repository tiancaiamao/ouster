package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/packet"
	"io"
	"net"
)

func main() {
	conn, err := Login()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		ln, err := net.Listen("tcp", ":8782")
		if err != nil {
			panic(err)
		}

		fd, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		ln.Close()
		defer fd.Close()

		for {
			// forward packet to server
			io.Copy(conn, fd)
		}
	}()

	// goroutine for display packet from server
	for {
		pkt, err := packet.Read(conn)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(pkt)
	}
}

func Login() (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1"+config.GameServerPort)
	if err != nil {
		return nil, err
	}

	packet.Write(conn, packet.PLogin, packet.LoginPacket{
		Username: "genius",
		Password: "0101001",
	})

	info, err := packet.Read(conn)
	if err != nil {
		return nil, err
	}
	if _, ok := info.(packet.CharactorInfoPacket); !ok {
		return nil, ouster.NewError("need a CharactorInfoPacket")
	}

	packet.Write(conn, packet.PSelectCharactor, packet.SelectCharactorPacket{
		Which: 0,
	})

	loginOk, err := packet.Read(conn)
	if _, ok := loginOk.(packet.LoginOkPacket); !ok {
		return nil, ouster.NewError("need a LoginOkPacket")
	}

	return conn, nil
}
