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
    Interval    time.Duration
    CastingTime time.Duration
    RunTime     time.Time
    Enable      bool
}

type VampireSkillSlot struct {
    Name        string
    SkillType   SkillType_t
    Interval    time.Duration
    CastingTime time.Duration
    RunTime     time.Time
}

type OusterSkillSlot struct {
    Name        string
    SkillType   SkillType_t
    ExpLevel    ExpLevel_t
    Interval    time.Duration
    CastingTime time.Duration
    RunTime     time.Time
}

type SkillSlot struct {
    SkillType uint16
    ExpLevel  uint16

    LastUse  time.Time
    Cooling  uint16
    Duration uint16

    Interval    uint32
    CastingTime uint32
}

func verifyRunTime(skillSlot SkillSlot) bool {
    return true
}
