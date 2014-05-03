package darkeden

type CLLoginPacket struct {
	Username string
	Password string
}

func (login *CLLoginPacket) Id() PacketID {
	return PACKET_CL_LOGIN
}

func (login *CLLoginPacket) String() string {
	return "Login"
}

func readLogin(buf []byte) (Packet, error) {
	buf = buf[1:]
	szUsername := int(buf[0])
	szPassword := int(buf[1+szUsername])
	return &CLLoginPacket{
		Username: string(buf[1 : 1+szUsername]),
		Password: string(buf[2+szUsername : 2+szUsername+szPassword]),
	}, nil
}

type CLVersionCheckPacket struct{}

func (v CLVersionCheckPacket) Id() PacketID {
	return PACKET_CL_VERSION_CHECK
}
func (v CLVersionCheckPacket) String() string {
	return "version check"
}

type CLGetWorldListPacket struct{}

func (worldList CLGetWorldListPacket) Id() PacketID {
	return PACKET_CL_GET_WORLD_LIST
}
func (w CLGetWorldListPacket) String() string {
	return "get world list"
}

type CLSelectWorldPacket uint8

func (sw CLSelectWorldPacket) Id() PacketID {
	return PACKET_CL_SELECT_WORLD
}
func (sw CLSelectWorldPacket) String() string {
	return "select world"
}

func readSelectWorld(buf []byte) (Packet, error) {
	return CLSelectWorldPacket(buf[0]), nil
}

type CLSelectServerPacket uint8

func (ss CLSelectServerPacket) Id() PacketID {
	return PACKET_CL_SELECT_SERVER
}
func (ss CLSelectServerPacket) String() string {
	return "select server"
}
func readSelectServer(buf []byte) (Packet, error) {
	return CLSelectServerPacket(buf[0]), nil
}
