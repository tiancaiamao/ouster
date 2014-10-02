package main

import (
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "time"
)

type SkillToSelfInterface interface {
    ExecuteToSelf(packet.CGSkillToSelfPacket, *Agent)
}

type SkillToObjectInterface interface {
    ExecuteToObject(CreatureInterface, CreatureInterface)
}

type SkillToTileInterface interface {
    ExecuteToTile(packet.CGSkillToTilePacket, *Agent)
}

type SkillHandler interface {
    ComputeOutput(*SkillInput, *SkillOutput)
}

type SkillInput struct {
    SkillLevel  int
    DomainLevel int
    STR         int
    DEX         int
    INT         int
    TargetType  int
    Ragge       int
    IClass      ItemClass
    PartySize   int
}

type SkillOutput struct {
    Damage   int
    Duration int
    Tick     int
    ToHit    int
    Range    int
    Delay    int
}

// TODO
func (skill BloodDrain) ExecuteToObject(sender CreatureInterface, target CreatureInterface) {
    log.Error("尚未实现")
}

// 注意:需要在agent的goroutine中执行的
func (melee AttackMelee) ExecuteToObject(sender CreatureInterface, target CreatureInterface) {
    rangeCheck := verifyDistance(sender, target)
    hitRoll := HitRoll(sender, target, 0)

    if rangeCheck && hitRoll {
        if agent, ok := sender.(*Agent); ok {
            damage := agent.computeDamage(target, false)
            // 这个伤害是要广播给地图周围玩家知道的
            agent.scene <- DamageMessage{
                Agent:    agent,
                target:   target,
                damage:   damage,
                critical: false,
            }

            // 发给攻击者，告诉他攻击成功了
            ok1 := packet.GCAttackMeleeOK1{
                ObjectID: target.CreatureInstance().ObjectID,
            }
            agent.sendPacket(ok1)

            if slayer, ok := agent.PlayerCreatureInterface.(*Slayer); ok {
                weapon := slayer.getWearItem(SLAYER_WEAR_RIGHTHAND)
                switch weapon.ItemClass() {
                case ITEM_CLASS_BLADE:
                    // increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
                case ITEM_CLASS_SWORD:
                    // increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
                case ITEM_CLASS_CROSS:
                    // increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
                case ITEM_CLASS_MACE:
                    // increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
                default:
                    log.Errorln("武器不对!")
                }
            }
        }

        if monster, ok := sender.(*Monster); ok {
            damage := monster.computeDamage(target, false)
            if agent, ok := target.(*Agent); ok {
                pc := agent.PlayerCreatureInstance()
                if pc.HP[ATTR_CURRENT] < HP_t(damage) {
                    // 玩家被打死了
                    log.Debugln("玩家被打死还没实现")
                } else {
                    pc.HP[ATTR_CURRENT] -= HP_t(damage)

                    log.Debugln("怪物攻击玩家，广播状态信息的攻击成功")
                    // 广播给所有玩家，攻击成功
                    ok3 := packet.GCAttackMeleeOK3{
                        ObjectID:       sender.CreatureInstance().ObjectID,
                        TargetObjectID: target.CreatureInstance().ObjectID,
                    }
                    pc.Scene.broadcastPacket(pc.X, pc.Y, ok3, agent)

                    // 广播给所有玩家，状态变化
                    status := packet.GCStatusCurrentHP{
                        ObjectID:  pc.ObjectID,
                        CurrentHP: pc.HP[ATTR_CURRENT],
                    }
                    pc.Scene.broadcastPacket(pc.X, pc.Y, status, nil)
                }
            } else {
                log.Errorln("参数不对")
            }
        }

        switch agent := target.(type) {
        case *Agent:
            agent.sendPacket(packet.GCAttackMeleeOK2{
                ObjectID: sender.CreatureInstance().ObjectID,
            })
        case *Monster:
            // monster := target.(*Monster)
            // monster.addEnemy(agent)
        }
    }
}

