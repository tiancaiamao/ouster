package packet

import (
    "bytes"
    "encoding/binary"
    "errors"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type CGConnectPacket struct {
    NotImplementWrite

    Key        uint32
    PCType     uint8
    PCName     string
    MacAddress [4]byte
}

func (connect *CGConnectPacket) PacketID() PacketID {
    return PACKET_CG_CONNECT
}

func (connect *CGConnectPacket) MarshalBinary(code uint8) ([]byte, error) {
    buf := &bytes.Buffer{}
    binary.Write(buf, binary.LittleEndian, connect.Key)
    binary.Write(buf, binary.LittleEndian, connect.PCType)
    binary.Write(buf, binary.LittleEndian, uint8(len(connect.PCName)))
    io.WriteString(buf, connect.PCName)
    buf.Write(connect.MacAddress[:])
    return buf.Bytes(), nil
}

func (packet *CGConnectPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &packet.Key)
    binary.Read(reader, binary.LittleEndian, &packet.PCType)
    var szName uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &szName)
    _, err := reader.Read(buf[:szName])
    if err != nil {
        return err
    }
    packet.PCName = string(buf[:szName])
    _, err = reader.Read(packet.MacAddress[:])
    if err != nil {
        return err
    }
    return nil
}

type CGReadyPacket struct {
    NotImplementWrite
}

func (ready CGReadyPacket) PacketID() PacketID {
    return PACKET_CG_READY
}
func (ready *CGReadyPacket) Read(reader io.Reader, code uint8) error {
    return nil
}

type CGMovePacket struct {
    NotImplementWrite

    Dir uint8
    X   uint8
    Y   uint8
}

func (move CGMovePacket) PacketID() PacketID {
    return PACKET_CG_MOVE
}

const (
    dirLEFT      = 0
    dirRIGHT     = 4
    dirUP        = 6
    dirDOWN      = 2
    dirLEFTUP    = 7
    dirRIGHTUP   = 5
    dirLEFTDOWN  = 1
    dirRIGHTDOWN = 3
)

func encryptDir(dir byte) (uint8, error) {
    var ret uint8
    switch dir {
    case 53:
        dir = dirLEFT
    case 49:
        dir = dirRIGHT
    case 51:
        dir = dirUP
    case 55:
        dir = dirDOWN
    case 50:
        dir = dirLEFTUP
    case 48:
        dir = dirRIGHTUP
    case 52:
        dir = dirLEFTDOWN
    case 54:
        dir = dirRIGHTDOWN
    default:
        return ret, errors.New("unknow dir")
    }
    return ret, nil
}

func SHUFFLE_STATEMENT_3(code uint8, A func(), B func(), C func()) {
    switch code % 3 {
    case 0:
        A()
        B()
        C()
    case 1:
        B()
        C()
        A()
    case 2:
        C()
        A()
        B()
    }
    return
}

func SHUFFLE_STATEMENT_2(code uint8, A func(), B func()) {
    switch code % 2 {
    case 0:
        A()
        B()
    case 1:
        B()
        A()
    }
    return
}

func SHUFFLE_STATEMENT_4(code uint8, A func(), B func(), C func(), D func()) {
    switch code % 4 {
    case 0:
        A()
        B()
        C()
        D()
    case 1:
        B()
        C()
        D()
        A()
    case 2:
        C()
        D()
        A()
        B()
    case 3:
        D()
        A()
        C()
        B()
    }
    return
}

func (move *CGMovePacket) Read(reader io.Reader, code uint8) error {
    A := func() {
        binary.Read(reader, binary.LittleEndian, &move.X)
        move.X ^= code
    }
    B := func() {
        binary.Read(reader, binary.LittleEndian, &move.Y)
        move.Y ^= code
    }
    C := func() {
        binary.Read(reader, binary.LittleEndian, &move.Dir)
        move.Dir ^= code
    }
    // encryption...fuck
    SHUFFLE_STATEMENT_3(code, A, B, C)
    if move.Dir >= 8 {
        return errors.New("Dir out of range")
    }
    return nil
}

type CGVerifyTimePacket struct {
    NotImplementWrite
}

func (verifyTime CGVerifyTimePacket) PacketID() PacketID {
    return PACKET_CG_VERIFY_TIME
}

func (verifyTime *CGVerifyTimePacket) Read(reader io.Reader, code uint8) error {
    return nil
}

type CGAttackPacket struct {
    NotImplementWrite

    ObjectID ObjectID_t
    X        uint8
    Y        uint8
    Dir      uint8
}

func (attack CGAttackPacket) PacketID() PacketID {
    return PACKET_CG_ATTACK
}

func (attack *CGAttackPacket) Read(reader io.Reader, code uint8) error {
    // [188 251 55 82 48 0 0]
    A := func() {
        binary.Read(reader, binary.LittleEndian, &attack.ObjectID)
        attack.ObjectID ^= ObjectID_t(code)
    }
    B := func() {
        binary.Read(reader, binary.LittleEndian, &attack.X)
        attack.X ^= code
    }
    C := func() {
        binary.Read(reader, binary.LittleEndian, &attack.Y)
        attack.Y ^= code
    }
    D := func() {
        binary.Read(reader, binary.LittleEndian, &attack.Dir)
        attack.Dir ^= code
    }
    SHUFFLE_STATEMENT_4(code, A, B, C, D)
    return nil
}

