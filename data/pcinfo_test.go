package data

import (
    "bytes"
    . "github.com/tiancaiamao/ouster/util"
    "testing"
)

func TestPCOusterInfoWrite(t *testing.T) {
    info := PCOusterInfo{
        ObjectID: 1,
        Name:     "test",
        Level:    150,

        HairColor: 101,

        Alignment: 7500,
        STR:       [3]Attr_t{10, 10, 10},
        DEX:       [3]Attr_t{25, 25, 25},
        INI:       [3]Attr_t{10, 10, 10},

        HP: [2]HP_t{315, 315},
        MP: [2]MP_t{186, 111},

        Rank:    50,
        RankExp: 10700,
        Exp:     125,

        Fame:       500,
        Gold:       68,
        Sight:      13,
        Bonus:      9999,
        SkillBonus: 9999,

        Competence: 1,
        GuildID:    66,

        GuildMemberRank:  4,
        AdvancementLevel: 100,
    }

    buf := &bytes.Buffer{}
    info.Write(buf)

    right := []byte{1, 0, 0, 0, 4, 116, 101, 115, 116, 150, 0, 101, 0, 0, 76, 29, 0, 0, 10, 0, 10, 0, 10, 0, 25, 0, 25, 0, 25, 0, 10, 0, 10, 0, 10, 0, 59, 1, 59, 1, 186, 0, 111, 0, 50, 204, 41, 0, 0, 125, 0, 0, 0, 244, 1, 0, 0, 68, 0, 0, 0, 13, 15, 39, 15, 39, 0, 0, 1, 66, 0, 0, 4, 0, 0, 0, 0, 100, 0, 0, 0, 0}
    if !bytes.Equal(buf.Bytes(), right) {
        t.Failed()
    }

    slayerInfo := &PCSlayerInfo{
        ObjectID:          0x76dd,
        Name:              "test",
        Sex:               1,
        HairStyle:         1,
        HairColor:         0x17d,
        SkinColor:         0x1e0,
        MasterEffectColor: 0x0,
        Alignment:         7500,
        STR:               [3]Attr_t{0xd, 0xd, 0xd},
        DEX:               [3]Attr_t{0xe, 0xe, 0xe},
        INI:               [3]Attr_t{0x14, 0x14, 0x14},
        STRExp:            0x4f8,
        DEXExp:            0x5dd,
        INIExp:            0xfd3,
        Rank:              0x1,
        RankExp:           0x29cc,
        HP:                [2]HP_t{0x36, 0x36},
        MP:                [2]MP_t{0x32, 0x32},
        Fame:              0x0,
        Gold:              0x1f4,
        DomainLevels:      [6]SkillLevel_t{0x0, 0x0, 0x0, 0x7, 0x0, 0x0},
        DomainExps:        [6]SkillExp_t{0x32, 0x32, 0x28, 0x366, 0x1e, 0xf4240},
        Sight:             0xd,
        HotKey:            [4]SkillType_t{0x4972, 0x39, 0x0, 0xd748},
        Competence:        0x1,
        GuildID:           0x63, GuildName: "",
        GuildMemberRank:    0x0,
        UnionID:            0x4,
        AdvancementLevel:   0x0,
        AdvancementGoalExp: 0x0,
        AttrBonus:          0x0,
        ZoneID:             0x0,
        ZoneX:              0x0,
        ZoneY:              0x0,
    }
    buf.Reset()
    slayerInfo.Write(buf)
    right = []byte{221, 118, 0, 0, 4, 116, 101, 115, 116, 1, 1, 125, 1, 224, 1, 0, 76, 29, 0, 0, 13, 0, 13, 0, 13, 0, 14, 0, 14, 0, 14, 0, 20, 0, 20, 0, 20, 0, 1, 204, 41, 0, 0, 248, 4, 0, 0, 221, 5, 0, 0, 211, 15, 0, 0, 54, 0, 54, 0, 50, 0, 50, 0, 0, 0, 0, 0, 244, 1, 0, 0, 0, 50, 0, 0, 0, 0, 50, 0, 0, 0, 0, 40, 0, 0, 0, 7, 102, 3, 0, 0, 0, 30, 0, 0, 0, 0, 64, 66, 15, 0, 13, 114, 73, 57, 0, 0, 0, 72, 215, 1, 99, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    if !bytes.Equal(buf.Bytes(), right) {
        t.Failed()
    }
}
