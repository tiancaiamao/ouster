package packet

import (
    // "bytes"
    "encoding/binary"
    "errors"
    "github.com/tiancaiamao/ouster/data"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type DummyRead struct{}

func (_ DummyRead) Read(reader io.Reader, code uint8) error {
    return errors.New("dummy reader")
}

type GCMoveOKPacket struct {
    DummyRead

    Dir uint8
    X   uint8
    Y   uint8
}

func (moveOk GCMoveOKPacket) PacketID() PacketID {
    return PACKET_GC_MOVE_OK
}

func (moveOK GCMoveOKPacket) PacketSize() uint32 {
    return 3
}

func (moveOk GCMoveOKPacket) String() string {
    return "move ok"
}
func (moveOk GCMoveOKPacket) Write(buf io.Writer, code uint8) error {
    ret := make([]byte, 3)
    offset := 0
    A := func() {
        ret[offset] = moveOk.X ^ code
        offset++
    }
    B := func() {
        ret[offset] = moveOk.Y ^ code
        offset++
    }
    C := func() {
        ret[offset] = moveOk.Dir ^ code
        offset++
    }
    SHUFFLE_STATEMENT_3(code, A, B, C)
    return nil
}

type GCMoveErrorPacket struct {
    DummyRead

    X   uint8
    Y   uint8
}

func (moveError GCMoveErrorPacket) PacketID() PacketID {
    return PACKET_GC_MOVE_ERROR
}
func (moveError GCMoveErrorPacket) PacketSize() uint32 {
    return 2
}
func (moveError GCMoveErrorPacket) String() string {
    return "move error"
}
func (moveError GCMoveErrorPacket) Write(buf io.Writer, code uint8) error {
    ret := make([]byte, 2)
    offset := 0
    A := func() {
        ret[offset] = moveError.X ^ code
        offset++
    }
    B := func() {
        ret[offset] = moveError.Y ^ code
        offset++
    }
    SHUFFLE_STATEMENT_2(code, A, B)
    return nil
}

type GCMovePacket struct {
    DummyRead

    ObjectID ObjectID_t
    X        Coord_t
    Y        Coord_t
    Dir      Dir_t
}

func (move GCMovePacket) PacketID() PacketID {
    return PACKET_GC_MOVE
}

func (move GCMovePacket) PacketSize() uint32 {
    return 7
}

func (move GCMovePacket) Write(buf io.Writer, code uint8) error {
    ret := []byte{0, 0, 0, 0, byte(move.X), byte(move.Y), byte(move.Dir)}
    binary.LittleEndian.PutUint32(ret[:], uint32(move.ObjectID))
    return nil
}

type NPCType uint16

type RideMotorcycleInfo struct{}

func (info RideMotorcycleInfo) Dump(writer io.Writer) {
    // TODO
    binary.Write(writer, binary.LittleEndian, uint8(0))
}

type GameTimeType struct {
    Year  uint16
    Month uint8
    Day   uint8

    Hour   uint8
    Minute uint8
    Second uint8
}

func (time *GameTimeType) Size() uint32 {
    return 7
}

func (time GameTimeType) Dump(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, time.Year)
    binary.Write(writer, binary.LittleEndian, time.Month)
    binary.Write(writer, binary.LittleEndian, time.Day)
    binary.Write(writer, binary.LittleEndian, time.Hour)
    binary.Write(writer, binary.LittleEndian, time.Minute)
    binary.Write(writer, binary.LittleEndian, time.Second)
    return
}

type GCUpdateInfoPacket struct {
    // 'V'或者'O'或者'S'
    PCType             byte
    PCInfo             data.PCInfo
    InventoryInfo      data.InventoryInfo
    GearInfo           data.GearInfo
    ExtraInfo          data.ExtraInfo
    EffectInfo         data.EffectInfo
    hasMotorcycle      bool
    RideMotorcycleInfo RideMotorcycleInfo

    ZoneID   ZoneID_t
    ZoneX    Coord_t
    ZoneY    Coord_t
    GameTime GameTimeType

    Weather      Weather
    WeatherLevel WeatherLevel_t

    DarkLevel  DarkLevel_t
    LightLevel LightLevel_t

    NPCTypes     []NPCType_t
    MonsterTypes []MonsterType_t

    NPCInfos []data.NPCInfo

    ServerStat   uint8
    Premium      uint8
    SMSCharge    uint32
    NicknameInfo data.NicknameInfo

    NonPK              bool
    GuildUnionID       uint32
    GuildUnionUserType uint8
    BloodBibleSignInfo data.BloodBibleSignInfo
    PowerPoint         uint32
}

func (info *GCUpdateInfoPacket) PacketID() PacketID {
    return PACKET_GC_UPDATE_INFO
}
func (info *GCUpdateInfoPacket) PacketSize() uint32 {
    var sz uint32
    sz = info.PCInfo.Size() +
        info.InventoryInfo.Size() +
        info.GearInfo.Size() +
        info.ExtraInfo.Size() + 1

    if info.hasMotorcycle {
        panic("not implement yet")
        // sz += info.RideMotorcycleInfo.Size()
    }

    sz = sz + 2 + 1 + 1 +
        info.GameTime.Size() + 1 + 1 + 1

    sz = sz + 1 + uint32(len(info.NPCTypes))*2 + 1 + uint32(len(info.MonsterTypes))*2 + 1

    for i := 0; i < len(info.NPCInfos); i++ {
        sz += info.NPCInfos[i].Size()
    }
    sz += 6
    sz += info.BloodBibleSignInfo.Size()
    sz += 4
    return sz
}