func (invisibility Invisibility) ExecuteToSelf(skill packet.CGSkillToSelfPacket, agent *Agent) {
    // 基类的函数
    // if !invisibility.CheckAndDecreaseMana(skill.SkillType, agent, skillSlot) {
    //     receiver.executeSkillFailNormal(ouster, skill.SkillType)
    //     return
    // }

    // input := input(agent)
    var input SkillInput
    var output SkillOutput

    invisibility.ComputeOutput(&input, &output)
    effect := new(EffectFadeOut)
    // effect.Duration = output.Duration
    // effect.Deadline = 40

    pc := agent.PlayerCreatureInstance()
    pc.addEffect(effect)
    pc.setFlag(EFFECT_CLASS_FADE_OUT)

    ok1 := &packet.GCSkillToSelfOK1{
        SkillType: SKILL_INVISIBILITY,
        CEffectID: skill.CEffectID,
        Duration:  0,
    }
    // ok.Short = make(map[packet.ModifyType]uint16)
    // ok.Short[12] = 180 + 256
    agent.sendPacket(ok1)

    ok2 := SkillBroadcastMessage{
        Agent: agent,
        // Packet: packet.GCSkillToSelfOK2{
        //     ObjectID:  agent.PlayerCreatureInstance().ObjectID,
        //     SkillType: SKILL_INVISIBILITY,
        //     Duration:  0,
        // },
    }

    agent.scene <- ok2
}

func (meteor MeteorStrike) ExecuteToObject(skill packet.CGSkillToObjectPacket, agent *Agent) {
    // target := agent.NearbyAgent(ObjectID_t(skill.TargetObjectID))
    meteor.ExecuteToTile(packet.CGSkillToObjectPacket{
    // X:  target.X,
    // Y:  target.Y,
    }, agent)
}

func (meteor MeteorStrike) ExecuteToTile(skill packet.CGSkillToObjectPacket, agent *Agent) {
    // 基类的函数
    // if !invisiblity.CheckAndDecreaseMana(skill.SkillType, agent, skillSlot) {
    //     receiver.executeSkillFailNormal(ouster, skill.SkillType)
    //     return
    // }

    // var (
    //     ok1 packet.GCSkillToTileOK1
    //     ok2 packet.GCSkillToTileOK2
    //     ok3 packet.GCSkillToTileOK3
    //     ok4 packet.GCSkillToTileOK4
    //     ok5 packet.GCSkillToTileOK5
    //     ok6 packet.GCSkillToTileOK6
    // )

    // input := input(agent)
    var input SkillInput
    var output SkillOutput
    meteor.ComputeOutput(&input, &output)

    // agent.sendPacket(ok1)
    agent.scene <- MeteorStrikeMessage{
    // UserObjectID: agent.PlayerCreatureInstance().ObjectID,
    // Damage:       output.Damage,
    // NextTime:     output.Duration,
    // X:            skill.X,
    // Y:            skill.Y,
    }
}

func (paralyze Paralyze) ExecuteToObject(skill packet.CGSkillToObjectPacket, agent *Agent) {

}