type CGBloodDrainPacket struct {
    NotImplementWrite

    ObjectID uint32
}

func (bloodDrain CGBloodDrainPacket) PacketID() PacketID {
    return PACKET_CG_BLOOD_DRAIN
}
func (bloodDrain CGBloodDrainPacket) String() string {
    return "blood drain"
}
func (packet *CGBloodDrainPacket) Read(reader io.Reader, code uint8) error {
    return binary.Read(reader, binary.LittleEndian, &packet.ObjectID)
}

type CGLearnSkillPacket struct {
    NotImplementWrite

    SkillType       SkillType_t
    SkillDomainType uint8
}

func (learnSkill CGLearnSkillPacket) PacketID() PacketID {
    return PACKET_CG_LEARN_SKILL
}

func (learnSkill CGLearnSkillPacket) String() string {
    return "learn skill"
}

func (learn *CGLearnSkillPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &learn.SkillType)
    binary.Read(reader, binary.LittleEndian, &learn.SkillDomainType)
    return nil
}

type CGSkillToObjectPacket struct {
    NotImplementWrite

    SkillType      SkillType_t
    CEffectID      uint16
    TargetObjectID ObjectID_t
}

func (skill CGSkillToObjectPacket) PacketID() PacketID {
    return PACKET_CG_SKILL_TO_OBJECT
}

func (skill CGSkillToObjectPacket) String() string {
    return "skill to object"
}

func (ret *CGSkillToObjectPacket) Read(reader io.Reader, code uint8) error {
    A := func() {
        binary.Read(reader, binary.LittleEndian, &ret.SkillType)
        ret.SkillType ^= SkillType_t(code)
    }
    B := func() {
        binary.Read(reader, binary.LittleEndian, &ret.CEffectID)
        ret.CEffectID ^= uint16(code)
    }
    C := func() {
        binary.Read(reader, binary.LittleEndian, &ret.TargetObjectID)
        ret.TargetObjectID ^= ObjectID_t(code)
    }
    SHUFFLE_STATEMENT_3(code, A, B, C)
    return nil
}

type CGSkillToSelfPacket struct {
    NotImplementWrite

    SkillType SkillType_t
    CEffectID uint16
}

func (skill CGSkillToSelfPacket) PacketID() PacketID {
    return PACKET_CG_SKILL_TO_SELF
}

func (skill CGSkillToSelfPacket) String() string {
    return "skill to self"
}

func (ret *CGSkillToSelfPacket) Read(reader io.Reader, code uint8) error {
    A := func() {
        binary.Read(reader, binary.LittleEndian, &ret.SkillType)
        ret.SkillType ^= SkillType_t(code)
    }
    B := func() {
        binary.Read(reader, binary.LittleEndian, &ret.CEffectID)
        ret.CEffectID ^= uint16(code)
    }
    SHUFFLE_STATEMENT_2(code, A, B)
    return nil
}

type CGSkillToTilePacket struct {
    NotImplementWrite

    SkillType SkillType_t
    CEffectID uint16
    X         Coord_t
    Y         Coord_t
}

func (skill CGSkillToTilePacket) PacketID() PacketID {
    return PACKET_CG_SKILL_TO_TILE
}

func (skill CGSkillToTilePacket) String() string {
    return "skill to tile"
}

func (ret *CGSkillToTilePacket) Read(reader io.Reader, code uint8) error {
    A := func() {
        binary.Read(reader, binary.LittleEndian, &ret.SkillType)
        ret.SkillType ^= SkillType_t(code)
    }
    B := func() {
        binary.Read(reader, binary.LittleEndian, &ret.CEffectID)
        ret.CEffectID ^= uint16(code)
    }
    C := func() {
        binary.Read(reader, binary.LittleEndian, &ret.X)
        ret.X ^= Coord_t(code)
    }
    D := func() {
        binary.Read(reader, binary.LittleEndian, &ret.Y)
        ret.Y ^= Coord_t(code)
    }
    SHUFFLE_STATEMENT_4(code, A, B, C, D)
    return nil
}

type CGSayPacket struct {
    NotImplementWrite

    Color   uint32
    Message string
}

func (say *CGSayPacket) PacketID() PacketID {
    return PACKET_CG_SAY
}
func (say *CGSayPacket) String() string {
    return "say"
}
func (say *CGSayPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &say.Color)
    var sz uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &sz)
    _, err := reader.Read(buf[:sz])
    if err != nil {
        return err
    }
    say.Message = string(buf[:sz])
    return nil
}

type CGLogoutPacket struct {
    NotImplementWrite
}

func (_ CGLogoutPacket) PacketID() PacketID {
    return PACKET_CG_LOGOUT
}

func (_ *CGLogoutPacket) Read(reader io.Reader, code uint8) error {
    return nil
}

func (_ CGLogoutPacket) String() string {
    return "logout"
}
