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

    registerID ObjectID_t

    objects []ObjectInterface

    // 玩家管理
    players map[ObjectID_t]*Agent

    // NPC管理

    // 怪物管理
    monsterManager MonsterManager

    // Effect管理
    // 天气管理

    quit  chan struct{}
    event chan interface{}
    agent chan AgentMessage
}

func (scene *Scene) registerObject(obj ObjectInterface) {
    // 不需要加锁，但要保证这个函数只在scene的goroutine中运行
    scene.registerID++
    obj.ObjectInstance().ObjectID = scene.registerID
}

func NewScene(smp *data.SMP, ssi data.SSI) *Scene {
    ret := new(Scene)
    ret.registerID = 10000
    ret.load(smp, ssi)

    // num := 0
    //  for _, mi := range ret.MonsterInfo {
    //      num += int(mi.Count)
    //  }
    // monsters := make([]Monster, num)
    // ret.objects = make([]ObjectInterface, 0, num+200)

    // idx := 0
    //    for _, mi := range ret.MonsterInfo {
    //        // tp := data.MonsterType2MonsterInfo[mi.MonsterType]
    //        for i := 0; i < int(mi.Count); i++ {
    //            var x, y int
    //            // set monster's position and so on
    //            for {
    //                x = rand.Intn(int(ret.Width))
    //                y = rand.Intn(int(ret.Height))
    //
    //                flag := ret.Data[x*int(ret.Width)+y]
    //                if flag == 0x0 && !ret.Blocked(uint16(x), uint16(y)) {
    //                    break
    //                }
    //						}

    // monster := &monsters[idx]
    // id := ret.AddObject(monster)
    // monster.Entity = aoi.Add(uint8(x), uint8(y), id)
    // monster.MonsterType = mi.MonsterType
    // monster.STR[ATTR_CURRENT], monster.STR[ATTR_BASE] = tp.STR, tp.STR
    // monster.DEX[ATTR_CURRENT], monster.DEX[ATTR_BASE] = tp.DEX, tp.DEX
    // monster.INT[ATTR_CURRENT], monster.DEX[ATTR_BASE] = tp.INTE, tp.INTE
    // monster.Defense = monster.DEX[ATTR_CURRENT] / 2
    // monster.Protection = monster.STR[ATTR_CURRENT]
    // monster.HP[ATTR_MAX] = tp.STR*4 + uint16(tp.Level)
    // monster.HP[ATTR_CURRENT] = monster.HP[ATTR_MAX]
    // idx++
    // }
    // }

    ret.players = make(map[ObjectID_t]*Agent)
    // ret.monsters = monsters
    ret.quit = make(chan struct{})
    ret.event = make(chan interface{})
    ret.agent = make(chan AgentMessage, 200)

    return ret
}

func (m *Scene) Blocked(x, y uint16) bool {
    return false
}

func (m *Scene) Login(agent *Agent) {
    m.registerObject(agent)

    c := agent.CreatureInstance()
    m.players[c.ObjectID] = agent
    c.Scene = m
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

func (m *Scene) processAgentMessage(msg AgentMessage) {
    switch raw := msg.(type) {
    case LoginMessage:
        m.Login(raw.Agent)
        raw.wg.Done()
    case MoveMessage:
        m.movePC(raw.Agent, ZoneCoord_t(raw.X), ZoneCoord_t(raw.Y), Dir_t(raw.Dir))
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
        tile.AddEffect(&raw.EffectMeteorStrike)

        if tile.HasCreature(MOVE_MODE_WALKING) {
            target := tile.GetCreature(MOVE_MODE_WALKING)
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
            ObjectID: uint32(raw.PlayerCreatureInstance().ObjectID),
            // SkillType: raw.SkillType,
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
        }

        ok4 := &packet.GCSkillToTileOK4{
            // SkillType: raw.SkillType,
            X:  uint8(pc.X),
            Y:  uint8(pc.Y),
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

            scene := NewScene(smp, ssi)
            if scene != nil {
                go scene.Loop()
                g_Scenes[scene.ZoneID] = scene
                log.Infof("加载地图%d成功:%s", scene.ZoneID, info.Name())
            } else {
                log.Infof("加载地图失败:%s", info.Name())
            }
        }
    }
}
