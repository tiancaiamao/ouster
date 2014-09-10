package main

import (
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "net"
    "time"
)

type Agent struct {
    PlayerCreatureInterface
    Player

    scene       chan<- AgentMessage
    computation chan func()
}

func NewAgent(conn net.Conn) *Agent {
    agent := new(Agent)
    InitPlayer(&agent.Player, conn)
    return agent
}

func (agent *Agent) Loop() {
    heartbeat := time.Tick(100 * time.Millisecond)
    for {
        select {
        case msg, ok := <-agent.client:
            if !ok {
                log.Debugln("客户端关了")
                return
            }
            log.Debugln("agent收到一个packet:", msg)
            agent.handleClientMessage(msg)
        case <-heartbeat:
            if agent.PlayerCreatureInterface != nil {
                agent.PlayerCreatureInstance().heartbeat()
            }
        case f, _ := <-agent.computation:
            f()
        }
    }
}

func (agent *Agent) handleClientMessage(pkt packet.Packet) {
    if pkt == nil {
        log.Errorln("不应该呀 怎么可能返回一个空")
    }
    handler, ok := packetHandlers[pkt.PacketID()]
    if !ok {
        log.Errorln("packet的handler未实现：", pkt.PacketID())
        return
    }

    handler(pkt, agent)
}

// called in scene
func (this *Player) handleAoiMessage(id uint32) {
    // obj := this.Scene.objects[id]
    // if _, ok := obj.(*Monster); ok {
    //     log.Println("it's a monster...send message")
    //     monster := obj.(*Monster)
    //     if _, ok := this.nearby[id]; !ok {
    //         this.nearby[id] = struct{}{}
    //
    //         addMonster := &packet.GCAddMonster{
    //             ObjectID:    uint32(id),
    //             MonsterType: monster.MonsterType,
    //             MonsterName: "test",
    //             X:           monster.X(),
    //             Y:           monster.Y(),
    //             Dir:         2,
    //             CurrentHP:   monster.HP[ATTR_CURRENT],
    //             MaxHP:       monster.MaxHP(),
    //         }
    //
    //         this.send <- addMonster
    //         monster.flag |= flagActive
    //         log.Println("monster ", id, "set to active", monster.flag)
    //         monster.Enemies = append(monster.Enemies, this.PacketID())
    //     } else {
    //
    //     }
    // }
}

func decreaseConsumeMP(agent *Agent) int {
    // TODO
    return 0
}

func decreaseMana(agent *Agent, mana int) {
    // TODO
}

func (agent *Agent) hasRankBonus() bool {
    // TODO
    return false
}

func hasEnoughMana(agent *Agent, requireMP int) bool {
    // TODO
    return true
}

func (agent *Agent) NearbyAgent(id ObjectID_t) *Agent {
    return nil
}
