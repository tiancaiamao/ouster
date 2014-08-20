package main

import (
    "github.com/tiancaiamao/ouster/packet"
    "log"
)

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

func CGReadyHandler(pkt packet.Packet, agent *Agent) {
    // say := pkt.(*packet.CGSayPacket)
    // log.Println("say:", say.Message)
}

func CGMoveHandler(pkt packet.Packet, agent *Agent) {
    // player.Scene.agent <- AgentMessage{
    // Player: player,
    // Msg:    pkt,
    // }
}

func CGSkillToSelfHandler(pkt packet.Packet, agent *Agent) {
    pcItf := agent.pc
    pc := pcItf.PlayerCreatureInstance()
    player := pc.Player

    skill := pkt.(packet.CGSkillToSelfPacket)
    switch skill.SkillType {
    case SKILL_INVISIBILITY:
        ok := &packet.GCSkillToSelfOK1{
            SkillType: SKILL_INVISIBILITY,
            CEffectID: 181,
            Duration:  0,
            Grade:     0,
        }
        ok.Short = make(map[packet.ModifyType]uint16)
        ok.Short[12] = 180 + 256
        player.send <- ok
    default:
        log.Println("unknown SkillToSelf type:", skill.SkillType)
    }
}

func CGSkillToObjectHandler(pkt packet.Packet, agent *Agent) {
    // skill := pkt.(packet.CGSkillToObjectPacket)
    // player.SkillToObject(skill)

}

func CGSkillToTileHandler(pkt packet.Packet, agent *Agent) {
    // skill := pkt.(packet.CGSkillToTilePacket)
    // player.SkillToTile(skill)
}

func CGConnectHandler(pkt packet.Packet, agent *Agent) {
    pcItf := agent.pc
    pc := pcItf.PlayerCreatureInstance()
    player := pc.Player

    raw := pkt.(*packet.CGConnectPacket)
    player.Load(raw.PCName, packet.PCType(raw.PCType))

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

    player.send <- info
    player.send <- &packet.GCPetInfoPacket{}
}
