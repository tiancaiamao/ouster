package darkeden

import (
	"github.com/tiancaiamao/ouster/packet"
)

type CLLoginPacket struct {
	Username string
	Password string
}

func (login *CLLoginPacket) Id() packet.PacketID {
	return PACKET_CL_LOGIN
}

func (login *CLLoginPacket) String() string {
	return "Login"
}

func readLogin(buf []byte, code uint8) (packet.Packet, error) {
	szUsername := int(buf[0])
	szPassword := int(buf[1+szUsername])
	return &CLLoginPacket{
		Username: string(buf[1 : 1+szUsername]),
		Password: string(buf[2+szUsername : 2+szUsername+szPassword]),
	}, nil
}

type CLVersionCheckPacket struct{}

func (v CLVersionCheckPacket) Id() packet.PacketID {
	return PACKET_CL_VERSION_CHECK
}
func (v CLVersionCheckPacket) String() string {
	return "version check"
}

type CLGetWorldListPacket struct{}

func (worldList CLGetWorldListPacket) Id() packet.PacketID {
	return PACKET_CL_GET_WORLD_LIST
}
func (w CLGetWorldListPacket) String() string {
	return "get world list"
}
func readGetWorldList(buf []byte, code uint8) (packet.Packet, error) {
	return CLGetWorldListPacket{}, nil
}

type CLSelectWorldPacket uint8

func (sw CLSelectWorldPacket) Id() packet.PacketID {
	return PACKET_CL_SELECT_WORLD
}
func (sw CLSelectWorldPacket) String() string {
	return "select world"
}

func readSelectWorld(buf []byte, code uint8) (packet.Packet, error) {
	return CLSelectWorldPacket(buf[0]), nil
}

type CLSelectServerPacket uint8

func (ss CLSelectServerPacket) Id() packet.PacketID {
	return PACKET_CL_SELECT_SERVER
}
func (ss CLSelectServerPacket) String() string {
	return "select server"
}
func readSelectServer(buf []byte, code uint8) (packet.Packet, error) {
	return CLSelectServerPacket(buf[0]), nil
}

type PCType uint8
type CLSelectPcPacket struct {
	Name string
	Type PCType
}

func (sp *CLSelectPcPacket) Id() packet.PacketID {
	return PACKET_CL_SELECT_PC
}
func (sp *CLSelectPcPacket) String() string {
	return "select pc"
}
func readSelectPc(buf []byte, code uint8) (packet.Packet, error) {
	//	[8 178 187 212 217 209 218 202 206 0]
	sz := buf[0]
	return &CLSelectPcPacket{
		Type: PCType(buf[len(buf)-1]),
		Name: string(buf[1 : 1+sz]),
	}, nil
}
