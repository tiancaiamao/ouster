package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/packet/msgpack"
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
			// forward msgpack to server
			io.Copy(conn, fd)
		}
	}()

	// goroutine for display msgpack from server
	for {
		pkt, err := msgpack.Read(conn)
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

	msgpack.Write(conn, msgpack.PLogin, msgpack.Loginmsgpack{
		Username: "genius",
		Password: "0101001",
	})

	info, err := msgpack.Read(conn)
	if err != nil {
		return nil, err
	}
	if _, ok := info.(msgpack.CharactorInfomsgpack); !ok {
		return nil, ouster.NewError("need a CharactorInfomsgpack")
	}

	msgpack.Write(conn, msgpack.PSelectCharactor, msgpack.SelectCharactormsgpack{
		Which: 0,
	})

	loginOk, err := msgpack.Read(conn)
	if _, ok := loginOk.(msgpack.LoginOkmsgpack); !ok {
		return nil, ouster.NewError("need a LoginOkmsgpack")
	}

	return conn, nil
}
