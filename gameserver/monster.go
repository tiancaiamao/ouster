package main

import (
    "fmt"
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    "github.com/tiancaiamao/ouster/packet"
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

    ticker uint16

    MonsterType MonsterType_t
    Name        string

    MainColor Color_t
    SubColor  Color_t

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

    MoveMode MoveMode

    Exp   Exp_t
    Moral Moral_t

    Delay       time.Duration
    AttackDelay time.Duration
    AccuDelay   time.Duration

    Enemies  map[ObjectID_t]CreatureInterface
    NextTurn time.Time
    Brain    *MonsterAI

    SilverDamage Silver_t

    LastHitCreatureClass CreatureClass

    LastKiller ObjectID_t

    RegenAmount   HP_t
    NextRegenTime time.Time
}

func (m *Monster) computeHP() {
    m.HP[ATTR_MAX] = HP_t(m.STR*2 + Attr_t(m.Level))
}

func (m *Monster) computeToHit() {
    m.ToHit = ToHit_t(float64(m.DEX) / 2 * (1.0 + float64(m.Level)/100))
}

func (m *Monster) computeDefense() {
    m.Defense = Defense_t(float64(m.DEX) / 2 * (1.0 + float64(m.Level)/100))
}

func (m *Monster) computeProtection() {
    m.Protection = Protection_t(float64(m.STR) / (5.0 - float64(m.Level)/100))
}

func (m *Monster) computeMinDamage() {
    m.Damage[ATTR_CURRENT] = Damage_t(float64(m.STR) / (6.0 - float64(m.Level)/100))
}

func (m *Monster) computeMaxDamage() {
    m.Damage[ATTR_MAX] = Damage_t(float64(m.STR) / (4.0 - float64(m.Level)/100))
}

func (m *Monster) getProtection() Protection_t {
    // TODO: 加入夜间强化
    return m.Protection
}

func NewMonster(monsterType MonsterType_t) *Monster {
    info, ok := data.MonsterInfoTable[monsterType]
    if !ok {
        log.Warnln("未知的monster类型:", monsterType)
        return nil
    }

    ret := &Monster{
        MonsterType:  monsterType,
        Name:         info.Name,
        STR:          info.STR,
        DEX:          info.DEX,
        INI:          info.INTE,
        Exp:          info.Exp,
        MeleeRange:   info.MeleeRange,
        MissileRange: info.MissileRange,
        MainColor:    info.MColor,
        SubColor:     info.SColor,
        MoveMode:     info.MoveMode,
    }

    ret.Sight = info.Sight
    ret.Level = info.Level
    ret.Delay = time.Duration(info.Delay) * time.Millisecond
    ret.AttackDelay = ret.Delay

    ret.computeHP()
    ret.HP[ATTR_CURRENT] = ret.HP[ATTR_MAX]

    ret.computeToHit()
    ret.computeDefense()
    ret.computeProtection()
    ret.computeMinDamage()
    ret.computeMaxDamage()

    ret.Brain = NewMonsterAI(ret, info.AIType)
    ret.Enemies = make(map[ObjectID_t]CreatureInterface)
    ret.Creature.Init()

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
    return max(abs(int(x-m.X)), abs(int(y-m.Y)))
}

func (m *Monster) getHP(attr int) HP_t {
    return m.HP[attr]
}

func (m *Monster) computeDamage(creature CreatureInterface, critical bool) Damage_t {
    minDamage := m.Damage[ATTR_CURRENT]
    maxDamage := m.Damage[ATTR_MAX]
    scope := int(maxDamage - minDamage)
    realDamage := max(1, int(minDamage)+rand.Intn(scope))

    // timeband    = getZoneTimeband(pMonster->getZone());
    timeband := 0
    realDamage = getPercentValue(realDamage, MonsterTimebandFactor[timeband])
    protection := creature.getProtection()
    return computeFinalDamage(minDamage, maxDamage, Damage_t(realDamage), protection, critical)
}

// TODO
func (m *Monster) getFleeRange() int {
    return 0
}

func (m *Monster) getMeleeRange() int {
    return m.MeleeRange
}

func (m *Monster) getZone() *Zone {
    return &m.Scene.Zone
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

func (m *Monster) canMove(nx ZoneCoord_t, ny ZoneCoord_t) bool {
    if !m.Creature.canMove(nx, ny) {
        return false
    }
    if m.Scene == nil {
        return false
    }

    tile := m.Scene.Tile(int(nx), int(ny))
    if tile == nil {
        return false
    }

    if tile.HasCreature(m.MoveMode) {
        return false
    }
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
    Count     int
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
            x, y := m.findPosition(&scene.Zone, monsterType)
            monster := NewMonster(monsterType)
            // log.Debugf("%p\n", monster)
            scene.registeObject(monster)
            scene.addCreature(monster, x, y, Dir_t(rand.Intn(DIR_MAX)))
        }
    }

    return nil
}

