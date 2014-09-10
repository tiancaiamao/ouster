package packet

import (
    "bytes"
    "github.com/tiancaiamao/ouster/data"
    . "github.com/tiancaiamao/ouster/util"
    "testing"
)

func TestPCOusterInfo(t *testing.T) {
    info := data.PCOusterInfo{
        ObjectID:           12202,
        Name:               "ouster",
        Level:              150,
        Sex:                FEMALE,
        HairColor:          101,
        MasterEffectColor:  0,
        Alignment:          7500,
        STR:                [3]Attr_t{10, 10, 10},
        DEX:                [3]Attr_t{25, 25, 25},
        INI:                [3]Attr_t{10, 10, 10},
        HP:                 [2]HP_t{315, 315},
        MP:                 [2]MP_t{186, 111},
        Rank:               50,
        RankExp:            10700,
        Exp:                125,
        Fame:               500,
        Gold:               92,
        Sight:              13,
        Bonus:              9999,
        SkillBonus:         9994,
        SilverDamage:       0,
        Competence:         1,
        GuildID:            66,
        GuildName:          "",
        GuildMemberRank:    4,
        UnionID:            0,
        AdvancementLevel:   100,
        AdvancementGoalExp: 0,
    }
    buf := &bytes.Buffer{}
    info.Dump(buf)

    get := buf.Bytes()
    want := []byte{170, 47, 0, 0, 6, 'o', 'u', 's', 't', 'e', 'r', 150, 0, 101, 0, 0, 76, 29, 0, 0, 10, 0, 10, 0, 10, 0, 25, 0, 25, 0, 25, 0, 10, 0, 10, 0, 10, 0, 59, 1, 59, 1, 186, 0, 111, 0, 50, 204, 41, 0, 0, 125, 0, 0, 0, 244, 1, 0, 0, 92, 0, 0, 0, 13, 15, 39, 10, 39, 0, 0, 1, 66, 0, 0, 4, 0, 0, 0, 0, 100, 0, 0, 0, 0}
    if !bytes.Equal(want, get) {
        t.Errorf("want: %v\n, get:%v\n", want, get)
    }
}

func TestGCUpdateInfo(t *testing.T) {
}
