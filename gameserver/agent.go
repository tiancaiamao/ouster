package main

import (
    "github.com/tiancaiamao/ouster/packet"
    "log"
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

func (agent *Agent) handleClientMessage(pkt packet.Packet) {
    pcItf := agent.pc
    pc := pcItf.PlayerCreatureInstance()
    player := pc.Player

    switch pkt.Id() {
    case packet.PACKET_CG_CONNECT:
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
        player.send <- &skillInfo
    case packet.PACKET_CG_MOVE:
        // player.Scene.agent <- AgentMessage{
        // Player: player,
        // Msg:    pkt,
        // }
    case packet.PACKET_CG_SAY:
        // say := pkt.(*packet.CGSayPacket)
        // log.Println("say:", say.Message)
    case packet.PACKET_CG_ATTACK:
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
        //    }
    case packet.PACKET_CG_SKILL_TO_SELF:
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
    case packet.PACKET_CG_SKILL_TO_OBJECT:
        skill := pkt.(packet.CGSkillToObjectPacket)
        player.SkillToObject(skill)
    case packet.PACKET_CG_SKILL_TO_TILE:
        skill := pkt.(packet.CGSkillToTilePacket)
        player.SkillToTile(skill)
    case packet.PACKET_CG_BLOOD_DRAIN:
    case packet.PACKET_CG_VERIFY_TIME:
    case packet.PACKET_CG_LOGOUT:
        // player.Save()
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
    //         monster.Enemies = append(monster.Enemies, this.Id())
    //     } else {
    //
    //     }
    // }
}
