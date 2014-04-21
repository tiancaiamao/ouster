package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/packet"
	"net"
)

func main() {
	conn, err := Login()
	if err != nil {
		panic(err)
	}

	go func(c net.Conn) {
		for {
			pkt, err := packet.Read(c)
			if err != nil {
				fmt.Println(err)
				c.Close()
			}

			fmt.Println(pkt)
		}
	}(conn)
}

func Login() (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1"+config.ServerPort)
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
