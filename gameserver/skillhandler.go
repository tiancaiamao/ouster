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

func (receiver Invisibility) ExecuteToSelf(skill packet.CGSkillToSelfPacket, agent *Agent) {
    ok := &packet.GCSkillToSelfOK1{
        SkillType: SKILL_INVISIBILITY,
        CEffectID: 181,
        Duration:  0,
        Grade:     0,
    }
    ok.Short = make(map[packet.ModifyType]uint16)
    ok.Short[12] = 180 + 256
    agent.pc.PlayerCreatureInstance().sendPacket(ok)
}
