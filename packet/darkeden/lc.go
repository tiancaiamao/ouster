package darkeden

type LCLoginOKPacket struct{}

func (loginOk LCLoginOKPacket) Id() PacketID {
	return PACKET_LC_LOGIN_OK
}

func (loginOk LCLoginOKPacket) String() string {
	return "loginOK"
}

func (loginOk LCLoginOKPacket) Bytes() []byte {
	return []byte{1, 0, 0, 231, 254, 255}
}

type LCVersionCheckOKPacket struct{}

func (v LCVersionCheckOKPacket) Id() PacketID {
	return PACKET_LC_VERSION_CHECK_OK
}

func (v LCVersionCheckOKPacket) String() string {
	return "version check ok"
}
func (v LCVersionCheckOKPacket) Bytes() []byte {
	return []byte{0, 0}
}

type LCWorldListPacket struct{}

func (wl LCWorldListPacket) Id() PacketID {
	return PACKET_LC_WORLD_LIST
}
func (wl LCWorldListPacket) String() string {
	return "world list"
}
func (wl LCWorldListPacket) Bytes() []byte {
	return []byte{2, 1, 1, 1, 8, 185, 237, 247, 200, 193, 182, 211, 252, 0}
}

type LCServerListPacket struct {
	CurrentWorld uint8
	Size         uint8
	list         []string
}

func (sl *LCServerListPacket) Id() PacketID {
	return PACKET_LC_SERVER_LIST
}
func (sl *LCServerListPacket) String() string {
	return "server list"
}
func (sl *LCServerListPacket) Bytes() []byte {
	return []byte{3, 1, 2, 0, 6, 183, 226, 178, 226, 199, 248, 0, 1, 8, 185, 237, 247, 200, 193, 182, 211, 252, 0}
}

type PCInfo struct{}
type LCPCListPacket struct {
	list []PCInfo
}

func (pl *LCPCListPacket) Id() PacketID {
	return PACKET_LC_PC_LIST
}
func (pl *LCPCListPacket) String() string {
	return "pc list"
}
func (pl *LCPCListPacket) Bytes() []byte {
	// ret := make([]byte, 164)

	return []byte{
		4, 83, 79, 86,
		8, 178, 187, 212, 217, 209, 218, 202, 206, 0,
		76, 29, 0, 0, 9, 0, 11, 0, 10, 0, 50, 170, 9, 0, 0, 53, 15, 0, 0, 47, 12, 0, 0, 18, 0, 18, 0, 20, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 144, 1, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 8, 181, 216, 211, 252, 214, 174, 195, 197, 1, 76, 29, 0, 0, 0, 121, 1, 101, 0, 121, 1, 0, 0, 8, 10, 0, 25, 0, 10, 0, 59, 1, 59, 1, 0, 0, 0, 0, 150, 50, 0, 0, 0, 0, 68, 0, 0, 0, 15, 39, 15, 39, 100, 4, 183, 232, 191, 241, 2, 76, 29, 0, 0, 0, 0, 0, 164, 1, 1, 121, 1, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 150, 50, 125, 0, 0, 0, 217, 0, 0, 0, 15, 39, 100,
	}
}
