package main

import (
    // "github.com/tiancaiamao/ouster/packet"
    "fmt"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
    "math/rand"
    "time"
)

const (
    flagDead = 1 << iota
    flagActive
)

type Monster struct {
    Creature

    Zone *Zone

    ticker uint16

    MonsterType MonsterType_t
    Name        string

    STR Attr_t
    DEX Attr_t
    INI Attr_t

    HP         [2]HP_t
    Defense    Defense_t
    Protection Protection_t
    ToHit      ToHit_t
    Damage     [2]Damage_t

    MeleeRange   int
    MissileRange int

    EffectManager EffectManager

    IsEventMonster bool
    IsChief        bool
    isMaster       bool
    bTreasure      bool

    Exp   Exp_t
    Moral Moral_t

    Delay       Turn_t
    AttackDelay Turn_t
    AccuDelay   time.Duration

    Enemies  []ObjectID_t
    NextTurn time.Time
    Brain    *MonsterAI

    SilverDamage Silver_t

    LastHitCreatureClass CreatureClass

    LastKiller ObjectID_t

    RegenAmount   HP_t
    NextRegenTime time.Time
}

func NewMonster(monsterType MonsterType_t) *Monster {
    info, ok := data.MonsterInfoTable[monsterType]
    if !ok {
        log.Warnln("未知的monster类型:", monsterType)
        return nil
    }

    ret := &Monster{
        MonsterType: monsterType,
        Name:        info.Name,
        // Level:        info.Level,
        STR:          info.STR,
        DEX:          info.DEX,
        INI:          info.INTE,
        Exp:          info.Exp,
        MeleeRange:   info.MeleeRange,
        MissileRange: info.MissileRange,
    }

    ret.HP[ATTR_CURRENT] = info.HP
    ret.HP[ATTR_MAX] = info.HP

    return ret
}

func (m *Monster) CreatureClass() CreatureClass {
    return CREATURE_CLASS_MONSTER
}

func (c *Monster) isFlag(effect uint) bool {
    return c.Creature.IsFlag(effect)
}

func (m *Monster) addAccuDelay(time.Duration) {
    // TODO
}

func (m *Monster) getDistance(x, y ZoneCoord_t) int {
    // TODO
    return 0
}

// TODO
func (m *Monster) getFleeRange() int {
    return 0
}

func (m *Monster) getMeleeRange() int {
    // TODO
    return 0
}

func (m *Monster) getZone() *Zone {
    return m.Zone
}

func (m *Monster) getX() ZoneCoord_t {
    return m.X
}

func (m *Monster) getY() ZoneCoord_t {
    return m.Y
}

func (m *Monster) getDir() Dir_t {
    return m.Dir
}

// TODO
func canMove(nx ZoneCoord_t, ny ZoneCoord_t) bool {
    return true
}

func dir(dx int, dy int) uint8 {
    var ret uint8
    switch {
    case dx > 0 && dy > 0:
        ret = RIGHTDOWN
    case dx > 0 && dy == 0:
        ret = RIGHT
    case dx > 0 && dy < 0:
        ret = RIGHTUP
    case dx < 0 && dy > 0:
        ret = LEFTDOWN
    case dx < 0 && dy == 0:
        ret = LEFT
    case dx < 0 && dy > 0:
        ret = LEFTUP
    case dx == 0 && dy > 0:
        ret = DOWN
    case dx == 0 && dy < 0:
        ret = UP
    }
    return ret
}

type MonsterManager struct {
    Monsters  map[ObjectID_t]*Monster
    RegenTime time.Time
}

func NewMonsterManager() *MonsterManager {
    ret := &MonsterManager{
        Monsters:  make(map[ObjectID_t]*Monster),
        RegenTime: time.Now(),
    }
    return ret
}

func (m *MonsterManager) Init(scene *Scene) error {
    zoneInfo, ok := data.ZoneInfoTable[scene.ZoneID]
    if !ok {
        err := fmt.Errorf("ZoneInfo not found: %d", scene.ZoneID)
        return err
    }

    for _, v := range zoneInfo.MonsterList {
        monsterType := v.MonsterType
        count := v.Count

        for i := 0; i < count; i++ {
            monster := NewMonster(monsterType)
            scene.registeObject(monster)
            scene.addCreature(monster, monster.X, monster.Y, Dir_t(rand.Intn(DIR_MAX)))
        }
    }

    return nil
}

func (m *MonsterManager) addCreature(monster *Monster) {
    m.Monsters[monster.ObjectID] = monster
}

func (manager *MonsterManager) heartbeat() {
    now := time.Now()
    for key, monster := range manager.Monsters {
        monster.EffectManager.heartbeat(now)

        if monster.isAlive() {
            monster.heartbeat(now)
        } else {
            delete(manager.Monsters, key)
            manager.killCreature(monster)
        }

        if now.After(manager.RegenTime) {
            manager.regenerateCreatures()
            manager.RegenTime = now.Add(5 * time.Second)
        }
    }
}

