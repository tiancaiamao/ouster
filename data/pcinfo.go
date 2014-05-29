package data

import (
	"encoding/binary"
	"io"
)

const (
	ATTR_CURRENT = iota
	ATTR_MAX
	ATTR_BASIC
)

type PCInfo interface {
	Dump(writer io.Writer)
}

type PCOusterInfo struct {
	PCType   byte
	ObjectID uint32
	Name     string
	Level    uint8
	Sex      uint8

	HairColor         uint16
	MasterEffectColor uint8

	Alignment uint32
	STR       [3]uint16
	DEX       [3]uint16
	INT       [3]uint16

	HP [2]uint16
	MP [2]uint16

	Rank    uint8
	RankExp uint32
	Exp     uint32

	Fame       uint32
	Gold       uint32
	Sight      uint8
	Bonus      uint16
	SkillBonus uint16

	SilverDamage       uint16
	Competence         uint8
	GuildID            uint16
	GuildName          string
	GuildMemberRank    uint8
	UnionID            uint32
	AdvancementLevel   uint8
	AdvancementGoalExp uint32

	ZoneID uint16
	ZoneX  uint8
	ZoneY  uint8
}

func (info *PCOusterInfo) Dump(writer io.Writer) {
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
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_CURRENT])
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_MAX])
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_BASIC])

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

	return
}

type PCVampireInfo struct {
	PCType   byte
	ObjectID uint32
	Name     string
	Level    uint8
	Sex      uint8

	BatColor          uint16
	SkinColor         uint16
	MasterEffectColor uint8

	Alignment uint32
	STR       [3]uint16
	DEX       [3]uint16
	INT       [3]uint16

	HP [2]uint16

	Rank    uint8
	RankExp uint32

	Exp          uint32
	Fame         uint32
	Gold         uint32
	Sight        uint8
	Bonus        uint16
	HotKey       [8]uint16
	SilverDamage uint16

	Competence uint8
	GuildID    uint16
	GuildName  string

	GuildMemberRank uint8
	UnionID         uint32

	AdvancementLevel   uint8
	AdvancementGoalExp uint32

	ZoneID uint16
	ZoneX  uint8
	ZoneY  uint8
}

func (info *PCVampireInfo) Dump(writer io.Writer) {
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
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_CURRENT])
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_MAX])
	binary.Write(writer, binary.LittleEndian, info.INT[ATTR_BASIC])

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

	return
}
