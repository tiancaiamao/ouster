package main

import (
    . "github.com/tiancaiamao/ouster/util"
    "math/rand"
)

func computeFinalDamage(
    minDamage Damage_t,
    maxDamage Damage_t,
    realDamage Damage_t,
    protection Protection_t,
    bCritical bool) Damage_t {
    // 致命一击无视防御
    if bCritical {
        return realDamage
    }

    finalDamage := realDamage - (realDamage*(Damage_t(protection)/8))/100

    return finalDamage

    // avgDamage := (minDamage + maxDamage) / 2
    //
    //     	int      DamageRatio = 100;
    //
    //     	if (Protection < avgDamage)
    //     	{
    //     		DamageRatio = 100;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 150))
    //     	{
    //     		DamageRatio = 90;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 200))
    //     	{
    //     		DamageRatio = 80;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 250))
    //     	{
    //     		DamageRatio = 70;
    //     	}
    //     	else if (Protection < getPercentValue(avgDamage, 300))
    //     	{
    //     		DamageRatio = 60;
    //     	}
    //     	else
    //     	{
    //     		DamageRatio = 50;
    //     	}

    // return max(1, getPercentValue(realDamage, DamageRatio))
}

func (ignore MeteorStrike) ComputeOutput(c1 *SkillInput, c2 *SkillOutput) {
    return
    // Damage: int(float32(c1.Level)*0.8) +
    // int(c1.STR[ATTR_CURRENT]+c1.DEX[ATTR_CURRENT])/6,
}

func (ignore Paralyze) ComputeOutput(c1 *SkillInput, c2 *SkillOutput) {
    return
    // Duration: int((3 + c1.INT[ATTR_CURRENT]/15) * 10),

}

func HitRoll(pAttacker CreatureInterface, pDefender CreatureInterface, bonus int) bool {
    if pDefender.CreatureInstance().isFlag(EFFECT_CLASS_NO_DAMAGE) {
        return false
    }

    var (
        tohit   ToHit_t
        defense Defense_t
    )
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
        tohit = ToHit_t(getPercentValue(int(tohit), VampireTimebandFactor[timeband]))
    case *Monster:
        tohit = pAttacker.(*Monster).ToHit
        tohit = ToHit_t(getPercentValue(int(tohit), VampireTimebandFactor[timeband]))
    }

    switch pDefender.(type) {
    case *Slayer:
        defense = pAttacker.(*Slayer).Defense[ATTR_CURRENT]
    case *Ouster:
        defense = pAttacker.(*Ouster).Defense[ATTR_CURRENT]
    case *Vampire:
        defense = pAttacker.(*Vampire).Defense[ATTR_CURRENT]
        defense = Defense_t(getPercentValue(int(defense), VampireTimebandFactor[timeband]))
    case *Monster:
        defense = pAttacker.(*Monster).Defense
        defense = Defense_t(getPercentValue(int(defense), VampireTimebandFactor[timeband]))
    }

    randValue := rand.Intn(100)
    var result int

    if int(tohit) >= int(defense) {
        result = min(90, int(int((float64(int(tohit)-int(defense))/1.5))+60)+bonus)
    } else {
        if _, ok := pAttacker.(*Monster); ok {
            result = max(10, (int)(60-int((float64(defense)-float64(tohit))/1.5)+bonus))
        } else {
            result = max(20, (int)(60-int((float64(defense)-float64(tohit))/1.5)+bonus))
        }
    }

    if randValue <= result {
        return true
    }

    return false
}
