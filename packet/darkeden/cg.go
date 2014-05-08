package darkeden

import (
	"encoding/binary"
)

type CGConnectPacket struct {
	Key uint32
	// Slayer or Vampire?
	PCType     uint8
	PCName     string
	MacAddress [4]byte
}

func (connect *CGConnectPacket) Id() PacketID {
	return PACKET_CG_CONNECT
}
func (connect *CGConnectPacket) String() string {
	return "connect"
}
func readConnect(buf []byte) (Packet, error) {
	// [ 0 0 0 0 240 1 4 183 232 191 241 0 80 86 192 0 8]
	ret := new(CGConnectPacket)
	ret.Key = binary.LittleEndian.Uint32(buf[:4])
	ret.PCType = buf[5]
	length := buf[6]
	ret.PCName = string(buf[7 : 7+length])
	copy(ret.MacAddress[:], buf[8+length:12+length])
	return ret, nil
}

type CGReadyPacket struct{}

func (ready CGReadyPacket) Id() PacketID {
	return PACKET_CG_READY
}
func (ready CGReadyPacket) String() string {
	return "ready"
}

type CGMovePacket struct {
	X   uint8
	Y   uint8
	Dir uint8
}

func (move CGMovePacket) Id() PacketID {
	return PACKET_CG_MOVE
}
func (move CGMovePacket) String() string {
	return "move"
}
func readMove(buf []byte) (Packet, error) {
	return nil, nil
}

type CGVerifyTimePacket struct{}

func (verifyTime CGVerifyTimePacket) Id() PacketID {
	return PACKET_CG_VERIFY_TIME
}
func (verifyTime CGVerifyTimePacket) String() string {
	return "verify time"
}
