package packet

import (
    "io"
)

type PacketID uint16

type Packet interface {
    Id() PacketID
    String() string
}

type PacketWriter interface {
    Write(writer io.Writer, pkt Packet) error
}

type PacketReader interface {
    Read(reader io.Reader) (ret Packet, err error)
}
