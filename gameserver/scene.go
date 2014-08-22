package main

import (
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/packet"
    "log"
    "math/rand"
    "time"
)

type AgentMessage struct {
    Player *Player
    Msg    interface{}
}

// Scene是一个运行起来的地图场景，包含一个Zone成员
// Scene负责channel通信相关，然后调用Zone中对应的方法
type Scene struct {
    objects []Object

    // 玩家管理
    // NPC管理
    // 怪物管理
    // Effect管理
    // 天气管理

    players map[ObjectID_t]*Agent

    monsters []Monster

    // 可以看作aoi管理
    zone *Zone

    quit  chan struct{}
    event chan interface{}
    agent chan AgentMessage
}

func (scene *Scene) AddObject(obj Object) uint32 {
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
    ret.objects = make([]Object, 0, num+200)

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
        case <-s.agent:
            // m.processPlayerInput(data.Player.Id(), data.Msg)
        case <-s.quit:
        case <-s.event:
        case <-heartbeat:
            s.zone.heartbeat()
        }
    }
}

func (m *Scene) processPlayerInput(playerId uint32, msg interface{}) {
    switch msg.(type) {
    case packet.CGMovePacket:
        move := msg.(packet.CGMovePacket)
        log.Println("scene receive a CGMovePacket:", move.X, move.Y, move.Dir)
        // obj := m.objects[playerId]
        // player := obj.(*Player)

        // if move.Dir >= 8 {
        //     moveErr := packet.GCMoveErrorPacket{
        //         player.X(),
        //         player.Y(),
        //     }
        //     player.send <- moveErr
        // }

        move.X = uint8(int(move.X) + dirMoveMask[move.Dir].X)
        move.Y = uint8(int(move.Y) + dirMoveMask[move.Dir].Y)

        // m.Update(player.Entity, move.X, move.Y)
        // player.send <- packet.GCMoveOKPacket{
        //     X:   move.X,
        //     Y:   move.Y,
        //     Dir: move.Dir,
        // }
    case *packet.GCFastMovePacket:
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
    case SkillOutput:
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