func (info *GCUpdateInfoPacket) String() string {
    return "update info"
}
func (info *GCUpdateInfoPacket) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, info.PCType)
    info.PCInfo.Write(buf)

    info.InventoryInfo.Write(buf)
    info.GearInfo.Write(buf)
    info.ExtraInfo.Write(buf)
    info.EffectInfo.Write(buf)
    if info.hasMotorcycle {
        binary.Write(buf, binary.LittleEndian, uint8(1))
        info.RideMotorcycleInfo.Dump(buf)
    } else {
        binary.Write(buf, binary.LittleEndian, uint8(0))
    }

    // write zone info
    binary.Write(buf, binary.LittleEndian, info.ZoneID)
    binary.Write(buf, binary.LittleEndian, info.ZoneX)
    binary.Write(buf, binary.LittleEndian, info.ZoneY)

    info.GameTime.Dump(buf)
    binary.Write(buf, binary.LittleEndian, info.Weather)
    binary.Write(buf, binary.LittleEndian, info.WeatherLevel)
    binary.Write(buf, binary.LittleEndian, info.DarkLevel)
    binary.Write(buf, binary.LittleEndian, info.LightLevel)

    binary.Write(buf, binary.LittleEndian, uint8(len(info.NPCTypes)))
    for i := 0; i < len(info.NPCTypes); i++ {
        binary.Write(buf, binary.LittleEndian, info.NPCTypes[i])
    }

    binary.Write(buf, binary.LittleEndian, uint8(len(info.MonsterTypes)))
    for i := 0; i < len(info.MonsterTypes); i++ {
        binary.Write(buf, binary.LittleEndian, info.MonsterTypes[i])
    }

    binary.Write(buf, binary.LittleEndian, uint8(len(info.NPCInfos)))
    for i := 0; i < len(info.NPCInfos); i++ {
        info.NPCInfos[i].Write(buf)
    }

    binary.Write(buf, binary.LittleEndian, info.ServerStat)
    binary.Write(buf, binary.LittleEndian, info.Premium)
    binary.Write(buf, binary.LittleEndian, info.SMSCharge)

    info.NicknameInfo.Write(buf)

    if info.NonPK {
        binary.Write(buf, binary.LittleEndian, uint8(1))
    } else {
        binary.Write(buf, binary.LittleEndian, uint8(0))
    }

    binary.Write(buf, binary.LittleEndian, info.GuildUnionID)
    binary.Write(buf, binary.LittleEndian, info.GuildUnionUserType)

    info.BloodBibleSignInfo.Write(buf)

    binary.Write(buf, binary.LittleEndian, info.PowerPoint)

    return nil

    // return

    // buf.Write([]byte{190, 7,
    // 	 3,
    // 	 19,
    // 	 16,
    // 	10,
    // 	40,
    //
    // 	0, Weather
    // 	0, WeatherLevel
    // 	13, DarkLevel
    // 	2, LightLevel
    //
    // 	0, nNPCS
    //
    // 	5, nMonsters
    // 	9, 0, monsterTypes
    // 	61, 0,
    // 	62, 0,
    // 	64, 0,
    // 	163, 0,
    //
    // 	0, NPCInfoCount
    // 	0, ServerStat
    // 	17, Premium
    // 	0, 0, 0, 0, SMSCharge
    // 	NickNameInfo
    // 	24, 125, 0,
    //
    // 	0, NonPK
    // 	0, 0, 0, 0, GuildUnionID
    // 	2, GuildUnionUserType
    // 	1, 0, 0, 0, 0,
    //  0, 0, 0, 0})	PowerPoint
    // return nil

    // return []byte{86, 117, 48, 0, 0, 4, 183, 232, 191, 241, 150, 0, 0, 0, 164, 1, 0, 76, 29, 0, 0,
    // 		20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 20, 0, 216, 1, 216, 1, 50, 204, 41, 0, 0, 125, 0, 0,
    // 		0, 0, 0, 0, 0, 26, 1, 0, 0, 13, 15, 39, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 4, 0, 0,
    // 		0, 0, 100, 0, 0, 0, 0, 6, 118, 48, 0, 0, 30, 0, 0, 0, 232, 3, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 119,
    // 		48, 0, 0, 44, 0, 0, 2, 16, 1, 136, 19, 0, 0, 0, 0, 3, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 120, 48, 0, 0, 34, 5, 0, 0,
    // 		1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 8, 0, 0, 0, 0, 1, 121, 48, 0, 0, 32, 0, 0, 2, 53, 43, 232, 3, 0, 0,
    // 		0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 122, 48, 0, 0, 32, 1, 0, 0, 232, 3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 3, 123, 48, 0,
    // 		0, 44, 0, 0, 2, 58, 38, 32, 28, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4, 0, 0, 2, 146, 1, 54, 66, 109, 0, 246, 224, 0, 21, 0, 145, 237, 190, 7, 3, 19, 16,
    // 		10, 40, 0, 0, 13, 2, 0, 5, 9, 0, 61, 0, 62, 0, 64, 0, 163, 0, 0, 0, 17, 0, 0, 0, 0, 24, 125, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0},
    // 	nil
}

func (update *GCUpdateInfoPacket) Read(reader io.Reader, code uint8) error {
    binary.Read(reader, binary.LittleEndian, &update.PCType)
    switch update.PCType {
    case 'S':
        update.PCInfo = &data.PCSlayerInfo{}
    case 'V':
        update.PCInfo = &data.PCVampireInfo{}
    case 'O':
        update.PCInfo = &data.PCOusterInfo{}
    default:
        return errors.New("invalid pc type")
    }
    err := update.PCInfo.Read(reader)
    if err != nil {
        return err
    }

    update.InventoryInfo.Read(reader)
    update.GearInfo.Read(reader)
    update.ExtraInfo.Read(reader)
    update.EffectInfo.Read(reader)

    binary.Read(reader, binary.LittleEndian, &update.hasMotorcycle)
    if update.hasMotorcycle {
        return errors.New("not implement yet!!!")
    }

    binary.Read(reader, binary.LittleEndian, &update.ZoneID)
    binary.Read(reader, binary.LittleEndian, &update.ZoneX)
    binary.Read(reader, binary.LittleEndian, &update.ZoneY)

    binary.Read(reader, binary.LittleEndian, &update.GameTime)

    binary.Read(reader, binary.LittleEndian, &update.Weather)
    binary.Read(reader, binary.LittleEndian, &update.WeatherLevel)

    binary.Read(reader, binary.LittleEndian, &update.DarkLevel)
    binary.Read(reader, binary.LittleEndian, &update.LightLevel)

    var nNPC uint8
    binary.Read(reader, binary.LittleEndian, &nNPC)
    update.NPCTypes = make([]NPCType_t, nNPC)
    for i := 0; i < int(nNPC); i++ {
        binary.Read(reader, binary.LittleEndian, &update.NPCTypes[i])
    }

    var nMonster uint8
    binary.Read(reader, binary.LittleEndian, &nMonster)
    update.MonsterTypes = make([]MonsterType_t, nMonster)
    for i := 0; i < int(nMonster); i++ {
        binary.Read(reader, binary.LittleEndian, &update.MonsterTypes[i])
    }

    var nNPCInfo uint8
    binary.Read(reader, binary.LittleEndian, &nNPCInfo)
    update.NPCInfos = make([]data.NPCInfo, nNPCInfo)
    for i := 0; i < int(nNPCInfo); i++ {
        // TODO
        update.NPCInfos[i].Read(reader)
    }

    binary.Read(reader, binary.LittleEndian, &update.ServerStat)
    binary.Read(reader, binary.LittleEndian, &update.Premium)
    binary.Read(reader, binary.LittleEndian, &update.SMSCharge)

    update.NicknameInfo.Read(reader)

    binary.Read(reader, binary.LittleEndian, &update.NonPK)
    binary.Read(reader, binary.LittleEndian, &update.GuildUnionID)
    binary.Read(reader, binary.LittleEndian, &update.GuildUnionUserType)

    update.BloodBibleSignInfo.Read(reader)
    binary.Read(reader, binary.LittleEndian, &update.PowerPoint)
    return nil
}

