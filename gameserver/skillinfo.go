package main

type SkillInfo struct {
    Type        uint
    Name        string
    MinDamage   uint
    MaxDamage   uint
    MinDelay    uint
    MaxDelay    uint
    MinCastTime uint
    MaxCastTime uint

    MinDuration int
    MaxDuration int

    ConsumeMP uint

    MaxRange uint
    MinRange uint

    // 0x01 : burrowing
    // 0x02 : walking
    // 0x04 : flying
    Target uint

    SubSkill uint
    Point    uint
    Domain   byte

    MagicDomain   int
    ElementDomain int

    SkillPoint   int
    LevelUpPoint int

    RequireSkills  []SkillType_t
    RequiredSkills []SkillType_t

    CanDelete byte

    RequireFire  Elemental_t
    RequireWater Elemental_t
    RequireEarth Elemental_t
    RequireWind  Elemental_t
    RequireSum   Elemental_t

    RequireWristletElemental ElementalType
    RequireStone1Elemental   ElementalType
    RequireStone2Elemental   ElementalType
    RequireStone3Elemental   ElementalType
    RequireStone4Elemental   ElementalType
}
