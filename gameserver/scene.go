package main

import (
    "github.com/tiancaiamao/ouster/config"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
    . "github.com/tiancaiamao/ouster/util"
    "os"
    "path"
    "strings"
    "time"
)

// Scene是一个运行起来的地图场景，包含一个Zone成员
// Scene负责channel通信相关，然后调用Zone中对应的方法
type Scene struct {
    // 可以看作aoi管理
    Zone

    // 对象管理
    registerID ObjectID_t
    objects    map[ObjectID_t]ObjectInterface

    // 玩家管理
    players map[ObjectID_t]*Agent

    // NPC管理
    npcManager *NPCManager

    // 怪物管理
    monsterManager *MonsterManager

    // Effect管理
    effectManager *EffectManager

    // 天气管理

    quit  chan struct{}
    event chan interface{}
    agent chan AgentMessage
}

// 不需要加锁，但要保证这个函数只在scene的goroutine中运行
func (scene *Scene) registeObject(obj ObjectInterface) {
    scene.registerID++
    obj.ObjectInstance().ObjectID = scene.registerID
    scene.objects[scene.registerID] = obj
}

func (scene *Scene) getCreature(objectID ObjectID_t) CreatureInterface {
    monster := scene.monsterManager.getCreature(objectID)
    if monster != nil {
        return monster
    }

    agent, ok := scene.players[objectID]
    if ok {
        return agent
    }

    log.Debugln("没有找到creature")
    // pCreature = m_pNPCManager->getCreature(objectID);

    return nil
}

func NewScene(smp *data.SMP, ssi data.SSI) (ret *Scene, err error) {
    ret = new(Scene)
    ret.registerID = 10000
    ret.load(smp, ssi)
    ret.monsterManager = NewMonsterManager()
    ret.npcManager = NewNPCManager()
    ret.effectManager = NewEffectManager()
    ret.players = make(map[ObjectID_t]*Agent)
    ret.objects = make(map[ObjectID_t]ObjectInterface)

    ret.quit = make(chan struct{})
    ret.event = make(chan interface{})
    ret.agent = make(chan AgentMessage, 200)

    return
}

func (s *Scene) Init() error {
    err := s.monsterManager.Init(s)
    return err
}

func (m *Scene) Blocked(x, y uint16) bool {
    return false
}

func addPCToTile(scene *Scene, x, y int, pc *PlayerCreature, agent *Agent) {
    for sum := 0; sum < 20; sum++ {
        for i := -10; i < 10; i++ {
            for j := -10; j < 10; j++ {
                if abs(i)+abs(j) < sum {
                    if x+i >= int(scene.Width) || x+i < 0 {
                        continue
                    }
                    if y+i >= int(scene.Height) || y+i < 0 {
                        continue
                    }

                    tile := scene.Tile(x+i, y+j)
                    if !tile.HasCreature(pc.MoveMode) {
                        tile.addCreature(agent)
                        log.Debugf("login player: %d to (%d %d) tile=%#v\n", pc.ObjectID, x+i, y+j, tile)
                        pc.X = ZoneCoord_t(x + i)
                        pc.Y = ZoneCoord_t(y + j)
                        return
                    }
                }
            }
        }
    }
    panic("should not reach here!!")
}

func (m *Scene) Login(agent *Agent) {
    m.registeObject(agent)

    pc := agent.PlayerCreatureInstance()
    m.players[pc.ObjectID] = agent
    pc.Scene = m

    addPCToTile(m, int(pc.X), int(pc.Y), pc, agent)

    // log.Debugln("login的时候应该是加入了的: ", pc.ObjectID, m.players[pc.ObjectID])
    // obj := m.getCreature(pc.ObjectID)
    // log.Debugln("立马能取出来: ", pc.ObjectID, obj)
}

func (m *Scene) Logout(agent *Agent) {
    c := agent.CreatureInstance()

    m.Tile(int(c.X), int(c.Y)).deleteCreature(c.ObjectID)
    delete(m.players, c.ObjectID)

    gcDeleteObject := &packet.GCDeleteObjectPacket{
        ObjectID: c.ObjectID,
    }

    m.broadcastPacket(c.X, c.Y, gcDeleteObject, agent)

    // 最后写一条退出消息
    agent.sendPacket(packet.GCDisconnect{})
    // 然后关闭chan使得写的goroutine退出
    close(agent.send)
}

func (s *Scene) Loop() {
    heartbeat := time.Tick(200 * time.Millisecond)
    for {
        select {
        case msg, ok := (<-s.agent):
            if !ok {

            }
            s.processAgentMessage(msg)
        case <-s.quit:
        case <-s.event:
        case <-heartbeat:
            s.heartbeat()
        }
    }
}

