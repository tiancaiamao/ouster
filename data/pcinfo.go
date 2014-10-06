package data

import (
    "encoding/binary"
    "errors"
    "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

const (
    ATTR_CURRENT = iota
    ATTR_MAX
    ATTR_BASIC
)

type PCInfo interface {
    Write(writer io.Writer) error
    Read(reader io.Reader) error
    Size() uint32
}

type PCOusterInfo struct {
    ObjectID ObjectID_t
    Name     string
    Level    Level_t
    Sex      Sex_t

    HairColor         Color_t
    MasterEffectColor uint8

    Alignment Alignment_t
    STR       [3]Attr_t
    DEX       [3]Attr_t
    INI       [3]Attr_t

    HP  [2]HP_t
    MP  [2]MP_t

    Rank    Rank_t
    RankExp RankExp_t
    Exp     Exp_t

    Fame       Fame_t
    Gold       Gold_t
    Sight      Sight_t
    Bonus      Bonus_t
    SkillBonus SkillBonus_t

    SilverDamage       Silver_t
    Competence         uint8
    GuildID            GuildID_t
    GuildName          string
    GuildMemberRank    GuildMemberRank_t
    UnionID            uint32
    AdvancementLevel   Level_t
    AdvancementGoalExp Exp_t

    ZoneID ZoneID_t
    ZoneX  ZoneCoord_t
    ZoneY  ZoneCoord_t
}

func (info *PCOusterInfo) Size() uint32 {
    sz := 4 + 1 + uint32(len(info.Name)) + 2
    sz = sz + 2 + 1 + 4 + 9*2
    sz = sz + 4*2 + 1 + 4 + 4
    sz = sz + 4 + 4 + 1 + 2 + 2
    sz = sz + 2 + 1 + 2 + 1 + uint32(len(info.GuildName)) + 1 + 4 + 1 + 4
    return sz
}

func (info *PCOusterInfo) Read(reader io.Reader) error {
    binary.Read(reader, binary.LittleEndian, &info.ObjectID)
    var szName uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &szName)
    reader.Read(buf[:szName])
    info.Name = string(buf[:szName])
    binary.Read(reader, binary.LittleEndian, &info.Level)
    binary.Read(reader, binary.LittleEndian, &info.Sex)

    binary.Read(reader, binary.LittleEndian, &info.HairColor)
    binary.Read(reader, binary.LittleEndian, &info.MasterEffectColor)

    binary.Read(reader, binary.LittleEndian, &info.Alignment)

    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_BASIC])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_BASIC])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_BASIC])

    binary.Read(reader, binary.LittleEndian, &info.HP[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.HP[ATTR_MAX])

    binary.Read(reader, binary.LittleEndian, &info.MP[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.MP[ATTR_MAX])

    binary.Read(reader, binary.LittleEndian, &info.Rank)
    binary.Read(reader, binary.LittleEndian, &info.RankExp)

    binary.Read(reader, binary.LittleEndian, &info.Exp)
    binary.Read(reader, binary.LittleEndian, &info.Fame)
    binary.Read(reader, binary.LittleEndian, &info.Gold)
    binary.Read(reader, binary.LittleEndian, &info.Sight)

    binary.Read(reader, binary.LittleEndian, &info.Bonus)
    binary.Read(reader, binary.LittleEndian, &info.SkillBonus)

    binary.Read(reader, binary.LittleEndian, &info.SilverDamage)
    binary.Read(reader, binary.LittleEndian, &info.Competence)
    binary.Read(reader, binary.LittleEndian, &info.GuildID)

    var szGuildName uint8
    binary.Read(reader, binary.LittleEndian, &szGuildName)
    reader.Read(buf[:szGuildName])
    info.GuildName = string(buf[:szGuildName])

    binary.Read(reader, binary.LittleEndian, &info.GuildMemberRank)
    binary.Read(reader, binary.LittleEndian, &info.UnionID)
    binary.Read(reader, binary.LittleEndian, &info.AdvancementLevel)
    binary.Read(reader, binary.LittleEndian, &info.AdvancementGoalExp)

    return nil
}

func (info *PCOusterInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.ObjectID)
    binary.Write(writer, binary.LittleEndian, uint8(len(info.Name)))
    io.WriteString(writer, info.Name)
    binary.Write(writer, binary.LittleEndian, info.Level)
    binary.Write(writer, binary.LittleEndian, info.Sex)

    binary.Write(writer, binary.LittleEndian, info.HairColor)
    binary.Write(writer, binary.LittleEndian, info.MasterEffectColor)

    binary.Write(writer, binary.LittleEndian, info.Alignment)

    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_BASIC])

    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_MAX])

    binary.Write(writer, binary.LittleEndian, info.MP[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.MP[ATTR_MAX])

    binary.Write(writer, binary.LittleEndian, info.Rank)
    binary.Write(writer, binary.LittleEndian, info.RankExp)

    binary.Write(writer, binary.LittleEndian, info.Exp)
    binary.Write(writer, binary.LittleEndian, info.Fame)
    binary.Write(writer, binary.LittleEndian, info.Gold)
    binary.Write(writer, binary.LittleEndian, info.Sight)

    binary.Write(writer, binary.LittleEndian, info.Bonus)
    binary.Write(writer, binary.LittleEndian, info.SkillBonus)

    binary.Write(writer, binary.LittleEndian, info.SilverDamage)
    binary.Write(writer, binary.LittleEndian, info.Competence)
    binary.Write(writer, binary.LittleEndian, info.GuildID)

    binary.Write(writer, binary.LittleEndian, uint8(len(info.GuildName)))
    io.WriteString(writer, info.GuildName)

    binary.Write(writer, binary.LittleEndian, info.GuildMemberRank)
    binary.Write(writer, binary.LittleEndian, info.UnionID)
    binary.Write(writer, binary.LittleEndian, info.AdvancementLevel)
    binary.Write(writer, binary.LittleEndian, info.AdvancementGoalExp)

    return nil
}