type GCPetInfoPacket struct {
    DummyRead

    ObjectID   uint32
    PetInfo    []struct{}
    SummonInfo uint8
}

func (pet *GCPetInfoPacket) PacketID() PacketID {
    return PACKET_GC_PET_INFO
}
func (pet *GCPetInfoPacket) PacketSize() uint32 {
    return 5
}
func (pet *GCPetInfoPacket) String() string {
    return "pet info"
}
func (pet *GCPetInfoPacket) Write(buf io.Writer, code uint8) error {
    buf.Write([]byte{0, 117, 48, 0, 0})
    return nil
}

type GCSetPositionPacket struct {
    DummyRead

    X   uint8
    Y   uint8
    Dir uint8
}

func (setPosition GCSetPositionPacket) PacketID() PacketID {
    return PACKET_GC_SET_POSITION
}
func (setPosition GCSetPositionPacket) PacketSize() uint32 {
    return 24
}
func (setPosition GCSetPositionPacket) String() string {
    return "set position"
}
func (setPosition GCSetPositionPacket) Write(buf io.Writer, code uint8) error {
    buf.Write([]byte{setPosition.X, setPosition.Y, setPosition.Dir})
    return nil
}

type GCAddBat struct {
    DummyRead

    ObjectID    uint32
    MonsterName string
    ItemType    uint16
    X           uint8
    Y           uint8
    Dir         uint8
    CurrentHP   uint16
    MaxHP       uint16
    GuildID     uint16
    Color       uint16
}

func (bat *GCAddBat) PacketID() PacketID {
    return PACKET_GC_ADD_BAT
}
func (bat *GCAddBat) PacketSize() uint32 {
    return 4 + 1 + uint32(len(bat.MonsterName)) + 2 + 1 + 1 + 1 + 2 + 2 + 2 + 2
}
func (bat *GCAddBat) String() string {
    return "add bat"
}
func (bat *GCAddBat) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, bat.ObjectID)
    binary.Write(buf, binary.LittleEndian, uint8(len(bat.MonsterName)))
    io.WriteString(buf, bat.MonsterName)
    binary.Write(buf, binary.LittleEndian, bat.ItemType)
    binary.Write(buf, binary.LittleEndian, bat.X)
    binary.Write(buf, binary.LittleEndian, bat.Y)
    binary.Write(buf, binary.LittleEndian, bat.Dir)
    binary.Write(buf, binary.LittleEndian, bat.CurrentHP)
    binary.Write(buf, binary.LittleEndian, bat.MaxHP)
    binary.Write(buf, binary.LittleEndian, bat.GuildID)
    binary.Write(buf, binary.LittleEndian, bat.Color)
    return nil
}

type GCAddMonsterFromBurrowing struct {
    ObjectID    uint32
    MonsterType uint16
    MonsterName string
    MainColor   uint16
    SubColor    uint16
    X           uint8
    Y           uint8
    Dir         uint8
    EffectInfo  []data.EffectInfo
    CurrentHP   uint16
    MaxHP       uint16
}

func (monster *GCAddMonsterFromBurrowing) PacketID() PacketID {
    return PACKET_GC_ADD_MONSTER_FROM_BURROWING
}
func (monster *GCAddMonsterFromBurrowing) PacketSize() uint32 {
    sz := 4 + 2 + 1 + uint32(len(monster.MonsterName)) + 2 + 2 + 1 + 1 + 1 + 2 + 2
    for i := 0; i < len(monster.EffectInfo); i++ {
        sz += monster.EffectInfo[i].Size()
    }
    return sz
}

func (monster *GCAddMonsterFromBurrowing) String() string {
    return "add monster from burrowing"
}
func (monster *GCAddMonsterFromBurrowing) Write(buf io.Writer, code uint8) error {
    buf.Write([]byte{62, 48, 0, 0, 213, 0, 8, 185, 197, 181, 194, 203, 185, 182, 161, 53, 0, 0, 0, 137, 238, 0, 0, 54, 1, 54, 1})
    return nil
}

type GCAddMonster struct {
    DummyRead

    ObjectID    ObjectID_t
    MonsterType MonsterType_t
    MonsterName string
    MainColor   Color_t
    SubColor    Color_t
    X           Coord_t
    Y           Coord_t
    Dir         Dir_t
    EffectInfo  []data.EffectInfo
    CurrentHP   HP_t
    MaxHP       HP_t
    FromFlag    byte
}

func (monster *GCAddMonster) PacketID() PacketID {
    return PACKET_GC_ADD_MONSTER
}

func (monster *GCAddMonster) PacketSize() uint32 {
    sz := 4 + 2 + 1 + uint32(len(monster.MonsterName)) + 2 + 2 + 1 + 1 + 1 + 2 + 2 + 1
    for i := 0; i < len(monster.EffectInfo); i++ {
        sz += monster.EffectInfo[i].Size()
    }
    return sz
}

