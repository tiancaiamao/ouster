package main

import (
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "sync"
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
        packet.PACKET_CG_VERIFY_TIME:     CGVerifyTimeHandler,
        packet.PACKET_CG_LOGOUT:          CGLogoutHandler,
    }
}

func CGLogoutHandler(pkt packet.Packet, agent *Agent) {
    // 保存
    agent.save()
    // 从Zone去除
    agent.scene <- LogoutMessage{
        Agent: agent,
    }
}

func CGVerifyTimeHandler(pkt packet.Packet, agent *Agent) {}

func CGAttackHandler(pkt packet.Packet, agent *Agent) {
    fail := packet.GCSkillFailed1Packet{
        SkillType: SKILL_ATTACK_MELEE,
    }
    if agent.PlayerStatus != GPS_NORMAL {
        agent.sendPacket(&fail)
        return
    }

    pc := agent.PlayerCreatureInstance()
    zoneLevel := pc.Scene.getZoneLevel(pc.X, pc.Y)
    if (zoneLevel&ZoneLevel_t(COMPLETE_SAFE_ZONE)) != 0 ||
        pc.isFlag(EFFECT_CLASS_PARALYZE) ||
        pc.isFlag(EFFECT_CLASS_CAUSE_CRITICAL_WOUNDS) ||
        pc.isFlag(EFFECT_CLASS_EXPLOSION_WATER) ||
        pc.isFlag(EFFECT_CLASS_COMA) {
        agent.sendPacket(&fail)
        return
    }

    attack := pkt.(packet.CGAttackPacket)
    target, ok := pc.Scene.objects[attack.ObjectID]
    if !ok {
        agent.sendPacket(&fail)
        return
    }

    if target.ObjectClass() != OBJECT_CLASS_CREATURE {
        agent.sendPacket(&fail)
        return
    }
    targetCreature := target.(CreatureInterface)

    // ok3 := packet.GCAttackMeleeOK3{}

    // skillslot = agent.hasSkill(SKILL_ATTACK_MELEE)
    // timeCheck := verifyRunTime(skillslot)
    rangeCheck := verifyDistance(agent, targetCreature)
    hitRoll := HitRoll(agent, targetCreature, 0)

    if rangeCheck && hitRoll {
        damage := agent.PlayerCreatureInterface.computeDamage(targetCreature, false)

        // 这个伤害是要广播给地图周围玩家知道的
        agent.scene <- DamageMessage{
            target:   targetCreature,
            damage:   damage,
            critical: false,
        }

        if slayer, ok := agent.PlayerCreatureInterface.(*Slayer); ok {
            weapon := slayer.getWearItem(SLAYER_WEAR_RIGHTHAND)
            switch weapon.ItemClass() {
            case ITEM_CLASS_BLADE:
                increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
            case ITEM_CLASS_SWORD:
                increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
            case ITEM_CLASS_CROSS:
                increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
            case ITEM_CLASS_MACE:
                increaseDomainExp(slayer, SKILL_DOMAIN_BLADE, 1, packet.GCAttackMeleeOK1{}, targetCreature.CreatureInstance().Level)
            default:
                log.Errorln("武器不对!")
            }
        }

        ok1 := packet.GCAttackMeleeOK1{
            ObjectID: attack.ObjectID,
        }
        agent.sendPacket(ok1)

        switch target.(type) {
        case *Agent:
            targetAgent := target.(*Agent)
            targetAgent.sendPacket(packet.GCAttackMeleeOK2{
                ObjectID: agent.ObjectInstance().ObjectID,
            })
        case *Monster:
            monster := target.(*Monster)
            monster.addEnemy(agent)
        }

        // skillslot.setRunTime()
    } else {
        // 执行失败处理
    }

    // agent.setLastTarget(target.ObjectID)

    // if monster, ok := target.(*Monster); ok {
    //        hit := HitTest(player.ToHit, monster.Defense)
    //        if hit {
    //            player.send <- packet.GCAttackMeleeOK1{
    //								ObjectID: monster.Id(),
    //						}
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
    pcItf, zid, err := LoadPlayerCreature(raw.PCName, packet.PCType(raw.PCType))
    if err != nil {
        log.Errorln("LoadPlayerCreature失败了:", err)
    }

    agent.PlayerCreatureInterface = pcItf
    scene, ok := g_Scenes[zid]
    if !ok {
        log.Errorln("加载的agent所在的scene不存在:", zid)
        agent.ErrorClose()
        return
    }
    agent.scene = scene.agent

    msg := LoginMessage{
        Agent: agent,
        wg:    &sync.WaitGroup{},
    }
    // 向scene发消息并等待其返回
    msg.wg.Add(1)
    agent.scene <- msg
    msg.wg.Wait()

    log.Debugln("坐标：", agent.CreatureInstance().X, agent.CreatureInstance().Y)
    info := &packet.GCUpdateInfoPacket{
        PCInfo: agent.PCInfo(),
        ZoneID: zid,
        ZoneX:  Coord_t(agent.CreatureInstance().X),
        ZoneY:  Coord_t(agent.CreatureInstance().Y),

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

        MonsterTypes: []MonsterType_t{5, 6, 7, 8},

        Premium: 17,
        NicknameInfo: packet.NicknameInfo{
            NicknameID: 32560,
        },

        GuildUnionUserType: 2,
    }
    switch agent.PlayerCreatureInterface.(type) {
    case *Vampire:
        info.PCType = 'V'
    case *Ouster:
        info.PCType = 'O'
    case *Slayer:
        info.PCType = 'S'
    default:
        log.Errorln("agent类型不对!!")
    }

    code := Encrypt(uint16(agent.CreatureInstance().Scene.ZoneID), 1)
    agent.packetReader.Code = code
    agent.packetWriter.Code = code

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
