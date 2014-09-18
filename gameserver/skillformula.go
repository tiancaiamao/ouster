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

func HitRoll(pAttacker pDefender, CreatureInterface, bonus int) bool {
    if pDefender.CreatureInstance().isFlag(EFFECT_CLASS_NO_DAMAGE) {
        return false
    }

    tohit := 0
    defense := 0
    //TODO
    // timeband = pZone->getTimeband();
    timeband := 0

    switch pAttacker.(type) {
    case *Slayer:
        tohit = pAttacker.(*Slayer).ToHit[ATTR_CURRENT]
    case *Ouster:
        tohit = pAttacker.(*Ouster).ToHit[ATTR_CURRENT]
    case *Vampire:
        tohit = pAttacker.(*Vampire).ToHit[ATTR_CURRENT]
        tohit = getPrecentValue(tohit, VampireTimebandFactor[timeband])
    case *Monster:
        tohit = pAttacker.(*Monster).ToHit[ATTR_CURRENT]
        tohit = getPrecentValue(tohit, VampireTimebandFactor[timeband])
    }

    switch pDefender.(type) {
    case *Slayer:
        defense = pAttacker.(*Slayer).Defense[ATTR_CURRENT]
    case *Ouster:
        defense = pAttacker.(*Ouster).Defense[ATTR_CURRENT]
    case *Vampire:
        defense = pAttacker.(*Vampire).Defense[ATTR_CURRENT]
        defense = getPrecentValue(defense, VampireTimebandFactor[timeband])
    case *Monster:
        defense = pAttacker.(*Monster).Defense[ATTR_CURRENT]
        defense = getPrecentValue(defense, VampireTimebandFactor[timeband])
    }

    randValue = rand.Intn(100)
    result = 0

    if tohit >= defense {
        Result = min(90, int(((tohit-defense)/1.5)+60)+bonus)
    } else {
        if isMonster {
            Result = max(10, (int)(60-((defense-tohit)/1.5)+bonus))
        } else {
            Result = max(20, (int)(60-((defense-tohit)/1.5)+bonus))
        }
    }

    if randValue <= result {
        return true
    }

    return false
}