func (monster *GCAddMonster) String() string {
    return "add monster"
}
func (monster *GCAddMonster) Write(buf io.Writer, code uint8) error {
    //[218 47 0 0 223 0 6 196 218 185 254 203 185 7 0 174 0 102 79 5 0 133 0 133 0 0]
    //[166 47 0 0 72 0 4 192 188 197 181 5 137 133 0 164 214 6 0 156 0 156 0 0]
    //[24 47 0 0 8 0 10 203 185 196 170 191 203 206 172 198 230 53 48 48 58 137 192 6 0 156 0 156 0 0]

    binary.Write(buf, binary.LittleEndian, monster.ObjectID)
    binary.Write(buf, binary.LittleEndian, monster.MonsterType)
    binary.Write(buf, binary.LittleEndian, uint8(len(monster.MonsterName)))
    io.WriteString(buf, monster.MonsterName)
    binary.Write(buf, binary.LittleEndian, monster.MainColor)
    binary.Write(buf, binary.LittleEndian, monster.SubColor)
    binary.Write(buf, binary.LittleEndian, monster.X)
    binary.Write(buf, binary.LittleEndian, monster.Y)
    binary.Write(buf, binary.LittleEndian, monster.Dir)
    binary.Write(buf, binary.LittleEndian, uint8(0))
    binary.Write(buf, binary.LittleEndian, monster.CurrentHP)
    binary.Write(buf, binary.LittleEndian, monster.MaxHP)
    binary.Write(buf, binary.LittleEndian, monster.FromFlag)
    return nil
}

type GCStatusCurrentHP struct {
    DummyRead

    ObjectID  ObjectID_t
    CurrentHP HP_t
}

func (status GCStatusCurrentHP) PacketID() PacketID {
    return PACKET_GC_STATUS_CURRENT_HP
}
func (status GCStatusCurrentHP) PacketSize() uint32 {
    return 6
}

func (status GCStatusCurrentHP) String() string {
    return "status current HP"
}
func (status GCStatusCurrentHP) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, status.ObjectID)
    binary.Write(buf, binary.LittleEndian, status.CurrentHP)
    return nil
}

// 发送给攻击者的，告诉他攻击成功
// ObjectID是被攻击目标的ObjectID
// ModifyInfo是攻击者自身的耗蓝
type GCAttackMeleeOK1 struct {
    DummyRead

    ModifyInfo
    ObjectID ObjectID_t
}

func (ok1 GCAttackMeleeOK1) PacketSize() uint32 {
    return 4 + ok1.ModifyInfo.Size()
}

func (attackOk GCAttackMeleeOK1) PacketID() PacketID {
    return PACKET_GC_ATTACK_MELEE_OK_1
}

func (attackOk GCAttackMeleeOK1) String() string {
    return "attack melee ok 1"
}
func (attackOk GCAttackMeleeOK1) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, attackOk.ObjectID)
    attackOk.Dump(buf)
    return nil
}

// 发送给被攻击者的，告诉他被攻击了
// ObjectID是攻击者的ObjectID
type GCAttackMeleeOK2 struct {
    DummyRead

    ObjectID ObjectID_t
    ModifyInfo
}

func (attackOk GCAttackMeleeOK2) PacketID() PacketID {
    return PACKET_GC_ATTACK_MELEE_OK_2
}
func (ok GCAttackMeleeOK2) PacketSize() uint32 {
    return 4 + ok.ModifyInfo.Size()
}

func (attackOk GCAttackMeleeOK2) String() string {
    return "attack melee ok 2"
}
func (attackOk GCAttackMeleeOK2) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, attackOk.ObjectID)
    attackOk.Dump(buf)
    return nil
}

// 广播
// ObjectID是攻击者，Target是被攻击者
type GCAttackMeleeOK3 struct {
    DummyRead

    ObjectID       ObjectID_t
    TargetObjectID ObjectID_t
}

func (ok GCAttackMeleeOK3) PacketID() PacketID {
    return PACKET_GC_ATTACK_MELEE_OK_3
}

func (ok GCAttackMeleeOK3) PacketSize() uint32 {
    return 8
}

func (ok GCAttackMeleeOK3) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.TargetObjectID)
    return nil
}

type GCCannotUsePacket struct {
    DummyRead
    Data uint32
}

func (cannot GCCannotUsePacket) PacketSize() uint32 {
    return 4
}

func (cannot GCCannotUsePacket) PacketID() PacketID {
    return PACKET_GC_CANNOT_USE
}
func (cannot GCCannotUsePacket) String() string {
    return "cannot use"
}
func (cannot GCCannotUsePacket) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, cannot.Data)
    return nil
}

type ModifyType byte

const (
    MODIFY_BASIC_STR ModifyType = iota
    MODIFY_CURRENT_STR
    MODIFY_MAX_STR
    MODIFY_STR_EXP
    MODIFY_BASIC_DEX
    MODIFY_CURRENT_DEX
    MODIFY_MAX_DEX
    MODIFY_DEX_EXP
    MODIFY_BASIC_INT
    MODIFY_CURRENT_INT
    MODIFY_MAX_INT
    MODIFY_INT_EXP
    MODIFY_CURRENT_HP
    MODIFY_MAX_HP
    MODIFY_CURRENT_MP
    MODIFY_MAX_MP
    MODIFY_MIN_DAMAGE
    MODIFY_MAX_DAMAGE
    MODIFY_DEFENSE
    MODIFY_PROTECTION
    MODIFY_TOHIT
    MODIFY_VISION
    MODIFY_FAME
    MODIFY_GOLD
    MODIFY_SWORD_DOMAIN_LEVEL
    MODIFY_SWORD_DOMAIN_EXP
    MODIFY_SWORD_DOMAIN_GOAL_EXP
    MODIFY_BLADE_DOMAIN_LEVEL
    MODIFY_BLADE_DOMAIN_EXP
    MODIFY_BLADE_DOMAIN_GOAL_EXP
    MODIFY_HEAL_DOMAIN_LEVEL
    MODIFY_HEAL_DOMAIN_EXP
    MODIFY_HEAL_DOMAIN_GOAL_EXP
    MODIFY_ENCHANT_DOMAIN_LEVEL
    MODIFY_ENCHANT_DOMAIN_EXP
    MODIFY_ENCHANT_DOMAIN_GOAL_EXP
    MODIFY_GUN_DOMAIN_LEVEL
    MODIFY_GUN_DOMAIN_EXP
    MODIFY_GUN_DOMAIN_GOAL_EXP
    MODIFY_ETC_DOMAIN_LEVEL
    MODIFY_ETC_DOMAIN_EXP
    MODIFY_ETC_DOMAIN_GOAL_EXP
    MODIFY_SKILL_LEVEL
    MODIFY_LEVEL
    MODIFY_EFFECT_STAT
    MODIFY_DURATION
    MODIFY_BULLET
    MODIFY_BONUS_POINT
    MODIFY_DURABILITY
    MODIFY_NOTORIETY
    MODIFY_VAMP_GOAL_EXP
    MODIFY_SILVER_DAMAGE
    MODIFY_ATTACK_SPEED
    MODIFY_ALIGNMENT
    MODIFY_SILVER_DURABILITY
    MODIFY_REGEN_RATE
    MODIFY_GUILDID
    MODIFY_RANK
    MODIFY_RANK_EXP
    MODIFY_OUSTERS_GOAL_EXP
    MODIFY_SKILL_BONUS_POINT
    MODIFY_ELEMENTAL_FIRE
    MODIFY_ELEMENTAL_WATER
    MODIFY_ELEMENTAL_EARTH
    MODIFY_ELEMENTAL_WIND
    MODIFY_SKILL_EXP
    MODIFY_PET_HP
    MODIFY_PET_EXP
    MODIFY_LAST_TARGET
    MODIFY_UNIONID
    MODIFY_UNIONGRADE
    MODIFY_ADVANCEMENT_CLASS_LEVEL
    MODIFY_ADVANCEMENT_CLASS_GOAL_EXP
    MODIFY_MAX
)

