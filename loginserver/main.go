package main

import (
    "bytes"
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/packet"
    "log"
    "net"
)

func main() {
    ln, err := net.Listen("tcp", config.LoginServerPort)
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
    defer conn.Close()

    reader := packet.NewReader()
    writer := packet.NewWriter()

    for {
        pkt, err := reader.Read(conn)
        if err != nil {
            log.Println("read packet error in loginserver's serve:", err)
            return
        }

        // log.Println("read a packet: ", pkt.Id())

        switch pkt.Id() {
        case packet.PACKET_CL_GET_WORLD_LIST:
            writer.Write(conn, packet.LCWorldListPacket{})
            Debug(writer, packet.LCWorldListPacket{})
        case packet.PACKET_CL_LOGIN:
            writer.Write(conn, packet.LCLoginOKPacket{})
            Debug(writer, packet.LCLoginOKPacket{})
        case packet.PACKET_CL_SELECT_SERVER:
            writer.Write(conn, &packet.LCPCListPacket{})
            Debug(writer, &packet.LCPCListPacket{})
        case packet.PACKET_CL_SELECT_WORLD:
            writer.Write(conn, &packet.LCServerListPacket{})
            Debug(writer, &packet.LCServerListPacket{})
        case packet.PACKET_CL_VERSION_CHECK:
            writer.Write(conn, packet.LCVersionCheckOKPacket{})
            Debug(writer, packet.LCVersionCheckOKPacket{})
        case packet.PACKET_CL_SELECT_PC:
            reconnect := &packet.LCReconnectPacket{
                Ip:   config.GameServerIP,
                Port: 9998,
                Key:  []byte{0, 0, 0, 32, 6, 11},
            }
            writer.Write(conn, reconnect)
            Debug(writer, reconnect)
            return
        default:
            log.Printf("get a unknow packet: %d\n", pkt.Id())
        }
    }
}

func Debug(writer packet.PacketWriter, pkt packet.Packet) {
    stdout := &bytes.Buffer{}
    writer.Write(stdout, pkt)
    log.Println(stdout.Bytes())
}
