package main

import (
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "sync"
)

// 传interface相当于传指针，消息更小
// 增加了解包的代价
type AgentMessage interface {
    Sender() *Agent
}

func (agent *Agent) Sender() *Agent {
    return agent
}

type MoveMessage struct {
    *Agent
    *packet.CGMovePacket
}

type FastMoveMessage struct {
    *Agent
    X         ZoneCoord_t
    Y         ZoneCoord_t
    SkillType SkillType_t
}

type DamageMessage struct {
    *Agent
    target   CreatureInterface
    damage   Damage_t
    critical bool
}

type LoginMessage struct {
    *Agent
    wg  *sync.WaitGroup
}

type LogoutMessage struct {
    *Agent
    // wg  *sync.WaitGroup
}

type SkillBroadcastMessage struct {
    *Agent
    packet.Packet
}

type MeteorStrikeMessage struct {
    *Agent
    EffectMeteorStrike
    X   uint8
    Y   uint8
}
