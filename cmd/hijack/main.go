package main

import (
	//	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	client, err := ln.Accept()
	if err != nil {
		panic(err)
	}

	server, err := net.Dial("tcp", "192.168.1.123:9999")
	if err != nil {
		panic(err)
	}

	hajack(client, server)
}

func hajack(client, server net.Conn) {
	clientBuf := make([]byte, 39)
	serverBuf := make([]byte, 200)

	flag := true
	for flag {
		// read from client
		n, err := client.Read(clientBuf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端关了，没法读取了")
			} else {
				panic(err)
			}
		}

		fmt.Println("[C->G]", clientBuf[:n])

		if n > 0 {
			// forward the msg to server
			_, err = server.Write(clientBuf[:n])
			if err != nil {
				if err == io.EOF {
					fmt.Println("服务端关了，没法把客户端数据写过去")
				} else {
					panic(err)
				}
			}
		}

		// read from server
		n, err = server.Read(serverBuf)
		if err != nil {
			if err == io.EOF {
				flag = false
				fmt.Println("服务器那边关了，没法读取了")
			} else {
				panic(err)
			}
		}

		fmt.Println("[G->C]", serverBuf[:n])

		// write the msg to client
		if n > 0 {
			_, err = client.Write(serverBuf[:n])
			if err != nil {
				if err == io.EOF {
					flag = false
					fmt.Println("客户端关了，写不到客户端了")
				} else {
					panic(err)
				}
			}
		}
	}
}
