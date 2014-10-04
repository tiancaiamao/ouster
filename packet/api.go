package packet

import (
    "io"
)

type PacketID uint16

type Packet interface {
    PacketID() PacketID
    PacketSize() uint32
    Read(reader io.Reader, code uint8) error
    Write(writer io.Writer, code uint8) error
}

// type Writer interface {
//		 Write(writer io.Writer, pkt Packet) error
// }
//
// type Reader interface {
//		 Read(reader io.Reader) (ret Packet, err error)
// }
