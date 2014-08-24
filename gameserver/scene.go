package main

import (
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/packet"
    // "log"
    "math/rand"
    "time"
)

// Scene是一个运行起来的地图场景，包含一个Zone成员
// Scene负责channel通信相关，然后调用Zone中对应的方法
type Scene struct {
    objects []ObjectInterface

    // 玩家管理
    players map[ObjectID_t]*Agent
    // NPC管理
    // 怪物管理
    monsters []Monster
    // Effect管理
    // 天气管理

    // 可以看作aoi管理
    *Zone

    quit  chan struct{}
    event chan interface{}
    agent chan AgentMessage
}

func (scene *Scene) AddObject(obj ObjectInterface) uint32 {
    idx := len(scene.objects)
    if idx < 10000 {
        scene.objects = append(scene.objects, obj)
    } else {
    }
    return uint32(idx)
}

func NewScene(m *data.Map) *Scene {
    ret := new(Scene)

    // ret.Map = m
    // players := make([]*Player, 0, 200)

    num := 0
    for _, mi := range m.MonsterInfo {
        num += int(mi.Count)
    }
    // monsters := make([]Monster, num)
    ret.objects = make([]ObjectInterface, 0, num+200)

    idx := 0
    for _, mi := range m.MonsterInfo {
        // tp := data.MonsterType2MonsterInfo[mi.MonsterType]
        for i := 0; i < int(mi.Count); i++ {
            var x, y int
            // set monster's position and so on
            for {
                x = rand.Intn(int(m.Width))
                y = rand.Intn(int(m.Height))

                flag := m.Data[x*int(m.Width)+y]
                if flag == 0x0 && !ret.Blocked(uint16(x), uint16(y)) {
                    break
                }
            }

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
            idx++
        }
    }

    // ret.players = players
    // ret.monsters = monsters
    ret.quit = make(chan struct{})
    ret.event = make(chan interface{})
    ret.agent = make(chan AgentMessage, 200)

    return ret
}

func (m *Scene) Blocked(x, y uint16) bool {
    return false
}

func (m *Scene) Login(player *Player, zoneX uint8, zoneY uint8) error {
    // m.players = append(m.players, player)

    // id := m.AddObject(player)
    // player.Entity = m.Add(zoneX, zoneY, id)
    // player.Scene = m

    return nil
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
    maps      map[string]*Scene
    zoneTable map[uint16]*Scene
)

func Initialize() {
    maps = make(map[string]*Scene)
    maps["limbo_lair_se"] = NewScene(&data.LimboLairSE)
    maps["perona_nw"] = NewScene(&data.PeronaNW)

    // zoneTable = make(map[uint16]*Scene)
    // for _, m := range maps {
    //     zoneTable[m.ZoneID] = m
    //     go m.Loop()
    // }
}
