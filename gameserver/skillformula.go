package main

func (ignore MeteorStrike) ComputeOutput(c1 *SkillInput, c2 *SkillOutput) {
    return
    // Damage: int(float32(c1.Level)*0.8) +
    // int(c1.STR[ATTR_CURRENT]+c1.DEX[ATTR_CURRENT])/6,
}

func (ignore Paralyze) ComputeOutput(c1 *SkillInput, c2 *SkillOutput) {
    return
    // Duration: int((3 + c1.INT[ATTR_CURRENT]/15) * 10),

}