type PCVampireInfo struct {
    ObjectID ObjectID_t
    Name     string
    Level    Level_t
    Sex      Sex_t

    BatColor          Color_t
    SkinColor         Color_t
    MasterEffectColor uint8

    Alignment Alignment_t
    STR       [3]Attr_t
    DEX       [3]Attr_t
    INI       [3]Attr_t

    HP  [2]HP_t

    Rank    Rank_t
    RankExp RankExp_t

    Exp          Exp_t
    Fame         Fame_t
    Gold         Gold_t
    Sight        Sight_t
    Bonus        Bonus_t
    HotKey       [8]uint16
    SilverDamage Silver_t

    Competence uint8
    GuildID    GuildID_t
    GuildName  string

    GuildMemberRank GuildMemberRank_t
    UnionID         uint32

    AdvancementLevel   Level_t
    AdvancementGoalExp Exp_t

    ZoneID ZoneID_t
    ZoneX  Coord_t
    ZoneY  Coord_t
}

// TODO
func (info *PCVampireInfo) Size() uint32 {
    return 0
}

func (info *PCVampireInfo) Read(reader io.Reader) error {
    return errors.New("read PCVampireInfo not implement yet!!!")
}
func (info *PCVampireInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.ObjectID)
    binary.Write(writer, binary.LittleEndian, uint8(len(info.Name)))
    io.WriteString(writer, info.Name)
    binary.Write(writer, binary.LittleEndian, info.Level)
    binary.Write(writer, binary.LittleEndian, info.Sex)

    binary.Write(writer, binary.LittleEndian, info.BatColor)
    binary.Write(writer, binary.LittleEndian, info.SkinColor)
    binary.Write(writer, binary.LittleEndian, info.MasterEffectColor)

    binary.Write(writer, binary.LittleEndian, info.Alignment)

    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_BASIC])

    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_MAX])

    binary.Write(writer, binary.LittleEndian, info.Rank)
    binary.Write(writer, binary.LittleEndian, info.RankExp)

    binary.Write(writer, binary.LittleEndian, info.Exp)
    binary.Write(writer, binary.LittleEndian, info.Gold)

    binary.Write(writer, binary.LittleEndian, info.Fame)
    binary.Write(writer, binary.LittleEndian, info.Sight)
    binary.Write(writer, binary.LittleEndian, info.Bonus)

    for i := 0; i < 8; i++ {
        binary.Write(writer, binary.LittleEndian, info.HotKey[i])
    }

    binary.Write(writer, binary.LittleEndian, info.SilverDamage)
    binary.Write(writer, binary.LittleEndian, info.Competence)
    binary.Write(writer, binary.LittleEndian, info.GuildID)

    binary.Write(writer, binary.LittleEndian, uint8(len(info.GuildName)))
    io.WriteString(writer, info.GuildName)

    binary.Write(writer, binary.LittleEndian, info.GuildMemberRank)
    binary.Write(writer, binary.LittleEndian, info.UnionID)
    binary.Write(writer, binary.LittleEndian, info.AdvancementLevel)
    binary.Write(writer, binary.LittleEndian, info.AdvancementGoalExp)

    return nil
}