var ModifyType2String []string = []string{
    "BASIC_STR",
    "CURRENT_STR",
    "MAX_STR",
    "STR_EXP",
    "BASIC_DEX",
    "CURRENT_DEX",
    "MAX_DEX",
    "DEX_EXP",
    "BASIC_INT",
    "CURRENT_INT",
    "MAX_INT",
    "INT_EXP",
    "CURRENT_HP",
    "MAX_HP",
    "CURRENT_MP",
    "MAX_MP",
    "MIN_DAMAGE",
    "MAX_DAMAGE",
    "DEFENSE",
    "PROTECTION",
    "TOHIT",
    "VISION",
    "FAME",
    "GOLD",
    "SWORD_DOMAIN_LEVEL",
    "SWORD_DOMAIN_EXP",
    "SWORD_DOMAIN_GOAL_EXP",
    "BLADE_DOMAIN_LEVEL",
    "BLADE_DOMAIN_EXP",
    "BLADE_DOMAIN_GOAL_EXP",
    "HEAL_DOMAIN_LEVEL",
    "HEAL_DOMAIN_EXP",
    "HEAL_DOMAIN_GOAL_EXP",
    "ENCHANT_DOMAIN_LEVEL",
    "ENCHANT_DOMAIN_EXP",
    "ENCHANT_DOMAIN_GOAL_EXP",
    "GUN_DOMAIN_LEVEL",
    "GUN_DOMAIN_EXP",
    "GUN_DOMAIN_GOAL_EXP",
    "ETC_DOMAIN_LEVEL",
    "ETC_DOMAIN_EXP",
    "ETC_DOMAIN_GOAL_EXP",
    "SKILL_LEVEL",
    "LEVEL",
    "EFFECT_STAT",
    "DURATION",
    "BULLET",
    "BONUS_POINT",
    "DURABILITY",
    "NOTORIETY",
    "VAMP_EXP",
    "SILVER_DAMAGE",
    "ATTACK_SPEED",
    "ALIGNMENT",
    "SILVER_DURABILITY",
    "REGEN_RATE",
    "GUILDID",
    "RANK",
    "RANK_EXP",
    "MODIFY_OUSTERS_EXP",
    "MODIFY_SKILL_BONUS_POINT",
    "MODIFY_ELEMENTAL_FIRE",
    "MODIFY_ELEMENTAL_WATER",
    "MODIFY_ELEMENTAL_EARTH",
    "MODIFY_ELEMENTAL_WIND",
    "MODIFY_SKILL_EXP",
    "MODIFY_PET_HP",
    "MODIFY_PET_EXP",
    "MODIFY_LAST_TARGET",
    "MODIFY_UNIONID",
    "MODIFY_UNIONGRADE",
    "MODIFY_ADVANCEMENT_CLASS_LEVEL",
    "MODIFY_ADVANCEMENT_CLASS_GOAL_EXP",
    "MAX",
}

type ModifyInfo struct {
    Short map[ModifyType]uint16
    Long  map[ModifyType]uint32
}

func (modify *ModifyInfo) Dump(writer io.Writer) {
    szShort := uint8(len(modify.Short))
    binary.Write(writer, binary.LittleEndian, szShort)
    for k, v := range modify.Short {
        binary.Write(writer, binary.LittleEndian, k)
        binary.Write(writer, binary.LittleEndian, v)
    }

    szLong := uint8(len(modify.Long))
    binary.Write(writer, binary.LittleEndian, szLong)
    for k, v := range modify.Long {
        binary.Write(writer, binary.LittleEndian, k)
        binary.Write(writer, binary.LittleEndian, v)
    }
}

func (info *ModifyInfo) Size() uint32 {
    return uint32(2 + len(info.Short)*3 + len(info.Long)*5)
}

type GCBloodDrainOK1 struct {
    DummyRead

    Modify   ModifyInfo
    ObjectID uint32
}

func (bdo *GCBloodDrainOK1) PacketID() PacketID {
    return PACKET_GC_BLOOD_DRAIN_OK_1
}
func (bdo *GCBloodDrainOK1) PacketSize() uint32 {
    return 4 + bdo.Modify.Size()
}
func (bdo *GCBloodDrainOK1) String() string {
    return "blood drain ok 1"
}
func (bdo *GCBloodDrainOK1) Write(buf io.Writer, code uint8) error {
    // 237, 53, 0, 0, 2, 51, 0, 0, 12, 216, 1, 0}

    binary.Write(buf, binary.LittleEndian, bdo.ObjectID)
    bdo.Modify.Dump(buf)
    return nil
}

type GCModifyInformationPacket ModifyInfo

