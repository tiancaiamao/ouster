package packet

import (
    "encoding/binary"
    "errors"
    "github.com/tiancaiamao/ouster/log"
    "io"
)

type PacketID uint16

type Packet interface {
    PacketID() PacketID
    PacketSize() uint32
    Read(reader io.Reader, code uint8) error
    Write(writer io.Writer, code uint8) error
}

func Write(writer io.Writer, pkt Packet, code uint8, seq uint8) error {
    id := pkt.PacketID()
    err := binary.Write(writer, binary.LittleEndian, id)
    if err != nil {
        return err
    }

    sz := pkt.PacketSize()
    err = binary.Write(writer, binary.LittleEndian, sz)
    if err != nil {
        return err
    }

    err = binary.Write(writer, binary.LittleEndian, seq)
    if err != nil {
        return err
    }

    err = pkt.Write(writer, code)
    if err != nil {
        return err
    }

    return nil
}

func Read(reader io.Reader, code uint8) (ret Packet, seq uint8, err error) {
    var id PacketID
    var sz uint32

    err = ReadHeader(reader, &id, &sz, &seq)
    if err != nil {
        return
    }

    if id >= PACKET_MAX {
        err = errors.New("packet id too large!")
        return
    }

    ret = packetTable[id]
    if ret == nil {
        log.Debugln("reading a not implement packet:", id)
        var buf [500]byte
        raw := RawPacket{
            Id:  id,
            Seq: seq,
        }
        if sz > uint32(len(buf)) {
            err = errors.New("too large raw packet")
            return
        }
        _, err = reader.Read(buf[:sz])
        if err != nil {
            return
        }
        copy(raw.Data, buf[:sz])
        err = NotImplementError{}
        return
    }

    err = ret.Read(reader, code)
    return
}
