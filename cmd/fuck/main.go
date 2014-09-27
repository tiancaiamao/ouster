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

var MAX int = 5

func main() {
    ln, err := net.Listen("tcp", ":9999")
    if err != nil {
        panic(err)
    }

    gn, err := net.Listen("tcp", ":9998")
    if err != nil {
        panic(err)
    }

    notice := make(chan struct{})
    go func() {
        client, err := ln.Accept()
        if err != nil {
            panic(err)
        }
        hajackLoginServer(client, notice)
    }()

    client, err := gn.Accept()
    if err != nil {
        panic(err)
    }

    hajackGameServer(client, notice)
}

func hajackLoginServer(client net.Conn, notice chan<- struct{}) {
    server, err := net.Dial("tcp", trueLoginServer)
    if err != nil {
        panic(err)
    }

    go func() {
        clientReader := packet.NewReader()
        serverWriter := packet.NewWriter()
        var raw packet.RawPacket
        var buf [800]byte
        for {
            err := clientReader.ReadRaw(client, &raw, buf[:])
            if err != nil {
                if err == io.EOF {
                    fmt.Println("CL客户端关了，没法读取了")
                    return
                } else {
                    panic(err)
                }
            }

            fmt.Println("[C->L]", len(raw.Data), raw)

            serverWriter.Seq = raw.Seq
            err = serverWriter.Write(server, raw)
            if err != nil {
                if err == io.EOF {
                    fmt.Println("LC服务端关了，没法把客户端数据写过去")
                    return
                } else {
                    panic(err)
                }
            }

        }
    }()

    serverReader := packet.NewReader()
    clientWriter := packet.NewWriter()
    var raw packet.RawPacket
    var buf [800]byte
    for {
        err := serverReader.ReadRaw(server, &raw, buf[:])
        if err != nil {
            if err == io.EOF {
                fmt.Println("LC服务器那边关了，没法读取了")
                return
            } else {
                panic(err)
            }
        }
        clientWriter.Seq = raw.Seq

        fmt.Println("[L->C]", len(raw.Data), raw)

        // 劫持修改最后的服务器包
        if raw.Id == packet.PACKET_LC_RECONNECT {
            pkt, err := raw.Unmarshal(0)
            if err != nil {
                panic(err)
            }
            reconn := pkt.(*packet.LCReconnectPacket)
            reconn.Ip = trueGameServer
            reconn.Port = 9998

            clientWriter.Write(client, raw)
            close(notice)
            client.Close()
            return
        }

        err = clientWriter.Write(client, raw)
        if err != nil {
            if err == io.EOF {
                fmt.Println("CL客户端关了，写不到客户端了")
            } else {
                panic(err)
            }
        }
    }
}

func hajackGameServer(client net.Conn, notice <-chan struct{}) {
    <-notice
    fmt.Println("收到客户端向GameServer连接请求...")
    server, err := net.Dial("tcp", trueGameServer+":9998")
    if err != nil {
        panic(err)
    }

    go func() {
        clientReader := packet.NewReader()
        serverWriter := packet.NewWriter()
        var raw packet.RawPacket
        var buf [800]byte
        for {
            err := clientReader.ReadRaw(client, &raw, buf[:])
            if err != nil {
                if err == io.EOF {
                    fmt.Println("CG客户端关了，没法读取了")
                } else {
                    panic(err)
                }
            }

            fmt.Println("[C->G]", raw)

            serverWriter.Seq = raw.Seq
            err = serverWriter.Write(server, &raw)
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
    var raw packet.RawPacket
    var buf [800]byte
    var correct uint8
    for {
        err := serverReader.ReadRaw(server, &raw, buf[:])
        if err != nil {
            if err == io.EOF {
                fmt.Println("CG服务器那边关了，没法读取了")
            } else {
                panic(err)
            }
        }

        fmt.Println("[G->C]", &raw)

        clientWriter.Code = raw.Seq + correct
        err = clientWriter.Write(client, raw)
        checkError(err)

        if raw.Id == packet.PACKET_CG_SKILL_TO_OBJECT {
            pkt, err := raw.Unmarshal(0)
            skill := pkt.(packet.CGSkillToObjectPacket)

            // 冰矛多倍
            if skill.SkillType == 285 {
                for i := 0; i < MAX; i++ {
                    correct++
                    clientWriter.Seq = raw.Seq + correct
                    err = clientWriter.Write(client, raw)
                    checkError(err)
                }
            }
        }
    }
}

func checkError(err error) {
    if err != nil {
        if err == io.EOF {
            fmt.Println("CG客户端关了，写不到客户端了")
        } else {
            panic(err)
        }
    }
}