func (modify *GCModifyInformationPacket) PacketID() PacketID {
    return PACKET_GC_MODIFY_INFORMATION
}
func (modify *GCModifyInformationPacket) String() string {
    return "modify information"
}
func (modify *GCModifyInformationPacket) Write(buf io.Writer, code uint8) error {

    raw := (*ModifyInfo)(modify)
    raw.Dump(buf)
    return nil
}

type GCAddEffect struct {
    DummyRead

    ObjectID ObjectID_t
    EffectID EffectID_t
    Duration Duration_t
}

func (effect GCAddEffect) PacketID() PacketID {
    return PACKET_GC_ADD_EFFECT
}
func (effect GCAddEffect) PacketSize() uint32 {
    return 8
}
func (effect GCAddEffect) String() string {
    return "add effect"
}
func (effect GCAddEffect) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, effect.ObjectID)
    binary.Write(buf, binary.LittleEndian, effect.EffectID)
    binary.Write(buf, binary.LittleEndian, effect.Duration)
    return nil
}

type GCAddMonsterCorpse struct {
    DummyRead

    ObjectID    ObjectID_t
    MonsterType MonsterType_t
    MonsterName string

    X       Coord_t
    Y       Coord_t
    Dir     Dir_t
    HasHead bool

    TreasureCount uint8
    LastKiller    ObjectID_t
}

func (corpse *GCAddMonsterCorpse) PacketID() PacketID {
    return PACKET_GC_ADD_MONSTER_CORPSE
}
func (corpse *GCAddMonsterCorpse) PacketSize() uint32 {
    return uint32(6 + 1 + len(corpse.MonsterName) + 5 + 4)
}
func (corpse *GCAddMonsterCorpse) String() string {
    return "add monster corpse"
}
func (corpse *GCAddMonsterCorpse) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, corpse.ObjectID)
    binary.Write(buf, binary.LittleEndian, corpse.MonsterType)
    binary.Write(buf, binary.LittleEndian, uint8(len(corpse.MonsterName)))
    if len(corpse.MonsterName) > 0 {
        io.WriteString(buf, corpse.MonsterName)
    }
    binary.Write(buf, binary.LittleEndian, corpse.X)
    binary.Write(buf, binary.LittleEndian, corpse.Y)
    binary.Write(buf, binary.LittleEndian, corpse.Dir)
    if corpse.HasHead {
        binary.Write(buf, binary.LittleEndian, uint8(1))
    } else {
        binary.Write(buf, binary.LittleEndian, uint8(0))
    }

    binary.Write(buf, binary.LittleEndian, corpse.TreasureCount)
    binary.Write(buf, binary.LittleEndian, corpse.LastKiller)
    return nil
}

type GCCreatureDiedPacket struct {
    DummyRead
    ObjectID ObjectID_t
}

func (died GCCreatureDiedPacket) PacketID() PacketID {
    return PACKET_GC_CREATURE_DIED
}
func (died GCCreatureDiedPacket) PacketSize() uint32 {
    return 4
}
func (died GCCreatureDiedPacket) String() string {
    return "creature died"
}
func (died GCCreatureDiedPacket) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, died.ObjectID)
    return nil
}

type GCDeleteObjectPacket struct {
    DummyRead
    ObjectID ObjectID_t
}

func (obj GCDeleteObjectPacket) PacketID() PacketID {
    return PACKET_GC_DELETE_OBJECT
}
func (obj GCDeleteObjectPacket) PacketSize() uint32 {
    return 4
}
func (obj GCDeleteObjectPacket) String() string {
    return "delete object"
}
func (obj *GCDeleteObjectPacket) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, obj.ObjectID)
    return nil
}

type GCAddEffectPacket struct {
    DummyRead

    ObjectID uint32
    EffectID uint16
    Duration uint16
}

func (obj GCAddEffectPacket) PacketID() PacketID {
    return PACKET_GC_ADD_EFFECT
}
func (obj GCAddEffectPacket) PacketSize() uint32 {
    return 8
}
func (obj GCAddEffectPacket) String() string {
    return "add effect"
}
func (obj GCAddEffectPacket) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, obj.ObjectID)
    binary.Write(buf, binary.LittleEndian, obj.EffectID)
    binary.Write(buf, binary.LittleEndian, obj.Duration)
    return nil
}

type GCFastMovePacket struct {
    DummyRead

    ObjectID  ObjectID_t
    FromX     Coord_t
    FromY     Coord_t
    ToX       Coord_t
    ToY       Coord_t
    SkillType SkillType_t
}

func (fastMove *GCFastMovePacket) PacketID() PacketID {
    return PACKET_GC_FAST_MOVE
}
func (fastMove *GCFastMovePacket) PacketSize() uint32 {
    return 10
}
func (fastMove *GCFastMovePacket) String() string {
    return "fast move"
}
func (fastMove *GCFastMovePacket) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, fastMove.ObjectID)
    binary.Write(buf, binary.LittleEndian, fastMove.FromX)
    binary.Write(buf, binary.LittleEndian, fastMove.FromY)
    binary.Write(buf, binary.LittleEndian, fastMove.ToX)
    binary.Write(buf, binary.LittleEndian, fastMove.ToY)
    binary.Write(buf, binary.LittleEndian, fastMove.SkillType)
    return nil
}

type GCLearnSkillOK struct {
    DummyRead

    SkillType       uint16
    SkillDomainType uint8
}

func (ok *GCLearnSkillOK) PacketID() PacketID {
    return PACKET_GC_LEARN_SKILL_OK
}
func (ok *GCLearnSkillOK) PacketSize() uint32 {
    return 3
}
func (ok *GCLearnSkillOK) String() string {
    return "learn skill ok"
}
func (ok *GCLearnSkillOK) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.SkillDomainType)
    return nil
}

type GCRemoveEffect struct {
    DummyRead

    ObjectID   ObjectID_t
    EffectList []uint16
}

func (remove GCRemoveEffect) PacketID() PacketID {
    return PACKET_GC_REMOVE_EFFECT
}
func (remove GCRemoveEffect) PacketSize() uint32 {
    return uint32(4 + 1 + 2*len(remove.EffectList))
}
func (remove GCRemoveEffect) String() string {
    return "remove effect"
}
func (remove *GCRemoveEffect) Write(buf io.Writer, code uint8) error {
    binary.Write(buf, binary.LittleEndian, remove.ObjectID)
    binary.Write(buf, binary.LittleEndian, uint8(len(remove.EffectList)))
    for _, v := range remove.EffectList {
        binary.Write(buf, binary.LittleEndian, v)
    }

    return nil
}

