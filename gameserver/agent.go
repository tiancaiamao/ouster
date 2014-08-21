package main

import (
    "github.com/tiancaiamao/ouster/packet"
    "time"
)

type Agent struct {
    pc          PlayerCreatureInterface
    computation chan func()
}

func (agent *Agent) Loop() {
    pcItf := agent.pc
    pc := pcItf.PlayerCreatureInstance()
    player := pc.Player
    // creature := pc.Creature
    heartbeat := time.Tick(100 * time.Millisecond)
    for {
        select {
        case msg, ok := <-player.client:
            if !ok {
                // kick the player off...
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

type PacketHandler func(pkt packet.Packet, agent *Agent)

var packetHandlers map[packet.PacketID]PacketHandler

func init() {
    packetHandlers = map[packet.PacketID]PacketHandler{
        packet.PACKET_CG_CONNECT:         CGConnectHandler,
        packet.PACKET_CG_ATTACK:          CGAttackHandler,
        packet.PACKET_CG_SAY:             CGReadyHandler,
        packet.PACKET_CG_MOVE:            CGMoveHandler,
        packet.PACKET_CG_SKILL_TO_SELF:   CGSkillToSelfHandler,
        packet.PACKET_CG_SKILL_TO_OBJECT: CGSkillToObjectHandler,
        packet.PACKET_CG_SKILL_TO_TILE:   CGSkillToTileHandler,
    }
}

func (agent *Agent) handleClientMessage(pkt packet.Packet) {
    pcItf := agent.pc
    pc := pcItf.PlayerCreatureInstance()
    player := pc.Player

    handler, ok := packetHandlers[pkt.PacketID()]
    if !ok {

    }

    handler(pkt, agent)
    switch pkt.PacketID() {
    case packet.PACKET_CG_CONNECT:

    case packet.PACKET_CG_READY:
        // log.Println("get a CG Ready Packet!!!")
        player.send <- &packet.GCSetPositionPacket{
            // X:   player.X(),
            // Y:   player.Y(),
            Dir: 2,
        }

        var skillInfo packet.GCSkillInfoPacket
        // switch player.PCType {
        // case 'V':
        //     skillInfo.PCType = packet.PC_VAMPIRE
        // case 'O':
        //     skillInfo.PCType = packet.PC_OUSTER
        // case 'S':
        //     skillInfo.PCType = packet.PC_SLAYER
        // }
        // skillInfo.PCSkillInfoList = []packet.SkillInfo{
        // player.SkillInfo(),
        // }
        player.sendPacket(&skillInfo)
        return
    }
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