func (s *Scene) heartbeat() {
    s.monsterManager.heartbeat()
    s.npcManager.heartbeat()
    s.effectManager.heartbeat(time.Now())
}

func (scene *Scene) addCreature(creature CreatureInterface, cx ZoneCoord_t, cy ZoneCoord_t, dir Dir_t) {
    pt, err := findSuitablePosition(&scene.Zone, cx, cy, creature.CreatureInstance().MoveMode)
    if err == nil {
        if creature.CreatureClass() == CREATURE_CLASS_MONSTER {
            monster := creature.(*Monster)
            scene.monsterManager.addCreature(monster)
            // pkt = packet.AddMonsterPacket(pMonster, NULL)
            // broadcastPacket(cx, cy, pAddMonsterPacket, pMonster)
        } else if creature.CreatureClass() == CREATURE_CLASS_NPC {
            npc := creature.(*NPC)
            scene.npcManager.addCreature(npc)

            // gcAddNPC := packet.GCAddNPCPacket{}
            // broadcastPacket(pt.x, pt.y, &gcAddNPC)
        }

        scene.Tile(pt.X, pt.Y).addCreature(creature)

        c := creature.CreatureInstance()
        c.X = ZoneCoord_t(pt.X)
        c.Y = ZoneCoord_t(pt.Y)
        c.Dir = dir

        c.Scene = scene
    } else {
        log.Debugln("应该是运行到这里来了，是不是？")
    }
}

// 这个函数只能在scene的goroutine中调用
// 被攻击者是怪物，直接减血
// 被攻击者是玩家，发到agent的goroutine中去计算
func (m *Scene) setDamage(target CreatureInterface, agent *Agent, damage Damage_t) {
    var status packet.GCStatusCurrentHP
    switch target.(type) {
    case *Agent:
        // TODO
    case *Monster:
        monster := target.(*Monster)
        if monster.HP[ATTR_CURRENT] < HP_t(damage) {
            // 这里只是设置dead标志，在MonsterManager的heartbeat中kill
            monster.HP[ATTR_CURRENT] = 0
            monster.LastKiller = agent.ObjectInstance().ObjectID
        } else {
            monster.HP[ATTR_CURRENT] -= HP_t(damage)
            status = packet.GCStatusCurrentHP{
                ObjectID:  monster.ObjectID,
                CurrentHP: monster.HP[ATTR_CURRENT],
            }
        }
    default:
        log.Errorln("参数不对")
    }

    pc := agent.CreatureInstance()

    // // 广播给所有玩家，攻击成功，怪物状态变化
    // ok3 := packet.GCAttackMeleeOK3{
    //     ObjectID:       pc.ObjectID,
    //     TargetObjectID: target.CreatureInstance().ObjectID,
    // }
    // m.broadcastPacket(pc.X, pc.Y, ok3, agent)

    m.broadcastPacket(pc.X, pc.Y, status, nil)
}

