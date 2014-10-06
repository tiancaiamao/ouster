package packet

import (
    "encoding/binary"
    "errors"
    "io"
)

type NotImplementWrite struct{}

func (ign NotImplementWrite) Write(writer io.Writer, code uint8) error {
    return errors.New("not implement write method!!")
}

func (ign NotImplementWrite) PacketSize() uint32 {
    return 0
}

type CLLoginPacket struct {
    NotImplementWrite
    Username   string
    Password   string
    MacAddress [6]byte
    LoginMode  uint8
}

func (login *CLLoginPacket) PacketID() PacketID {
    return PACKET_CL_LOGIN
}

func (login *CLLoginPacket) Read(reader io.Reader, code uint8) error {
    var tmp uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &tmp)
    _, err := reader.Read(buf[:tmp])
    if err != nil {
        return err
    }
    login.Username = string(buf[:tmp])

    binary.Read(reader, binary.LittleEndian, &tmp)
    _, err = reader.Read(buf[:tmp])
    if err != nil {
        return err
    }
    login.Password = string(buf[:tmp])

    reader.Read(login.MacAddress[:])
    binary.Read(reader, binary.LittleEndian, &login.LoginMode)
    return nil
}

type CLVersionCheckPacket struct {
    NotImplementWrite
}

func (ign *CLVersionCheckPacket) Read(reader io.Reader, code uint8) error {
    var buf [4]byte
    reader.Read(buf[:])
    return nil
}

func (v CLVersionCheckPacket) PacketID() PacketID {
    return PACKET_CL_VERSION_CHECK
}

type CLGetWorldListPacket struct {
    NotImplementWrite
}

func (worldList CLGetWorldListPacket) PacketID() PacketID {
    return PACKET_CL_GET_WORLD_LIST
}

func (ign *CLGetWorldListPacket) Read(reader io.Reader, code uint8) error {
    return nil
}

type CLSelectWorldPacket struct {
    NotImplementWrite

    WorldID uint8
}

func (sw CLSelectWorldPacket) PacketID() PacketID {
    return PACKET_CL_SELECT_WORLD
}

func (v *CLSelectWorldPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &v.WorldID)
    return nil
}

type CLSelectServerPacket struct {
    NotImplementWrite
    Data uint8
}

func (ss CLSelectServerPacket) PacketID() PacketID {
    return PACKET_CL_SELECT_SERVER
}

func (v *CLSelectServerPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &v.Data)
    return nil
}

type PCType uint8
type CLSelectPcPacket struct {
    NotImplementWrite
    Name string
    Type PCType
}

func (sp *CLSelectPcPacket) PacketID() PacketID {
    return PACKET_CL_SELECT_PC
}
func (_ *CLSelectPcPacket) PacketSize() uint32 {
    return 0
}

func (ret *CLSelectPcPacket) Read(reader io.Reader, code uint8) error {
    //	[8 178 187 212 217 209 218 202 206 0]
    var sz uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &sz)
    _, err := reader.Read(buf[:sz])
    if err != nil {
        return err
    }
    ret.Name = string(buf[:sz])
    binary.Read(reader, binary.LittleEndian, &ret.Type)
    return nil
}