type GCSkillFailed1Packet struct {
    DummyRead

    ModifyInfo

    SkillType SkillType_t
    Grade     uint8
}

func (failed *GCSkillFailed1Packet) PacketID() PacketID {
    return PACKET_GC_SKILL_FAILED_1
}
func (failed *GCSkillFailed1Packet) PacketSize() uint32 {
    return 3 + failed.ModifyInfo.Size()
}
func (failed *GCSkillFailed1Packet) String() string {
    return "skill failed1"
}
func (failed *GCSkillFailed1Packet) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, failed.SkillType)
    binary.Write(buf, binary.LittleEndian, failed.Grade)
    failed.Dump(buf)
    return nil
}

type GCSkillFailed2 struct {
    DummyRead

    ObjectID       uint32
    TargetObjectID uint32
    SkillType      uint16
    Grade          uint8
}

func (failed GCSkillFailed2) PacketID() PacketID {
    return PACKET_GC_SKILL_FAILED_2
}
func (failed GCSkillFailed2) PacketSize() uint32 {
    return 11
}
func (failed GCSkillFailed2) String() string {
    return "skill failed1"
}
func (failed GCSkillFailed2) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, failed.ObjectID)
    binary.Write(buf, binary.LittleEndian, failed.TargetObjectID)
    binary.Write(buf, binary.LittleEndian, failed.SkillType)
    binary.Write(buf, binary.LittleEndian, failed.Grade)
    return nil
}

// send to player
type GCSkillToObjectOK1 struct {
    DummyRead

    SkillType      SkillType_t
    CEffectID      uint16
    TargetObjectID ObjectID_t
    Duration       uint16
    Grade          uint8
    ModifyInfo
}

func (ok *GCSkillToObjectOK1) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_OBJECT_OK_1
}
func (ok *GCSkillToObjectOK1) PacketSize() uint32 {
    return 11 + ok.ModifyInfo.Size()
}
func (ok *GCSkillToObjectOK1) String() string {
    return "skill to object ok 1"
}
func (ok *GCSkillToObjectOK1) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.CEffectID)
    binary.Write(buf, binary.LittleEndian, ok.TargetObjectID)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    ok.Dump(buf)
    return nil
}

type GCSkillToObjectOK3 struct {
    DummyRead

    ObjectID  uint32
    SkillType uint16
    TargetX   uint8
    TargetY   uint8
    Grade     uint8
}

func (ok *GCSkillToObjectOK3) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_OBJECT_OK_3
}
func (ok *GCSkillToObjectOK3) PacketSize() uint32 {
    return 9
}
func (ok *GCSkillToObjectOK3) String() string {
    return "skill to object ok 3"
}
func (ok *GCSkillToObjectOK3) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.TargetX)
    binary.Write(buf, binary.LittleEndian, ok.TargetY)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    return nil
}

// send to passive
type GCSkillToObjectOK4 struct {
    DummyRead

    ObjectID  uint32
    SkillType uint16
    Duration  uint16
    Grade     uint8
}

func (ok GCSkillToObjectOK4) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_OBJECT_OK_4
}
func (ok GCSkillToObjectOK4) PacketSize() uint32 {
    return 9
}
func (ok GCSkillToObjectOK4) String() string {
    return "skill to object ok 4"
}
func (ok GCSkillToObjectOK4) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    return nil
}

type GCSkillToSelfOK1 struct {
    DummyRead

    SkillType SkillType_t
    CEffectID uint16
    Duration  uint16
    Grade     uint8
    ModifyInfo
}

func (ok *GCSkillToSelfOK1) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_SELF_OK_1
}
func (ok *GCSkillToSelfOK1) PacketSize() uint32 {
    return 7 + ok.ModifyInfo.Size()
}
func (ok *GCSkillToSelfOK1) String() string {
    return "skill to self ok 1"
}
func (ok *GCSkillToSelfOK1) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.CEffectID)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    ok.Dump(buf)
    return nil
}

type GCSkillToSelfOK2 struct {
    DummyRead

    ObjectID  ObjectID_t
    SkillType SkillType_t
    Duration  Duration_t
    Grade     byte
}

func (ok2 *GCSkillToSelfOK2) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_SELF_OK_2
}
func (ok2 *GCSkillToSelfOK2) PacketSize() uint32 {
    return 9
}
func (ok *GCSkillToSelfOK2) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    return nil
}

type GCSkillToTileOK1 struct {
    DummyRead

    SkillType    SkillType_t
    CEffectID    uint16
    Duration     Duration_t
    Range        Range_t
    X            Coord_t
    Y            Coord_t
    CreatureList []ObjectID_t
    Grade        uint8
    ModifyInfo
}

func (ok *GCSkillToTileOK1) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_TILE_OK_1
}
func (ok *GCSkillToTileOK1) PacketSize() uint32 {
    return 10 + 1 + uint32(len(ok.CreatureList)*4) + ok.ModifyInfo.Size()
}
func (ok *GCSkillToTileOK1) String() string {
    return "skill to tile ok 1"
}
func (ok *GCSkillToTileOK1) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.CEffectID)
    binary.Write(buf, binary.LittleEndian, ok.X)
    binary.Write(buf, binary.LittleEndian, ok.Y)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Range)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    binary.Write(buf, binary.LittleEndian, uint8(len(ok.CreatureList)))
    for _, v := range ok.CreatureList {
        binary.Write(buf, binary.LittleEndian, v)
    }
    ok.Dump(buf)

    return nil
}

type GCSkillToTileOK2 struct {
    DummyRead

    ModifyInfo

    ObjectID  ObjectID_t
    SkillType SkillType_t
    X         Coord_t
    Y         Coord_t
    Range     Range_t
    Duration  Duration_t
    CList     []ObjectID_t
    Grade     byte
}

