package main

func (ignore MeteorStrikeHandler) ComputeOutput(c1 *Creature, c2 *Creature) SkillOutput {
    return SkillOutput{
    // Damage: int(float32(c1.Level)*0.8) + int(c1.STR[ATTR_CURRENT]+c1.DEX[ATTR_CURRENT])/6,
    }
}

func (ignore ParalyzeHandler) ComputeOutput(c1 *Creature, c2 *Creature) SkillOutput {
    return SkillOutput{
    // Duration: int((3 + c1.INT[ATTR_CURRENT]/15) * 10),
    }
}
