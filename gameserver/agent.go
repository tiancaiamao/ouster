package main

import (
    "github.com/tiancaiamao/ouster/packet"
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
    pc := agent.PlayerCreatureInstance()
    heartbeat := time.Tick(100 * time.Millisecond)
    for {
        select {
        case msg, ok := <-agent.client:
            if !ok {
                // 网络出问题了，将玩家踢下线
                return
            } else {
                agent.handleClientMessage(msg)
            }
        case <-heartbeat:
            pc.heartbeat()
        case f, _ := <-agent.computation:
            f()
        }
    }
}

func (agent *Agent) handleClientMessage(pkt packet.Packet) {
    handler, ok := packetHandlers[pkt.PacketID()]
    if !ok {
        // TODO
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
