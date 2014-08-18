package main

import (
    "encoding/json"
    // "fmt"
    "github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
    "os"
)

const (
    SKILL_RAPID_GLIDING uint16 = 203
    SKILL_METEOR_STRIKE uint16 = 180
    SKILL_INVISIBILITY  uint16 = 100
    SKILL_PARALYZE      uint16 = 89
    SKILL_BLOOD_SPEAR   uint16 = 97

    SKILL_ABSORB_SOUL       uint16 = 246
    SKILL_SUMMON_SYLPH      uint16 = 247
    SKILL_SHARP_HAIL        uint16 = 348 // 尖锐冰雹
    SKILL_FLOURISH          uint16 = 219 // 活跃攻击
    SKILL_DESTRUCTION_SPEAR uint16 = 298 //致命爆发
    SKILL_SHARP_CHAKRAM     uint16 = 295 // 税利之轮
    SKILL_EVADE             uint16 = 220 // 回避术

    SKILL_FIRE_OF_SOUL_STONE uint16 = 227
    SKILL_ICE_OF_SOUL_STONE  uint16 = 228
    SKILL_SAND_OF_SOUL_STONE uint16 = 229

    SKILL_TELEPORT       uint16 = 280 // 瞬间移动
    SKILL_DUCKING_WALLOP uint16 = 302 // 光速冲击
    SKILL_DISTANCE_BLITZ uint16 = 304 // 雷神斩

)

func main() {
    vampire := data.PCVampireInfo{
        PCType:           'V',
        Name:             "vampire",
        Level:            150,
        SkinColor:        420,
        Alignment:        7500,
        STR:              [3]uint16{20, 20, 20},
        DEX:              [3]uint16{20, 20, 20},
        INT:              [3]uint16{20, 20, 20},
        HP:               [2]uint16{472, 472},
        Rank:             50,
        RankExp:          10700,
        Exp:              125,
        Fame:             282,
        Sight:            13,
        Bonus:            9999,
        Competence:       1,
        GuildMemberRank:  4,
        AdvancementLevel: 100,
        ZoneID:           23,
        ZoneX:            145,
        ZoneY:            237,
    }

    skill := packet.VampireSkillInfo{
        LearnNewSkill: false,
        SubVampireSkillInfoList: []packet.SubVampireSkillInfo{
            packet.SubVampireSkillInfo{
                SkillType:   SKILL_RAPID_GLIDING,
                Interval:    50,
                CastingTime: 31,
            },
            packet.SubVampireSkillInfo{
                SkillType:   SKILL_METEOR_STRIKE,
                Interval:    10,
                CastingTime: 4160749567,
            },
            packet.SubVampireSkillInfo{
                SkillType:   SKILL_INVISIBILITY,
                Interval:    30,
                CastingTime: 11,
            },
            packet.SubVampireSkillInfo{
                SkillType:   SKILL_PARALYZE,
                Interval:    60,
                CastingTime: 41,
            },
            packet.SubVampireSkillInfo{
                SkillType:   SKILL_BLOOD_SPEAR,
                Interval:    60,
                CastingTime: 41,
            },
        },
    }

    // ouster := data.PCOusterInfo{
    // 		PCType:           'O',
    // 		Name:             "ouster",
    // 		Level:            150,
    // 		HairColor:        101,
    // 		Alignment:        7500,
    // 		STR:              [3]uint16{10, 10, 10},
    // 		DEX:              [3]uint16{25, 25, 25},
    // 		INT:              [3]uint16{10, 10, 10},
    // 		HP:               [2]uint16{315, 315},
    // 		MP:               [2]uint16{186, 111},
    // 		Rank:             50,
    // 		RankExp:          10700,
    // 		Exp:              125,
    // 		Fame:             500,
    // 		Gold:             68,
    // 		Sight:            13,
    // 		Bonus:            9999,
    // 		SkillBonus:       9999,
    // 		Competence:       1,
    // 		GuildID:          66,
    // 		GuildMemberRank:  4,
    // 		AdvancementLevel: 100,
    //
    // 		ZoneID: 62,
    // 		ZoneX:  62,
    // 		ZoneY:  58,
    // 	}
    //
    // 	skill :=  OusterSkillInfo{
    // 		LearnNewSkill: false,
    // 		SubOusterSkillInfoList: [] SubOusterSkillInfo{
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_FLOURISH,
    // 				ExpLevel:    1,
    // 				Interval:    10,
    // 				CastingTime: 6,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType: SKILL_ABSORB_SOUL,
    // 				ExpLevel:  1,
    // 				Interval:  5,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType: SKILL_SUMMON_SYLPH,
    // 				ExpLevel:  1,
    // 				Interval:  5,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_SHARP_HAIL,
    // 				ExpLevel:    1,
    // 				Interval:    112,
    // 				CastingTime: 107,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_DISTANCE_BLITZ,
    // 				ExpLevel:    1,
    // 				Interval:    70,
    // 				CastingTime: 65,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_DUCKING_WALLOP,
    // 				ExpLevel:    1,
    // 				Interval:    100,
    // 				CastingTime: 95,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_DESTRUCTION_SPEAR,
    // 				ExpLevel:    1,
    // 				Interval:    50,
    // 				CastingTime: 45,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_SHARP_CHAKRAM,
    // 				ExpLevel:    1,
    // 				Interval:    600,
    // 				CastingTime: 600,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_TELEPORT,
    // 				ExpLevel:    1,
    // 				Interval:    50,
    // 				CastingTime: 45,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_SUMMON_SYLPH,
    // 				ExpLevel:    1,
    // 				Interval:    5,
    // 				CastingTime: 0,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType: SKILL_ICE_OF_SOUL_STONE,
    // 				ExpLevel:  1,
    // 				Interval:  0,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType: SKILL_FIRE_OF_SOUL_STONE,
    // 				ExpLevel:  1,
    // 				Interval:  0,
    // 			},
    // 			 SubOusterSkillInfo{
    // 				SkillType:   SKILL_EVADE,
    // 				ExpLevel:    1,
    // 				Interval:    600,
    // 				CastingTime: 600,
    // 			},
    // 		},
    // 	}

    f, err := os.Create("/Users/genius/.ouster/player/vampire")
    if err != nil {
        panic(err)
    }

    encoder := json.NewEncoder(f)
    err = encoder.Encode(vampire)
    if err != nil {
        panic(err)
    }

    err = encoder.Encode(skill)
    if err != nil {
        panic(err)
    }
    return
}