func (m *MonsterManager) findPosition(zone *Zone, monsterType MonsterType_t) (x ZoneCoord_t, y ZoneCoord_t) {
    info, ok := data.MonsterInfoTable[monsterType]
    if !ok {
        log.Error("不对，找不到monster信息")
        return
    }
    for i := 0; i < 300; i++ {
        pt := zone.getRandomMonsterRegenPosition()
        tile := zone.Tile(int(pt.X), int(pt.Y))

        if !tile.isBlocked(info.MoveMode) &&
            !tile.hasPortal() &&
            (*zone.Level(int(pt.X), int(pt.Y))&ZoneLevel_t(SAFE_ZONE)) == 0 {
            x = ZoneCoord_t(pt.X)
            y = ZoneCoord_t(pt.Y)
            return
        }
    }
    log.Errorln("地图中找不到可以放怪物的点了...不科学！")
    return
}

func (m *MonsterManager) addCreature(monster *Monster) {
    if monster.ObjectID == 0 {
        panic("!!!")
    }
    m.Monsters[monster.ObjectID] = monster
    m.Count++
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

        pkt := packet.GCStatusCurrentHP{
            ObjectID:  m.ObjectID,
            CurrentHP: m.HP[ATTR_CURRENT],
        }

        m.Scene.broadcastPacket(m.X, m.Y, &pkt, nil)
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
            pMonsterInfo := data.MonsterInfoTable[m.MonsterType]
            if rand.Intn(128) < pMonsterInfo.UnburrowChance {
                // TODO
                // SkillHandler * pSkillHandler =
                // g_pSkillHandlerManager.getSkillHandler(SKILL_UN_BURROW)
                // pSkillHandler.execute(this)
                m.addAccuDelay(time.Second + 500*time.Millisecond)
            } else {
                m.Brain.setDelay(currentTime)
            }
        } else if !m.isMaster {
            diceResult := rand.Intn(128)
            if diceResult < 6 {
                // 让怪走一走，锻炼身体
                // direction := rand.Intn(DIR_MAX)
                //    nx := m.X + ZoneCoord_t(dirMoveMask[direction].X)
                //    ny := m.Y + ZoneCoord_t(dirMoveMask[direction].Y)
                //
                //    if nx < m.Scene.Width && nx >= 0 && ny < m.Scene.Height && ny >= 0 {
                //        if m.canMove(nx, ny) && (m.Scene.getZoneLevel(nx, ny)&SAFE_ZONE) == 0 {
                //            m.Scene.moveCreature(m, nx, ny, Dir_t(direction))
                //        }
                //		}
            }

            // m.Brain.move(m_pZone.getWidth()>>1, m_pZone.getHeight()>>1)

            // if (m_bScanEnemy || isFlag(EFFECT_CLASS_HALLUCINATION)) && currentTime > m_NextScanTurn {
            //     m_pZone.mScan(this, m_X, m_Y, m_Dir)
            //
            //     m.NextScanTurn.tv_sec = currentTime.tv_sec + 2
            //     m.NextScanTurn.tv_usec = currentTime.tv_usec
            // }

            m.Brain.setDelay(currentTime)
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
    zone := &monster.Scene.Zone
    corpse := new(MonsterCorpse)
    corpse.ObjectID = monster.ObjectID
    corpse.CreateType = CREATE_TYPE_MONSTER
    corpse.X = monster.X
    corpse.Y = monster.Y
    corpse.MonsterType = monster.MonsterType
    corpse.Name = monster.Name
    corpse.LastKiller = monster.LastKiller

    zone.addItem(corpse, monster.X, monster.Y)
    zone.broadcastPacket(monster.X, monster.Y, packet.GCCreatureDiedPacket(monster.ObjectID), nil)
}

func (c *Monster) isAlive() bool {
    return c.HP[ATTR_CURRENT] > 0
}

func (m *Monster) getSight() Sight_t {
    return m.Sight
}

func (m *Monster) deleteAllEnemy() {
    // TODO
}

func (m *Monster) addEnemy(creature CreatureInterface) {
    m.Enemies[creature.CreatureInstance().ObjectID] = creature
}

func (m *Monster) hasEnemy() bool {
    return len(m.Enemies) != 0
}

func (m *Monster) getPrimaryEnemy() CreatureInterface {
    for _, v := range m.Enemies {
        // 随机挑选一个返回
        return v
    }
    return nil
}

func (m *Monster) hasNextMonsterSummonInfo() bool {
    // TODO
    return false
}

func (m *Monster) getDelay() time.Duration {
    return m.Delay
}

func (m *Monster) getAccuDelay() time.Duration {
    return m.AccuDelay
}

func (m *Monster) getAttackDelay() time.Duration {
    return m.AttackDelay
}

func (m *Monster) clearAccuDelay() {
    m.AccuDelay = 0
}

func (m *Monster) setNextTurn(t time.Time) {
    m.NextTurn = t
}
