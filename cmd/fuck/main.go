package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster/packet"
	"io"
	"net"
)

const (
	trueLoginServer = "60.169.77.55:9999"
	trueGameServer  = "60.169.77.55"
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

		server, err := net.Dial("tcp", trueLoginServer)
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
	server, err := net.Dial("tcp", trueGameServer+":9998")
	if err != nil {
		panic(err)
	}

	hajackGameServer(client, server)
}

func hajackLoginServer(client, server net.Conn) {
	go func() {
		clientReader := packet.NewReader()
		serverWriter := packet.NewWriter()
		for {
			pkt, err := clientReader.Read(client)
			if err != nil {
				if _, ok := err.(packet.NotImplementError); !ok {
					if err == io.EOF {
						fmt.Println("CG客户端关了，没法读取了")
					} else {
						panic(err)
					}
				}
			}

			fmt.Printf("[C->L %d] %v\n", clientReader.Seq , pkt)

			// 劫持修改最后的服务器包
			if pkt.PacketID() == packet.PACKET_LC_RECONNECT {
				raw := pkt.(*packet.LCReconnectPacket)
				raw.Ip = trueGameServer
				raw.Port = 9998

			serverWriter.Write(server, raw)
				client.Close()
				return
			}

			err = serverWriter.Write(server, pkt)
			if err != nil {
				if err == io.EOF {
					fmt.Println("CG服务端关了，没法把客户端数据写过去")
				} else {
					panic(err)
				}
			}

		}
	}()

	serverReader := packet.NewReader()
	clientWriter := packet.NewWriter()
	for {
		pkt, err := serverReader.Read(server)
		if err != nil {
			if _, ok := err.(packet.NotImplementError); !ok {
				if err == io.EOF {
					fmt.Println("CG服务器那边关了，没法读取了")
				} else {
					panic(err)
				}
			}
		}

		fmt.Println("[L->C]", pkt)

		err = clientWriter.Write(client, pkt)
		if err != nil {
			if err == io.EOF {
				fmt.Println("CG客户端关了，写不到客户端了")
			} else {
				panic(err)
			}
		}
	}
}

func hajackGameServer(client, server net.Conn) {
	go func() {
		clientReader := packet.NewReader()
		serverWriter := packet.NewWriter()
		for {
			pkt, err := clientReader.Read(client)
			if err != nil {
				if _, ok := err.(packet.NotImplementError); !ok {
					if err == io.EOF {
						fmt.Println("CG客户端关了，没法读取了")
					} else {
						panic(err)
					}
				}
			}

			fmt.Println("[C->G]", pkt)

			err = serverWriter.Write(server, pkt)
			if err != nil {
				if err == io.EOF {
					fmt.Println("CG服务端关了，没法把客户端数据写过去")
				} else {
					panic(err)
				}
			}
		}
	}()

	serverReader := packet.NewReader()
	clientWriter := packet.NewWriter()
	for {
		pkt, err := serverReader.Read(server)
		if err != nil {
			if _, ok := err.(packet.NotImplementError); !ok {
				if err == io.EOF {
					fmt.Println("CG服务器那边关了，没法读取了")
				} else {
					panic(err)
				}
			}
		}

		fmt.Println("[G->C]", pkt)

		err = clientWriter.Write(client, pkt)
		if err != nil {
			if err == io.EOF {
				fmt.Println("CG客户端关了，写不到客户端了")
			} else {
				panic(err)
			}
		}
	}
}