func (sharphail SharpHail) ExecuteToTile(skill packet.CGSkillToTilePacket, agent *Agent) {
    // weapon := agent.getWearItem(OUSTER_WEAR_RIGHTHAND)
    // if weapon == nil || weapon.ItemClass() != ITEM_CLASS_OUSTERS_CHAKRAM {
    //     // TODO
    //     return
    // }

    if !sharphail.Check(skill.SkillType, agent) {
        // TODO
        return
    }

    pc := agent.PlayerCreatureInstance()
    // skillslot := pc.SkillSlot[skill.SkillType]

    var input SkillInput
    var output SkillOutput
    // input.SkillLevel = skillslot.Level
    // input.DomainLevel =
    input.STR = int(pc.STR[ATTR_CURRENT])
    input.DEX = int(pc.DEX[ATTR_CURRENT])
    input.INT = int(pc.INI[ATTR_CURRENT])

    sharphail.ComputeOutput(&input, &output)
    scene := pc.Scene

    for x := skill.X - 2; x < skill.X+2; x++ {
        for y := skill.Y - 2; y < skill.Y+2; y++ {
            if x < 0 || ZoneCoord_t(x) >= scene.Width || y < 0 || ZoneCoord_t(y) >= scene.Height {
                continue
            }

            tile := scene.Tile(int(x), int(y))
            if !tile.canAddEffect() {
                continue
            }

            var creatureItf CreatureInterface
            if tile.HasCreature(MOVE_MODE_WALKING) {
                creatureItf = tile.getCreature(MOVE_MODE_WALKING)
            }

            damage := output.Damage
            // 技能伤害叠加基础伤害
            damage += int(agent.computeDamage(creatureItf, false))
            effect := new(EffectSharpHail)
            effect.UserObjectID = pc.ObjectID
            effect.Deadline = time.Now().Add(time.Duration(output.Duration) * time.Millisecond)
            effect.NextTime = time.Now().Add(3 * time.Millisecond)
            effect.Tick = output.Tick
            effect.Damage = damage / 3
            // effect.Level = skillslot.ExpLevel

            // TODO 加锁
            scene.registeObject(effect)
            scene.effectManager.addEffect(effect)
            tile.addEffect(effect)
        }
    }

    // ok1是回复攻击者技能施放成功
    agent.sendPacket(&packet.GCSkillToTileOK1{
        SkillType: skill.SkillType,
        CEffectID: skill.CEffectID,
        X:         skill.X,
        Y:         skill.Y,
        Duration:  uint16(output.Duration),
        Range:     5,
    })

    // ok5是发给即能看到施放者，又能看到tile的玩家
    ok5 := &packet.GCSkillToTileOK5{
        ObjectID:  pc.ObjectID,
        SkillType: skill.SkillType,
        X:         skill.X,
        Y:         skill.Y,
        Duration:  uint16(output.Duration),
    }

    scene.broadcastSkillPacket(pc.X, pc.Y, ZoneCoord_t(skill.X), ZoneCoord_t(skill.Y), ok5, agent)

    // ok3是向施法者周围广播施放成功
    scene.broadcastPacket(pc.X, pc.Y, &packet.GCSkillToTileOK3{
        ObjectID:  pc.ObjectID,
        SkillType: skill.SkillType,
        X:         skill.X,
        Y:         skill.Y,
    }, agent)

    // ok4是向tile周围广播
    scene.broadcastPacket(ZoneCoord_t(skill.X), ZoneCoord_t(skill.Y), &packet.GCSkillToTileOK4{
        SkillType: skill.SkillType,
        X:         Coord_t(skill.X),
        Y:         Coord_t(skill.Y),
        Duration:  uint16(output.Duration),
    }, agent)
}

func (sharphail SharpHail) ExecuteToObject(skill packet.CGSkillToObjectPacket, agent *Agent) {

}

func (sharphail SharpHail) ComputeOutput(input *SkillInput, output *SkillOutput) {
    output.Delay = 20
    if input.SkillLevel <= 15 {
        output.Damage = int(30 + ((input.DEX+input.STR)/20.0)*(1+(input.SkillLevel/15.0)))
    } else {
        output.Damage = int(30 + (float64(input.DEX+input.STR)/20)*(float64(5.0/3.0)+(float64(input.SkillLevel)/10.0)))
    }
    output.Tick = 3
    output.Duration = 10
}

func (spear DestructionSpear) ExecuteToObject(sender CreatureInterface, target CreatureInterface) {
    // weapon := sender.getWearItem(OUSTER_WEAR_RIGHTHAND)
    // if weapon == nil {
    // 	// TODO
    // 	return
    // }

    agent := sender.(*Agent)
    // pc := agent.PlayerCreatureInstance()
    spear.Check(SKILL_DESTRUCTION_SPEAR, agent)

    var input SkillInput
    var output SkillOutput

    spear.ComputeOutput(&input, &output)

    agent.scene <- DamageMessage{
        Agent:    agent,
        target:   target,
        damage:   Damage_t(output.Damage),
        critical: false,
    }

    // 发给攻击者，告诉他攻击成功了
    ok1 := &packet.GCSkillToObjectOK1{
        SkillType: SKILL_DESTRUCTION_SPEAR,
        // CEffectID      uint16
        TargetObjectID: target.CreatureInstance().ObjectID,
        Duration:       uint16(output.Duration),
        // Grade          uint8
        // ModifyInfo
    }
    agent.sendPacket(ok1)

    tc := target.CreatureInstance()
    scene := tc.Scene
    for x := max(0, int(tc.X)-1); x < min(int(scene.Width)-1, int(tc.X)+1); x++ {
        for y := max(0, int(tc.Y)-1); y < min(int(scene.Height)-1, int(tc.Y)+1); y++ {
            tile := scene.Tile(x, y)

            if tile.hasCreature() {
                // for _, v := range tile.Objects {
                // if creature, ok := v.(CreatureInterface); ok {
                // TODO 群体伤害的
                // }
                // }
            }
        }
    }
}

