package main

import (
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    "net"
)

func main() {
    ln, err := net.Listen("tcp", config.LoginServerPort)
    if err != nil {
        panic(err)
    }
    log.Infoln("loginserver started")
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Errorln("accept err:", err)
            continue
        }

        log.Infoln("receive a connect request")
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
            if _, ok := err.(packet.NotImplementError); !ok {
                log.Errorln("read packet error in loginserver's serve:", err)
                return
            }
        }

        log.Debugln("read a packet: ", pkt.PacketID())

        switch pkt.PacketID() {
        case packet.PACKET_CL_GET_WORLD_LIST:
            writer.Write(conn, packet.LCWorldListPacket{})
        case packet.PACKET_CL_LOGIN:
            writer.Write(conn, packet.LCLoginOKPacket{})
        case packet.PACKET_CL_SELECT_SERVER:
            writer.Write(conn, &packet.LCPCListPacket{})
        case packet.PACKET_CL_SELECT_WORLD:
            writer.Write(conn, &packet.LCServerListPacket{})
        case packet.PACKET_CL_VERSION_CHECK:
            writer.Write(conn, packet.LCVersionCheckOKPacket{})
        case packet.PACKET_CL_SELECT_PC:
            reconnect := &packet.LCReconnectPacket{
                Ip:   config.GameServerIP,
                Port: 9998,
                Key:  82180,
            }
            writer.Write(conn, reconnect)
            return
        default:
            log.Errorf("get a unknow packet: %d\n", pkt.PacketID())
        }
    }
}

// func Debug(writer packet.Writer, pkt packet.Packet) {
//		 stdout := &bytes.Buffer{}
//		 writer.Write(stdout, pkt)
//		 log.Println(stdout.Bytes())
// }