type PCSlayerInfo struct {
    ObjectID          ObjectID_t
    Name              string
    Sex               Sex_t
    HairStyle         HairStyle
    HairColor         Color_t
    SkinColor         Color_t
    MasterEffectColor byte
    Alignment         Alignment_t
    STR               [3]Attr_t
    DEX               [3]Attr_t
    INI               [3]Attr_t
    STRExp            Exp_t
    DEXExp            Exp_t
    INIExp            Exp_t
    Rank              Rank_t
    RankExp           RankExp_t
    HP                [2]HP_t
    MP                [2]MP_t
    Fame              Fame_t
    Gold              Gold_t

    DomainLevels [6]SkillLevel_t
    DomainExps   [6]SkillExp_t

    Sight           Sight_t
    HotKey          [4]SkillType_t
    Competence      byte
    GuildID         GuildID_t
    GuildName       string
    GuildMemberRank GuildMemberRank_t

    UnionID            uint32
    AdvancementLevel   Level_t
    AdvancementGoalExp Exp_t
    AttrBonus          Bonus_t

    ZoneID ZoneID_t
    ZoneX  Coord_t
    ZoneY  Coord_t
}

// TODO
func (info *PCSlayerInfo) Size() uint32 {
    return 0
}

func (info *PCSlayerInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.ObjectID)
    binary.Write(writer, binary.LittleEndian, uint8(len(info.Name)))
    io.WriteString(writer, info.Name)
    binary.Write(writer, binary.LittleEndian, info.Sex)
    binary.Write(writer, binary.LittleEndian, info.HairStyle)
    binary.Write(writer, binary.LittleEndian, info.HairColor)
    binary.Write(writer, binary.LittleEndian, info.SkinColor)
    binary.Write(writer, binary.LittleEndian, info.MasterEffectColor)

    binary.Write(writer, binary.LittleEndian, info.Alignment)

    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.STR[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.DEX[ATTR_BASIC])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.INI[ATTR_BASIC])

    binary.Write(writer, binary.LittleEndian, info.Rank)
    binary.Write(writer, binary.LittleEndian, info.RankExp)

    binary.Write(writer, binary.LittleEndian, info.STRExp)
    binary.Write(writer, binary.LittleEndian, info.DEXExp)
    binary.Write(writer, binary.LittleEndian, info.INIExp)

    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.HP[ATTR_MAX])
    binary.Write(writer, binary.LittleEndian, info.MP[ATTR_CURRENT])
    binary.Write(writer, binary.LittleEndian, info.MP[ATTR_MAX])

    binary.Write(writer, binary.LittleEndian, info.Fame)
    binary.Write(writer, binary.LittleEndian, info.Gold)

    for i := 0; i < SKILL_DOMAIN_VAMPIRE; i++ {
        binary.Write(writer, binary.LittleEndian, info.DomainLevels[i])
        binary.Write(writer, binary.LittleEndian, info.DomainExps[i])
    }

    binary.Write(writer, binary.LittleEndian, info.Sight)

    for i := 0; i < 4; i++ {
        binary.Write(writer, binary.LittleEndian, info.HotKey[i])
    }

    binary.Write(writer, binary.LittleEndian, info.Competence)
    binary.Write(writer, binary.LittleEndian, info.GuildID)

    binary.Write(writer, binary.LittleEndian, uint8(len(info.GuildName)))
    io.WriteString(writer, info.GuildName)

    binary.Write(writer, binary.LittleEndian, info.GuildMemberRank)
    binary.Write(writer, binary.LittleEndian, info.UnionID)
    binary.Write(writer, binary.LittleEndian, info.AdvancementLevel)
    binary.Write(writer, binary.LittleEndian, info.AdvancementGoalExp)
    binary.Write(writer, binary.LittleEndian, info.AttrBonus)
    return nil
}

