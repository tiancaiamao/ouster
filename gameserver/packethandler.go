package main

import (
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
)

type PacketHandler func(pkt packet.Packet, agent *Agent)

var packetHandlers map[packet.PacketID]PacketHandler

func init() {
    packetHandlers = map[packet.PacketID]PacketHandler{
        packet.PACKET_CG_CONNECT:         CGConnectHandler,
        packet.PACKET_CG_READY:           CGReadyHandler,
        packet.PACKET_CG_ATTACK:          CGAttackHandler,
        packet.PACKET_CG_SAY:             CGReadyHandler,
        packet.PACKET_CG_MOVE:            CGMoveHandler,
        packet.PACKET_CG_SKILL_TO_SELF:   CGSkillToSelfHandler,
        packet.PACKET_CG_SKILL_TO_OBJECT: CGSkillToObjectHandler,
        packet.PACKET_CG_SKILL_TO_TILE:   CGSkillToTileHandler,
    }
}

func CGAttackHandler(pkt packet.Packet, agent *Agent) {
    // attack := pkt.(packet.CGAttackPacket)
    //    log.Println(" attack monster ", attack.ObjectID)
    //    target := player.Scene.objects[attack.ObjectID]
    //    if monster, ok := target.(*Monster); ok {
    //        hit := HitTest(player.ToHit, monster.Defense)
    //        if hit {
    //            player.send <- packet.GCAttackMeleeOK1{
    //                ObjectID: monster.Id(),
    //            }
    //
    //            damage := 1
    //            if player.Damage > monster.Protection {
    //                damage = int(player.Damage - monster.Protection)
    //            }
    //
    //            log.Println("send attack SkillOutput to scene..........")
    //            player.Scene.agent <- AgentMessage{
    //                Player: player,
    //                Msg: SkillOutput{
    //                    MonsterID: attack.ObjectID,
    //                    Damage:    damage,
    //                },
    //            }
    //        } else {
    //            player.send <- &packet.GCSkillFailed1Packet{}
    //        }
    //		}
}

func CGMoveHandler(pkt packet.Packet, agent *Agent) {
    agent.scene <- MoveMessage{
        Agent:        agent,
        CGMovePacket: pkt.(packet.CGMovePacket),
    }
}

func CGSkillToSelfHandler(pkt packet.Packet, agent *Agent) {
    if agent.PlayerStatus != GPS_NORMAL {
        return
    }

    // 检查变身狼状态一些技能不可用

    // if slayer, ok := agent.PlayerCreatureInterface.(Ouster); ok {

    // }

    skillPacket := pkt.(packet.CGSkillToSelfPacket)
    skillHandler, ok := skillTable[skillPacket.SkillType]
    if !ok {

    }

    if handler, ok := skillHandler.(SkillToSelfInterface); ok {
        handler.ExecuteToSelf(skillPacket, agent)
    }
}

func CGSkillToObjectHandler(pkt packet.Packet, agent *Agent) {
    if agent.PlayerStatus != GPS_NORMAL {
        return
    }

    // 检查变身狼状态一些技能不可用

    // if slayer, ok := agent.PlayerCreature.(Ouster); ok {
    //
    // }

    skillPacket := pkt.(packet.CGSkillToTilePacket)
    skillHandler, ok := skillTable[skillPacket.SkillType]
    if !ok {

    }

    if handler, ok := skillHandler.(SkillToTileInterface); ok {
        handler.ExecuteToTile(skillPacket, agent)
    }
}

func CGSkillToTileHandler(pkt packet.Packet, agent *Agent) {
    // skill := pkt.(packet.CGSkillToTilePacket)
    // player.SkillToTile(skill)
}

func CGConnectHandler(pkt packet.Packet, agent *Agent) {
    raw := pkt.(*packet.CGConnectPacket)
    pcItf, err := LoadPlayerCreature(raw.PCName, packet.PCType(raw.PCType))
    if err != nil {
        Log.Errorln("对CGConnectHandler的处理有问题")
    }
    agent.PlayerCreatureInterface = pcItf

    info := &packet.GCUpdateInfoPacket{
        // PCType: player.PCType,
        // PCInfo: player.PCInfo(),
        // ZoneID: player.Scene.ZoneID,
        // ZoneX:  player.X(),
        // ZoneY:  player.Y(),

        GameTime: packet.GameTimeType{
            Year:  1983,
            Month: 8,
            Day:   19,

            Hour:   12,
            Minute: 28,
            Second: 16,
        },

        DarkLevel:  13,
        LightLevel: 6,

        MonsterTypes: []uint16{5, 6, 7, 8},

        Premium: 17,
        NicknameInfo: packet.NicknameInfo{
            NicknameID: 32560,
        },

        GuildUnionUserType: 2,
    }

    // code := Encrypt(player.Scene.ZoneID, 1)
    // player.packetReader.Code = code
    // player.packetWriter.Code = code

    if info.PCType == 'O' {
        info.GearInfo = packet.GearInfo{
            GearSlotInfoList: []packet.GearSlotInfo{
                packet.GearSlotInfo{
                    ObjectID:   12494,
                    ItemClass:  59,
                    ItemType:   14,
                    Durability: 6700,
                    Grade:      4,
                    ItemNum:    1,

                    SlotID: 3,
                },
            },
        }
    }

    agent.sendPacket(info)
    agent.sendPacket(&packet.GCPetInfoPacket{})
    agent.PlayerStatus = GPS_WAITING_FOR_CG_READY
}

func CGReadyHandler(pkt packet.Packet, agent *Agent) {
    // pc := agent.pc.PlayerCreatureInstance()
    if agent.PlayerStatus != GPS_WAITING_FOR_CG_READY {

    }

    // scene := GetScene(pc.ZoneID)
    // scene.agent <- AgentMsg{}

    // var save chan<- AgentMsg
    // 地图切换
    if agent.ZoneID != 0 {
        // save = agent.send
    }

    agent.sendPacket(&packet.GCSetPositionPacket{
        // X:   player.X(),
        // Y:   player.Y(),
        Dir: 2,
    })

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
    agent.sendPacket(&skillInfo)
    agent.PlayerStatus = GPS_NORMAL
}