func (m *Scene) processAgentMessage(msg AgentMessage) {
    switch raw := msg.(type) {
    case LoginMessage:
        m.Login(raw.Agent)
        raw.wg.Done()
    case MoveMessage:
        m.movePC(raw.Agent, ZoneCoord_t(raw.X), ZoneCoord_t(raw.Y), Dir_t(raw.Dir))
    case FastMoveMessage:
        pc := raw.Agent.PlayerCreatureInstance()
        m.moveFastPC(raw.Agent, pc.X, pc.Y, ZoneCoord_t(raw.X), ZoneCoord_t(raw.Y), raw.SkillType)
    case LogoutMessage:
        m.Logout(raw.Agent)
    case DamageMessage:
        m.setDamage(raw.target, raw.Agent, raw.damage)

        // case *packet.GCFastMovePacket:
        // fastMove := msg.(*packet.GCFastMovePacket)
        // obj := m.objects[playerId]
        // player := obj.(*Player)
        // m.Update(player.Entity, fastMove.ToX, fastMove.ToY)
        // player.send <- fastMove
        // player.BroadcastPacket(player.X(), player.Y(), &packet.GCSkillToTileOK5{
        //     ObjectID:  playerId,
        //     SkillType: fastMove.SkillType,
        //     X:         player.X(),
        //     Y:         player.Y(),
        //     Duration:  10,
        // })
        // case SkillOutput:
        // skillOutput := msg.(SkillOutput)
        // id := skillOutput.MonsterID
        // obj := m.objects[id]
        // monster, _ := obj.(*Monster)
        // if monster.HP[ATTR_CURRENT] > uint16(skillOutput.Damage) {
        //     monster.HP[ATTR_CURRENT] -= uint16(skillOutput.Damage)
        //     m.BroadcastPacket(monster.X(), monster.Y(), packet.GCStatusCurrentHP{
        //         ObjectID:  monster.Id(),
        //         CurrentHP: monster.HP[ATTR_CURRENT],
        //     })
        // } else {
        //     m.BroadcastPacket(monster.X(), monster.Y(), &packet.GCAddMonsterCorpse{
        //         ObjectID:    id,
        //         MonsterType: monster.MonsterType,
        //         MonsterName: monster.Name,
        //         X:           monster.X(),
        //         Y:           monster.Y(),
        //         Dir:         2,
        //         LastKiller:  playerId,
        //     })
        //     m.BroadcastPacket(monster.X(), monster.Y(), packet.GCCreatureDiedPacket(monster.Id()))
        // }
    case SkillBroadcastMessage:
        pc := raw.Agent.PlayerCreatureInstance()
        m.broadcastPacket(pc.X, pc.Y, raw.Packet, raw.Agent)
    case MeteorStrikeMessage:
        tile := m.Tile(int(raw.X), int(raw.Y))
        tile.addEffect(&raw.EffectMeteorStrike)

        if tile.HasCreature(MOVE_MODE_WALKING) {
            target := tile.getCreature(MOVE_MODE_WALKING)
            // x := raw.X
            // y := raw.Y

            class := target.CreatureClass()
            if class == CREATURE_CLASS_SLAYER || class == CREATURE_CLASS_OUSTER {
                if canSee(target, msg.Sender()) {
                    // pc := target.CreatureInstance()
                    // ok2.ObjectID = pc.ObjectID
                    //   ok2.SkillType = raw.SkillType
                    //   ok2.X = pc.X
                    //   ok2.Y = pc.Y
                    //   ok2.Duration = raw.Duration
                    //   ok2.Range = Range
                    //   target.sendPacket(ok2)
                } else {

                    // target.sendPacket(ok6)
                }
            }
            if class == CREATURE_CLASS_MONSTER {
                // target.(*Monster).addEnemy()
            }
        }

        pc := raw.PlayerCreatureInstance()
        ok3 := &packet.GCSkillToTileOK3{
            ObjectID: raw.PlayerCreatureInstance().ObjectID,
            // SkillType: raw.SkillType,
            X:  Coord_t(pc.X),
            Y:  Coord_t(pc.Y),
        }

        ok4 := &packet.GCSkillToTileOK4{
            // SkillType: raw.SkillType,
            X:  Coord_t(pc.X),
            Y:  Coord_t(pc.Y),
            // Duration: raw.Duration,
        }

        m.broadcastPacket(pc.X, pc.Y, ok3, raw.Agent)
        m.broadcastPacket(ZoneCoord_t(raw.X), ZoneCoord_t(raw.Y), ok4, raw.Agent)
    }
}

var (
    g_Scenes map[ZoneID_t]*Scene = make(map[ZoneID_t]*Scene)
)

func Initialize() {
    dir, err := os.Open(config.DataFilePath)
    if err != nil {
        panic(err)
    }

    fi, err := dir.Readdir(0)
    if err != nil {
        panic(err)
    }

    for _, info := range fi {
        if info.IsDir() {
            continue
        }

        if strings.HasSuffix(info.Name(), ".smp") {
            smp, err := data.ReadSMP(path.Join(config.DataFilePath, info.Name()))
            if err != nil {
                log.Infof("加载SMP地图文件失败:%s, error: %s\n", info.Name(), err)
                continue
            }
            // log.Debugf("%s: %#v\n", info.Name(), smp)

            ssiFile := strings.Replace(info.Name(), ".smp", ".ssi", 0)
            ssi, err := data.ReadSSI(path.Join(config.DataFilePath, ssiFile))
            if err != nil {
                log.Infof("加载SSI地图文件失败:%s, error: %s\n", info.Name(), err)
                continue
            }

            if smp.ZoneID == 0 {
                // 暂时不加载自己改的的图,像L3 L4 L5 塔1 塔2
                continue
            }

            scene, err := NewScene(smp, ssi)
            if err != nil {
                log.Warnf("加载地图失败:%s, error:%s\n", info.Name(), err.Error())
                continue
            }

            err = scene.Init()
            if err != nil {
                log.Warnf("地图初始化失败:%s, error:%s\n", info.Name(), err.Error())
                continue
            }

            go scene.Loop()
            g_Scenes[scene.ZoneID] = scene
            log.Infof("加载地图%d成功:%s", scene.ZoneID, info.Name())
        }
    }
}