func (info *PCSlayerInfo) Read(reader io.Reader) error {
    binary.Read(reader, binary.LittleEndian, &info.ObjectID)
    var szName uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &szName)
    log.Debugln("szName=", szName)
    _, err := reader.Read(buf[:szName])
    if err != nil {
        return err
    }
    info.Name = string(buf[:szName])
    log.Debugln("sdkjfasdlgasjdfasdf", info.Name, szName)
    binary.Read(reader, binary.LittleEndian, &info.Sex)
    binary.Read(reader, binary.LittleEndian, &info.HairStyle)
    binary.Read(reader, binary.LittleEndian, &info.HairColor)
    binary.Read(reader, binary.LittleEndian, &info.SkinColor)
    binary.Read(reader, binary.LittleEndian, &info.MasterEffectColor)
    binary.Read(reader, binary.LittleEndian, &info.Alignment)
    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.STR[ATTR_BASIC])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.DEX[ATTR_BASIC])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.INI[ATTR_BASIC])

    binary.Read(reader, binary.LittleEndian, &info.Rank)
    binary.Read(reader, binary.LittleEndian, &info.RankExp)

    binary.Read(reader, binary.LittleEndian, &info.STRExp)
    binary.Read(reader, binary.LittleEndian, &info.DEXExp)
    binary.Read(reader, binary.LittleEndian, &info.INIExp)

    binary.Read(reader, binary.LittleEndian, &info.HP[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.HP[ATTR_MAX])
    binary.Read(reader, binary.LittleEndian, &info.MP[ATTR_CURRENT])
    binary.Read(reader, binary.LittleEndian, &info.MP[ATTR_MAX])

    binary.Read(reader, binary.LittleEndian, &info.Fame)
    binary.Read(reader, binary.LittleEndian, &info.Gold)

    for i := 0; i < SKILL_DOMAIN_VAMPIRE; i++ {
        binary.Read(reader, binary.LittleEndian, &info.DomainLevels[i])
        binary.Read(reader, binary.LittleEndian, &info.DomainExps[i])
    }

    binary.Read(reader, binary.LittleEndian, &info.Sight)

    for i := 0; i < 4; i++ {
        binary.Read(reader, binary.LittleEndian, &info.HotKey[i])
    }

    binary.Read(reader, binary.LittleEndian, &info.Competence)
    binary.Read(reader, binary.LittleEndian, &info.GuildID)

    var szGuildName uint8
    _, err = reader.Read(buf[:szGuildName])
    if err != nil {
        return err
    }
    info.GuildName = string(buf[:szGuildName])

    binary.Read(reader, binary.LittleEndian, &info.GuildMemberRank)
    binary.Read(reader, binary.LittleEndian, &info.UnionID)
    binary.Read(reader, binary.LittleEndian, &info.AdvancementLevel)
    binary.Read(reader, binary.LittleEndian, &info.AdvancementGoalExp)
    binary.Read(reader, binary.LittleEndian, &info.AttrBonus)
    return nil
}
