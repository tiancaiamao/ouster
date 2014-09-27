package main

import (
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/tiancaiamao/ouster/packet"
    "io"
    "net"
)

const (
    trueLoginServer = "60.169.77.55:9999"
    trueGameServer  = "60.169.77.55:9998"
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
        var raw packet.RawPacket
        var buf [400]byte
        for {
            err := readRaw(client, &raw, buf[:])
            if err != nil {
                if err == io.EOF {
                    fmt.Println("CL客户端关了，没法读取了")
                    return
                } else {
                    panic(err)
                }
            }

            fmt.Println("[C->L]", len(raw.Data), raw)

            err = writeRaw(server, &raw)
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

    var raw packet.RawPacket
    var buf [400]byte
    for {
        err := readRaw(server, &raw, buf[:])
        if err != nil {
            if err == io.EOF {
                fmt.Println("LC服务器那边关了，没法读取了")
                return
            } else {
                panic(err)
            }
        }

        fmt.Println("[L->C]", len(raw.Data), raw)

        // if raw.Id == packet.PACKET_LC_PC_LIST {
        //     fmt.Println("run here...")
        //     pkt := &packet.LCPCListPacket{}

        //     raw.Data, _ = pkt.MarshalBinary(0)
        //     fmt.Println("[实际发送L->C]", len(raw.Data), raw)
        // }

        // 劫持修改最后的服务器包
        if raw.Id == packet.PACKET_LC_RECONNECT {
            pkt, err := raw.Unmarshal(0)
            if err != nil {
                panic(err)
            }
            reconn := pkt.(*packet.LCReconnectPacket)
            reconn.Ip = "192.168.10.102"
            reconn.Port = 9998

            raw.Data, err = reconn.MarshalBinary(0)
            err = writeRaw(client, &raw)
            close(notice)
            client.Close()
            return
        }

        err = writeRaw(client, &raw)
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
    server, err := net.Dial("tcp", trueGameServer)
    if err != nil {
        panic(err)
    }

    go func() {
        var raw packet.RawPacket
        var buf [800]byte
        for {
            err := readRaw(client, &raw, buf[:])
            if err != nil {
                if err == io.EOF {
                    fmt.Println("CG客户端关了，没法读取了")
                } else {
                    panic(err)
                }
            }

            fmt.Println("[C->G]", raw)

            err = writeRaw(server, &raw)
            if err != nil {
                if err == io.EOF {
                    fmt.Println("CG服务端关了，没法把客户端数据写过去")
                } else {
                    panic(err)
                }
            }
        }
    }()

    var raw packet.RawPacket
    var buf [800]byte
    //    var correct uint8
    for {
        err := readRaw(server, &raw, buf[:])
        if err != nil {
            if err == io.EOF {
                fmt.Println("CG服务器那边关了，没法读取了")
            } else {
                panic(err)
            }
        }

        fmt.Println("[G->C]", raw)

        err = writeRaw(client, &raw)
        if err != nil {
            panic(err)
        }

        // if raw.Id == packet.PACKET_CG_SKILL_TO_OBJECT {
        //     pkt, err := raw.Unmarshal(0)
        //     skill := pkt.(packet.CGSkillToObjectPacket)

        //     // 冰矛多倍
        //     if skill.SkillType == 285 {
        //         for i := 0; i < MAX; i++ {
        //             correct++
        //             raw.Seq += correct
        //             err = writeRaw.Write(client, raw)
        //             if err != nil {
        //                 panic(err)
        //             }
        //         }
        //     }
        // }
    }
}

func writeRaw(writer io.Writer, raw *packet.RawPacket) error {
    err := binary.Write(writer, binary.LittleEndian, raw.Id)
    if err != nil {
        return err
    }

    err = binary.Write(writer, binary.LittleEndian, uint32(len(raw.Data)))
    if err != nil {
        return err
    }

    err = binary.Write(writer, binary.LittleEndian, raw.Seq)
    if err != nil {
        return err
    }

    _, err = writer.Write(raw.Data)
    if err != nil {
        return err
    }

    return nil
}

func readRaw(reader io.Reader, raw *packet.RawPacket, buf []byte) (err error) {
    var sz uint32

    err = binary.Read(reader, binary.LittleEndian, &raw.Id)
    if err != nil {
        return
    }

    err = binary.Read(reader, binary.LittleEndian, &sz)
    if err != nil {
        return
    }

    err = binary.Read(reader, binary.LittleEndian, &raw.Seq)
    if err != nil {
        return
    }

    if sz > uint32(len(buf)) {
        return errors.New("packet size too large")
    }

    n, err := io.ReadFull(reader, buf[:sz])
    if err != nil {
        return
    }
    if n != int(sz) {
        err = errors.New("read get less data than needed")
        return
    }

    raw.Data = buf[:sz]
    return nil
}
