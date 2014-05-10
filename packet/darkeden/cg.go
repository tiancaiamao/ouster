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
	// [ 0 0 0 240 1 4 183 232 191 241 0 80 86 192 0 8]
	ret := new(CGConnectPacket)
	ret.Key = binary.LittleEndian.Uint32(buf[:4])
	ret.PCType = buf[4]
	length := buf[5]
	ret.PCName = string(buf[6 : 6+length])
	copy(ret.MacAddress[:], buf[6+length:])
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
	Dir uint8
	X   uint8
	Y   uint8
}

func (move CGMovePacket) Id() PacketID {
	return PACKET_CG_MOVE
}
func (move CGMovePacket) String() string {
	return "move"
}
func readMove(buf []byte) (Packet, error) {
	ret := CGMovePacket{
		X:   buf[0],
		Dir: buf[1],
		Y:   buf[2],
	}
	return ret, nil
}

type CGVerifyTimePacket struct{}

func (verifyTime CGVerifyTimePacket) Id() PacketID {
	return PACKET_CG_VERIFY_TIME
}
func (verifyTime CGVerifyTimePacket) String() string {
	return "verify time"
}

type CGAttackPacket struct {
	ObjectID uint32
	X        uint8
	Y        uint8
	Dir      uint8
}

func (attack CGAttackPacket) Id() PacketID {
	return PACKET_CG_ATTACK
}
func (attack CGAttackPacket) String() string {
	return "attack"
}
func readAttack([]byte) (Packet, error) {
	var ret CGAttackPacket
	return ret, nil
}
func (attack CGAttackPacket) Bytes() []byte {
	// [55 218 53 0 0 39 189]
	ret := make([]byte, 7)
	binary.LittleEndian.PutUint32(ret, attack.ObjectID)
	ret[4] = attack.X
	ret[5] = attack.Y
	ret[6] = attack.Dir
	return ret
}
