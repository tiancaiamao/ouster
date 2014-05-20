package darkeden

import (
	"encoding/binary"
	"github.com/tiancaiamao/ouster/packet"
)

type LCLoginOKPacket struct{}

func (loginOk LCLoginOKPacket) Id() packet.PacketID {
	return PACKET_LC_LOGIN_OK
}

func (loginOk LCLoginOKPacket) String() string {
	return "loginOK"
}

func (loginOk LCLoginOKPacket) MarshalBinary() ([]byte, error) {
	return []byte{0, 0, 231, 254, 255}, nil
}

type LCVersionCheckOKPacket struct{}

func (v LCVersionCheckOKPacket) Id() packet.PacketID {
	return PACKET_LC_VERSION_CHECK_OK
}

func (v LCVersionCheckOKPacket) String() string {
	return "version check ok"
}
func (v LCVersionCheckOKPacket) MarshalBinary() ([]byte, error) {
	return []byte{}, nil
}

type LCWorldListPacket struct{}

func (wl LCWorldListPacket) Id() packet.PacketID {
	return PACKET_LC_WORLD_LIST
}
func (wl LCWorldListPacket) String() string {
	return "world list"
}
func (wl LCWorldListPacket) MarshalBinary() ([]byte, error) {
	return []byte{1, 1, 1, 8, 185, 237, 247, 200, 193, 182, 211, 252, 0}, nil
}

type LCServerListPacket struct {
	CurrentWorld uint8
	Size         uint8
	list         []string
}

func (sl *LCServerListPacket) Id() packet.PacketID {
	return PACKET_LC_SERVER_LIST
}
func (sl *LCServerListPacket) String() string {
	return "server list"
}
func (sl *LCServerListPacket) MarshalBinary() ([]byte, error) {
	return []byte{1, 2, 0, 6, 183, 226, 178, 226, 199, 248, 0, 1, 8, 185, 237, 247, 200, 193, 182, 211, 252, 0}, nil
}

type LCPCListPacket struct {
	list []PCInfo
}

func (pl *LCPCListPacket) Id() packet.PacketID {
	return PACKET_LC_PC_LIST
}
func (pl *LCPCListPacket) String() string {
	return "pc list"
}
func (pl *LCPCListPacket) MarshalBinary() ([]byte, error) {
	// ret := make([]byte, 164)

	return []byte{
		83, 79, 86,
		8, 178, 187, 212, 217, 209, 218, 202, 206, 0,
		76, 29, 0, 0, 9, 0, 11, 0, 10, 0, 50, 170, 9, 0, 0, 53, 15, 0, 0, 47, 12, 0, 0, 18, 0, 18, 0, 20, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 144, 1, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 8, 181, 216, 211, 252, 214, 174, 195, 197, 1, 76, 29, 0, 0, 0, 121, 1, 101, 0, 121, 1, 0, 0, 8, 10, 0, 25, 0, 10, 0, 59, 1, 59, 1, 0, 0, 0, 0, 150, 50, 0, 0, 0, 0, 68, 0, 0, 0, 15, 39, 15, 39, 100, 4, 183, 232, 191, 241, 2, 76, 29, 0, 0, 0, 0, 0, 164, 1, 1, 121, 1, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 150, 50, 125, 0, 0, 0, 217, 0, 0, 0, 15, 39, 100,
	}, nil
}

type LCReconnectPacket struct {
	Ip   string
	Port uint16
	Key  []byte
}

func (rc *LCReconnectPacket) Id() packet.PacketID {
	return PACKET_LC_RECONNECT
}
func (rc *LCReconnectPacket) String() string {
	return "reconnect"
}
func (rc *LCReconnectPacket) MarshalBinary() ([]byte, error) {
	//[13 49 57 50 46 49 54 56 46 49 46 49 50 51 14 39 0 0 0 32 6 11]
	sz := 1 + len(rc.Ip) + 2 + 6
	ret := make([]byte, sz)
	ret[0] = byte(len(rc.Ip))
	copy(ret[1:], rc.Ip[:])
	// for i := 1; i < len(rc.Ip); i++ {
	// 	ret[i+1] = rc.Ip[i]
	// }

	binary.LittleEndian.PutUint16(ret[1+len(rc.Ip):], rc.Port)
	copy(ret[3+len(rc.Ip):], rc.Key)
	// for i := 0; i < len(rc.Key); i++ {
	// 	ret[4+len(rc.Ip)+i] = rc.Key[i]
	// }
	return ret, nil
}