func (spear DestructionSpear) ComputeOutput(input *SkillInput, output *SkillOutput) {
    if input.SkillLevel <= 15 {
        output.Damage = min(60, 15+(input.DEX/20)+(input.SkillLevel/3))
        output.Duration = max(20, min(150, (int)((5.0+(float64(input.DEX)/30.0)*(1.0+(float64(input.SkillLevel)/22.5)))*10)))
    } else {
        output.Damage = min(60, 15+(input.DEX/20)+(input.SkillLevel/2))
        output.Duration = max(20, min(150, (int)((5.0+(float64(input.DEX)/30.0)*(float64(4.0/3.0)+(float64(input.SkillLevel)/45.0)))*10)))
    }

    output.Delay = 30 - input.DEX/6 - input.SkillLevel
    output.Delay = max(output.Delay, 6)
}

func (evade Evade) ExecuteToSelf(skill packet.CGSkillToSelfPacket, agent *Agent) {
    // weapon := agent.getWearItem(OUSTER_WEAR_RIGHTHAND)
    // if weapon == nil {
    // 	// TODO
    // 	return
    // }

    evade.Check(skill.SkillType, agent)

    var input SkillInput
    var output SkillOutput
    evade.ComputeOutput(&input, &output)

    pc := agent.PlayerCreatureInstance()

    eft := new(EffectEvade)
    eft.Deadline = time.Now().Add(time.Duration(output.Duration) * time.Millisecond)
    eft.Bonus = uint32(output.Damage)
    pc.addEffect(eft)
    pc.setFlag(EFFECT_CLASS_EVADE)

    agent.sendPacket(&packet.GCSkillToSelfOK1{
        SkillType: skill.SkillType,
        CEffectID: skill.CEffectID,
        Duration:  uint16(output.Duration),
    })

    pc.Scene.broadcastPacket(pc.X, pc.Y, &packet.GCSkillToSelfOK2{
        ObjectID:  pc.ObjectID,
        SkillType: skill.SkillType,
        Duration:  Duration_t(output.Duration),
    }, agent)

    pc.Scene.broadcastPacket(pc.X, pc.Y, &packet.GCAddEffect{
        ObjectID: pc.ObjectID,
        EffectID: EffectID_t(eft.EffectClass()),
        Duration: Duration_t(output.Duration),
    }, agent)

    // skillslot := pc.SkillSlot[skill.SkillType]
    // skillslot.setRunTime(output.Delay)
}

func ComputeOutput(input *SkillInput, output *SkillOutput) {
    if input.SkillLevel <= 15 {
        //		output.Damage = (int)((input.DEX / 10.0 ) * (1.0 + ((float)(input.SkillLevel) / 22.5 ) ));
        output.Damage = 10 + (int)(float64(input.SkillLevel)*16.0/9.0)
    } else {
        //		output.Damage = (int)((input.DEX / 10.0 ) * (4.0/3.0 + ((float)(input.SkillLevel) / 45.0 ) ));
        output.Damage = 10 + (int)(40.0/3.0+float64(input.SkillLevel)*8.0/9.0)
        if input.SkillLevel == 30 {
            output.Damage = (int)(float64(output.Damage) * 1.1)
        }
    }

    output.Damage = min(50, output.Damage)
    //	output.Duration = (int)((30.0 + ((float)(input.SkillLevel) * 3.0 ) ) * 10);
    output.Duration = (60 + input.DEX/10*(input.SkillLevel*2/3)) * 10
    output.Duration = min(1200, output.Duration)
    //	output.Delay = max (50 , output.Duration * 2 - (input.DEX * 2 ));
    output.Delay = output.Duration - (input.SkillLevel / 5)
}
