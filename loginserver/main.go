package main

import (
    "bytes"
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/packet"
    "github.com/tiancaiamao/ouster/packet/darkeden"
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

    reader := darkeden.NewReader()
    writer := darkeden.NewWriter()

    for {
        pkt, err := reader.Read(conn)
        if err != nil {
            log.Println("read packet error in loginserver's serve:", err)
            return
        }

        // log.Println("read a packet: ", pkt.Id())

        switch pkt.Id() {
        case darkeden.PACKET_CL_GET_WORLD_LIST:
            writer.Write(conn, darkeden.LCWorldListPacket{})
            Debug(writer, darkeden.LCWorldListPacket{})
        case darkeden.PACKET_CL_LOGIN:
            writer.Write(conn, darkeden.LCLoginOKPacket{})
            Debug(writer, darkeden.LCLoginOKPacket{})
        case darkeden.PACKET_CL_SELECT_SERVER:
            writer.Write(conn, &darkeden.LCPCListPacket{})
            Debug(writer, &darkeden.LCPCListPacket{})
        case darkeden.PACKET_CL_SELECT_WORLD:
            writer.Write(conn, &darkeden.LCServerListPacket{})
            Debug(writer, &darkeden.LCServerListPacket{})
        case darkeden.PACKET_CL_VERSION_CHECK:
            writer.Write(conn, darkeden.LCVersionCheckOKPacket{})
            Debug(writer, darkeden.LCVersionCheckOKPacket{})
        case darkeden.PACKET_CL_SELECT_PC:
            reconnect := &darkeden.LCReconnectPacket{
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
