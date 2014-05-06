package main

import (
	//	"bytes"
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	gn, err := net.Listen("tcp", ":9998")
	if err != nil {
		panic(err)
	}

	go func() {
		client, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		server, err := net.Dial("tcp", "192.168.1.123:9999")
		if err != nil {
			panic(err)
		}

		hajackLoginServer(client, server)
	}()

	client, err := gn.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Println("收到客户端向GameServer连接请求...")
	server, err := net.Dial("tcp", "192.168.1.123:9998")
	if err != nil {
		panic(err)
	}

	hajackGameServer(client, server)
}

func hajackLoginServer(client, server net.Conn) {
	clientBuf := make([]byte, 200)
	serverBuf := make([]byte, 200)

	flag := true
	for flag {
		// read from client
		n, err := client.Read(clientBuf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("CL客户端关了，没法读取了")
			} else {
				panic(err)
			}
		}

		fmt.Println("[C->L]", clientBuf[:n])

		if n > 0 {
			// forward the msg to server
			_, err = server.Write(clientBuf[:n])
			if err != nil {
				if err == io.EOF {
					fmt.Println("CL服务端关了，没法把客户端数据写过去")
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
				fmt.Println("CL服务器那边关了，没法读取了")
			} else {
				panic(err)
			}
		}

		fmt.Println("[L->C]", serverBuf[:n])

		if bytes.HasPrefix(serverBuf[:n], []byte{194, 1}) {
			//  [194 1 22 0 0 0 5 13 49 57 50 46 49 54 56 46 49 46 49 50 51 14 39 0 0 0 32 6 11]
			key := serverBuf[n-6 : n]
			client.Write([]byte{194, 1, 20, 0, 0, 0, 5, 11, 49, 57, 50, 46, 49, 54, 56, 46, 49, 46, 50, 14, 39, key[0], key[1], key[2], key[3], key[4], key[5]})
			client.Close()
			return
		}

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

func hajackGameServer(client, server net.Conn) {
	clientBuf := make([]byte, 200)
	serverBuf := make([]byte, 200)

	go func() {
		for {
			// read from client
			n, err := client.Read(clientBuf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("CG客户端关了，没法读取了")
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
						fmt.Println("CG服务端关了，没法把客户端数据写过去")
					} else {
						panic(err)
					}
				}
			}
		}

	}()

	flag := true
	for flag {
		// read from server
		n, err := server.Read(serverBuf)
		if err != nil {
			if err == io.EOF {
				flag = false
				fmt.Println("CG服务器那边关了，没法读取了")
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
					fmt.Println("CG客户端关了，写不到客户端了")
				} else {
					panic(err)
				}
			}
		}
	}
}
