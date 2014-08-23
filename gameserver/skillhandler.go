package main

import (
    "github.com/tiancaiamao/ouster/packet"
)

type SkillToSelfInterface interface {
    ExecuteToSelf(packet.CGSkillToSelfPacket, *Agent)
}

type SkillToObjectInterface interface {
    ExecuteToObject(packet.CGSkillToObjectPacket, *Agent)
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

func (invisibility Invisibility) ExecuteToSelf(skill packet.CGSkillToSelfPacket, agent *Agent) {
    // 基类的函数
    if !invisiblity.CheckAndDecreaseMana(skill.SkillType, agent, skillSlot) {
        receiver.executeSkillFailNormal(ouster, skill.SkillType)
        return
    }

    input := input(agent)
    var output SkillOutput

    invisibility.ComputeOutput(&input, &output)
    effect := new(EffectFadeOut)
    effect.Duration = output.Duration
    effect.Deadline = 40
    agent.addEffect(effect)
    agent.setFlag(EFFECT_CLASS_FADE_OUT)

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
        packet.GCSkillToSelfOK2{
            ObjectID:  agent.PlayerCreatureInstance().ObjectID,
            SkillType: SKILL_INVISIBILITY,
            Duration:  0,
        },
    }

    agent.scene <- ok2
}

func (meteor MeteorStrike) ExecuteToObject(skill packet.CGSkillToObjectPacket, agent *Agent) {
    target := agent.NearbyAgent(skill.ObjectID)
    meteor.ExecuteToTile(packet.CGSkillToObjectPacket{
        X:  target.X,
        Y:  target.Y,
    })
}

func (meteor MeteorStrike) ExecuteToTile(skill packet.CGSkillToObjectPacket, agent *Agent) {
    // 基类的函数
    if !invisiblity.CheckAndDecreaseMana(skill.SkillType, agent, skillSlot) {
        receiver.executeSkillFailNormal(ouster, skill.SkillType)
        return
    }

    var (
        ok1 packet.GCSkillToTileOK1
        ok2 packet.GCSkillToTileOK2
        ok3 packet.GCSkillToTileOK3
        ok4 packet.GCSkillToTileOK4
        ok5 packet.GCSkillToTileOK5
        ok6 packet.GCSkillToTileOK6
    )

    input := input(agent)
    var output SKillOutput
    meteor.ComputeOutput(&input, &output)

    agent.sendPacket(ok1)
    agent.scene <- MeteorStrikeMessage{
        UserObjectID: agent.PlayerCreatureInstance().ObjectID,
        Damage:       output.Damage,
        NextTime:     output.Duration,
        X:            skill.X,
        Y:            skill.Y,
    }
}

func (paralyze Paralyze) ExecuteToObject(skill packet.CGSkillToObjectPacket, agent *Agent) {

}