func (m *Monster) heartbeat(currentTime time.Time) {
    if m.RegenAmount != 0 && currentTime.After(m.NextRegenTime) {
        if float64(m.HP[ATTR_MAX])*0.3 > float64(m.HP[ATTR_CURRENT]) {
            switch m.MonsterType {
            case 724:
                fallthrough
            case 725:
                m.NextRegenTime.Add(time.Second)
                m.RegenAmount = 100
            case 717:
                m.NextRegenTime.Add(time.Second)
                m.RegenAmount = 200
            case 723:
                m.NextRegenTime.Add(time.Second)
                m.RegenAmount = 300
            default:
                m.NextRegenTime.Add(3 * time.Second)
            }
        } else {
            switch m.MonsterType {
            case 724:
                fallthrough
            case 725:
                fallthrough
            case 717:
                fallthrough
            case 723:
                m.NextRegenTime.Add(2 * time.Second)
            case 721:
            default:
                m.NextRegenTime.Add(7 * time.Second)
            }
        }

        m.HP[ATTR_CURRENT] += HP_t(min(int(m.RegenAmount), int(m.HP[ATTR_MAX]-m.HP[ATTR_CURRENT])))

        // var gcHP packet.GCStatusCurrentHP
        // gcHP.setObjectID(getObjectID())
        // gcHP.setCurrentHP(m_HP[ATTR_CURRENT])
        //
        // m.getZone().broadcastPacket(getX(), getY(), &gcHP)
    }

    if currentTime.Before(m.NextTurn) {
        return
    }

    if m.AccuDelay != 0 {
        m.NextTurn.Add(m.AccuDelay)
        // m.clearAccuDelay()
        return
    }

    m.verifyEnemies()

    if m.Brain == nil {
        m.addAccuDelay(time.Second + 500*time.Millisecond)
        return
    }

    if m.isFlag(EFFECT_CLASS_PARALYZE) ||
        m.isFlag(EFFECT_CLASS_COMA) ||
        m.isFlag(EFFECT_CLASS_CAUSE_CRITICAL_WOUNDS) ||
        m.isFlag(EFFECT_CLASS_SLEEP) ||
        m.isFlag(EFFECT_CLASS_ARMAGEDDON) ||
        m.isFlag(EFFECT_CLASS_TRAPPED) ||
        m.isFlag(EFFECT_CLASS_EXPLOSION_WATER) {
        m.Brain.setDelay(currentTime)
        return
    }

    if m.hasEnemy() {
        pEnemy := m.getPrimaryEnemy()
        if pEnemy != nil {
            m.Brain.Deal(pEnemy, currentTime)
        }
    } else {
        if m.isFlag(EFFECT_CLASS_HIDE) {
            // pMonsterInfo := g_pMonsterInfoManager.getMonsterInfo(m_MonsterType)
            //
            // if (rand() & 0x0000007F) < pMonsterInfo.getUnburrowChance() {
            //     SkillHandler * pSkillHandler = g_pSkillHandlerManager.getSkillHandler(SKILL_UN_BURROW)
            //     Assert(pSkillHandler != NULL)
            //
            //     pSkillHandler.execute(this)
            //
            //     addAccuDelay(time.Second + 500*time.MilliSecond)
            // } else {
            //     m.Brain.setDelay(currentTime)
            // }
        } else if !m.isMaster && m.Zone.ZoneID != 72 {
            // if m.RelicIndex == -1 {
            //     pt := POINT(m.X, m.Y)
            //
            //     VSRect * pOuterRect = m_pZone.getOuterRect()
            //     VSRect * pInnerRect = m_pZone.getInnerRect()
            //     VSRect * pCoreRect = m_pZone.getCoreRect()
            //
            //     if pCoreRect.ptInRect(pt) || pInnerRect.ptInRect(pt) {
            //
            //         diceResult := rand() & 0x0000007F //%100;
            //         if diceResult < 6 {
            //             direction := rand() & 0x00000007 //% 8;
            //             nx := pt.x + dirMoveMask[direction].x
            //             ny := pt.y + dirMoveMask[direction].y
            //
            //             if canMove(nx, ny) && !(m_pZone.getZoneLevel(nx, ny) & SAFE_ZONE) {
            //                 m_pZone.moveCreature(this, nx, ny, direction)
            //             }
            //         }
            //     } else if pOuterRect.ptInRect(pt) {
            //         m.Brain.move(m_pZone.getWidth()>>1, m_pZone.getHeight()>>1)
            //     }
            //
            //     if (m_bScanEnemy || isFlag(EFFECT_CLASS_HALLUCINATION)) && currentTime > m_NextScanTurn {
            //         m_pZone.mScan(this, m_X, m_Y, m_Dir)
            //
            //         m.NextScanTurn.tv_sec = currentTime.tv_sec + 2
            //         m.NextScanTurn.tv_usec = currentTime.tv_usec
            //     }
            // }
            // m.Brain.setDelay(currentTime)
        }
    }

    statSum := m.STR + m.DEX + m.INI
    if rand.Intn(600) < int(statSum) && m.isAlive() {
        m.HP[ATTR_CURRENT] = HP_t(min(int(m.HP[ATTR_CURRENT])+1, int(m.HP[ATTR_MAX])))
    }

    if m.HP[ATTR_CURRENT] == m.HP[ATTR_MAX] && m.isFlag(EFFECT_CLASS_BLOOD_DRAIN) {
        m.removeFlag(EFFECT_CLASS_BLOOD_DRAIN)
    }
}

// TODO
func (m *Monster) verifyEnemies() {
    // body
}

func (manager *MonsterManager) regenerateCreatures() {

}

func (manager *MonsterManager) killCreature(monster *Monster) {
    // 广播消息之类的 TODO

}

func (c *Monster) isAlive() bool {
    return c.HP[ATTR_CURRENT] > 0
}

// TODO
func (m *Monster) getSight() Sight_t {
    return 0
}

func (m *Monster) deleteAllEnemy() {
    // TODO
}

func (m *Monster) addEnemy(CreatureInterface) {
    // TODO
}

func (m *Monster) hasEnemy() bool {
    // TODO
    return false
}

// TODO
func (m *Monster) getPrimaryEnemy() CreatureInterface {
    return nil
}

func (m *Monster) hasNextMonsterSummonInfo() bool {
    // TODO
    return false
}
