package main

import (
    . "github.com/tiancaiamao/ouster/util"
    "time"
)

type SlayerSkillSlot struct {
    Name        string
    SkillType   SkillType_t
    Exp         Exp_t
    ExpLevel    ExpLevel_t
    Interval    Turn_t
    CastingTime Turn_t
    RunTime     time.Time
    Enable      bool
}

type VampireSkillSlot struct {
    Name        string
    SkillType   SkillType_t
    Interval    Turn_t
    CastingTime Turn_t
    RunTime     time.Time
}

type OusterSkillSlot struct {
    Name        string
    SkillType   SkillType_t
    ExpLevel    ExpLevel_t
    Interval    Turn_t
    CastingTime Turn_t
    RunTime     time.Time
}

func verifyRuntime(skillSlot SkillSlot) bool {
    return true
}
