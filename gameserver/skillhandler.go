package main

import (
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
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
    pc.SetFlag(EFFECT_CLASS_FADE_OUT)

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