func (ok GCSkillToTileOK2) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_TILE_OK_2
}
func (ok GCSkillToTileOK2) PacketSize() uint32 {
    return 11 + 1 + uint32(4*len(ok.CList)) + 1 + ok.ModifyInfo.Size()
}
func (ok *GCSkillToTileOK2) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.X)
    binary.Write(buf, binary.LittleEndian, ok.Y)
    binary.Write(buf, binary.LittleEndian, ok.Range)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    binary.Write(buf, binary.LittleEndian, uint8(len(ok.CList)))
    if len(ok.CList) > 0 {
        for _, v := range ok.CList {
            binary.Write(buf, binary.LittleEndian, v)
        }
    }
    ok.Dump(buf)
    return nil
}

type GCSkillToTileOK5 struct {
    DummyRead

    ObjectID     ObjectID_t
    SkillType    SkillType_t
    X            Coord_t
    Y            Coord_t
    Range        Range_t
    Duration     Duration_t
    CreatureList []ObjectID_t
    Grade        uint8
}

func (ok *GCSkillToTileOK5) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_TILE_OK_5
}
func (ok *GCSkillToTileOK5) PacketSize() uint32 {
    return uint32(12 + 1 + len(ok.CreatureList)*4)
}

func (ok *GCSkillToTileOK5) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.X)
    binary.Write(buf, binary.LittleEndian, ok.Y)
    binary.Write(buf, binary.LittleEndian, ok.Range)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    binary.Write(buf, binary.LittleEndian, uint8(len(ok.CreatureList)))
    for _, v := range ok.CreatureList {
        binary.Write(buf, binary.LittleEndian, v)
    }
    return nil
}

type GCSkillToTileOK4 struct {
    DummyRead

    SkillType    SkillType_t
    X            Coord_t
    Y            Coord_t
    Range        uint8
    Duration     uint16
    CreatureList []uint32
    Grade        uint8
}

func (ok *GCSkillToTileOK4) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_TILE_OK_4
}
func (ok *GCSkillToTileOK4) PacketSize() uint32 {
    return uint32(8 + 1 + len(ok.CreatureList)*4)
}
func (ok *GCSkillToTileOK4) String() string {
    return "skill to tile ok 4"
}
func (ok *GCSkillToTileOK4) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.X)
    binary.Write(buf, binary.LittleEndian, ok.Y)
    binary.Write(buf, binary.LittleEndian, ok.Range)
    binary.Write(buf, binary.LittleEndian, ok.Duration)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    binary.Write(buf, binary.LittleEndian, uint8(len(ok.CreatureList)))
    for _, v := range ok.CreatureList {
        binary.Write(buf, binary.LittleEndian, v)
    }
    return nil
}

type GCSkillToTileOK3 struct {
    DummyRead

    ObjectID  ObjectID_t
    SkillType SkillType_t
    X         Coord_t
    Y         Coord_t
    Grade     uint8
}

func (ok *GCSkillToTileOK3) PacketID() PacketID {
    return PACKET_GC_SKILL_TO_TILE_OK_3
}
func (ok *GCSkillToTileOK3) PacketSize() uint32 {
    return 9
}
func (ok *GCSkillToTileOK3) String() string {
    return "skill to tile ok 3"
}
func (ok *GCSkillToTileOK3) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, ok.ObjectID)
    binary.Write(buf, binary.LittleEndian, ok.SkillType)
    binary.Write(buf, binary.LittleEndian, ok.X)
    binary.Write(buf, binary.LittleEndian, ok.Y)
    binary.Write(buf, binary.LittleEndian, ok.Grade)
    return nil
}

const (
    SYSTEM_MESSAGE_NORMAL = iota
    SYSTEM_MESSAGE_OPERATOR
    SYSTEM_MESSAGE_MASTER_LAIR
    SYSTEM_MESSAGE_COMBAT
    SYSTEM_MESSAGE_INFO
    SYSTEM_MESSAGE_HOLY_LAND
    SYSTEM_MESSAGE_RANGER_SAY
    SYSTEM_MESSAGE_MAX
)

type GCSystemMessagePacket struct {
    DummyRead

    Message string
    Color   uint32
    Type    uint8
}

func (msg *GCSystemMessagePacket) PacketID() PacketID {
    return PACKET_GC_SYSTEM_MESSAGE
}
func (msg *GCSystemMessagePacket) PacketSize() uint32 {
    return uint32(1 + len(msg.Message) + 4 + 1)
}
func (msg *GCSystemMessagePacket) String() string {
    return "system message"
}
func (msg *GCSystemMessagePacket) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, uint8(len(msg.Message)))
    io.WriteString(buf, msg.Message)
    binary.Write(buf, binary.LittleEndian, msg.Color)
    binary.Write(buf, binary.LittleEndian, msg.Type)
    return nil
}

const (
    PC_SLAYER PCType = iota
    PC_VAMPIRE
    PC_OUSTER
)

type GCSkillInfoPacket struct {
    DummyRead

    PCType          PCType
    PCSkillInfoList []data.PCSkillInfo
}

func (info *GCSkillInfoPacket) PacketSize() uint32 {
    var sz uint32
    sz = 2
    for i := 0; i < len(info.PCSkillInfoList); i++ {
        sz += info.PCSkillInfoList[i].Size()
    }
    return sz
}
func (info *GCSkillInfoPacket) PacketID() PacketID {
    return PACKET_GC_SKILL_INFO
}
func (info *GCSkillInfoPacket) String() string {
    return "skill info"
}
func (info *GCSkillInfoPacket) Write(buf io.Writer, code uint8) error {

    binary.Write(buf, binary.LittleEndian, uint8(info.PCType))
    binary.Write(buf, binary.LittleEndian, uint8(len(info.PCSkillInfoList)))
    for _, v := range info.PCSkillInfoList {
        v.Write(buf)
    }
    return nil
}

type GCMoveOK struct {
    DummyRead

    X   uint8
    Y   uint8
    Dir uint8
}

type GCDisconnect struct {
    DummyRead

    Message string
}

func (disconn GCDisconnect) PacketID() PacketID {
    return PACKET_GC_DISCONNECT
}
func (disconn GCDisconnect) PacketSize() uint32 {
    return 3
}

func (disconn GCDisconnect) Write(buf io.Writer, code uint8) error {

    sz := uint8(len(disconn.Message))

    binary.Write(buf, binary.LittleEndian, sz)
    io.WriteString(buf, disconn.Message)
    return nil
}
