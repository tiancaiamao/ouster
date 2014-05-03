package main

import (
	"github.com/tiancaiamao/ouster/packet/darkeden"
	"log"
	"net"
	// "os"
)

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go serve(conn)
	}
}

func serve(conn net.Conn) {
	for {
		pkt, err := darkeden.Read(conn)
		if err != nil {
			log.Println("read packet error in loginserver's serve:", err)
			conn.Close()
			return
		}

		log.Println("read a packet: ", pkt.Id())

		switch pkt.Id() {
		case darkeden.PACKET_CL_GET_WORLD_LIST:
			darkeden.Write(conn, darkeden.LCWorldListPacket{})
		case darkeden.PACKET_CL_LOGIN:
			darkeden.Write(conn, darkeden.LCLoginOKPacket{})
		case darkeden.PACKET_CL_SELECT_SERVER:
			darkeden.Write(conn, &darkeden.LCPCListPacket{})
		case darkeden.PACKET_CL_SELECT_WORLD:
			darkeden.Write(conn, &darkeden.LCServerListPacket{})
		case darkeden.PACKET_CL_VERSION_CHECK:
			darkeden.Write(conn, darkeden.LCVersionCheckOKPacket{})
		case darkeden.PACKET_CL_SELECT_PC:
			reconnect := &darkeden.LCReconnectPacket{
				Ip:   "192.168.1.123",
				Port: 9998,
				Key:  []byte{0, 0, 0, 32, 6, 11},
			}
			darkeden.Write(conn, reconnect)
		default:
			log.Printf("get a unknow packet: %d\n", pkt.Id())
		}
	}
}
