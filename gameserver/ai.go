package main

import (
    "github.com/tiancaiamao/ouster/data"
    "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
    "math/rand"
    "time"
)

type MonsterAI struct {
    Body         *Monster
    DirectiveSet data.DirectiveSet
    LastAction   int
    MoveRule     MoveRule
    BlockedDir   Dir_t
    WallCount    int
    bDamaged     bool
    Panic        int
    PanicMax     int
    Courage      int
    CourageMax   int
}

func NewMonsterAI(body *Monster, aiType int) *MonsterAI {
    ret := &MonsterAI{
        Body:       body,
        LastAction: LAST_ACTION_NONE,
        MoveRule:   MOVE_RULE_NORMAL,
        BlockedDir: DIR_NONE,
        Panic:      5,
        PanicMax:   5,
        Courage:    20,
        CourageMax: 20,
    }

    ret.DirectiveSet = data.DirectiveSetTable[aiType]
    return ret
}

type ConditionCheckFunction func(*Monster, CreatureInterface) bool

var (
    condChecker [DIRECTIVE_COND_MAX]ConditionCheckFunction
)

type MoveRule int8

const (
    MOVE_RULE_NORMAL MoveRule = iota
    MOVE_RULE_LEFTWALL
    MOVE_RULE_RIGHTWALL
    MOVE_RULE_MAX
)

const (
    LAST_ACTION_NONE = iota
    LAST_ACTION_MOVE
    LAST_ACTION_SKILL
    LAST_ACTION_SKILL_MAX
)

func init() {
    condChecker[DIRECTIVE_COND_ENEMY_RANGE_MELEE] = checkEnemyRangeMelee
    condChecker[DIRECTIVE_COND_ENEMY_RANGE_MISSILE] = checkEnemyRangeMissile
    condChecker[DIRECTIVE_COND_ENEMY_RANGE_CLOSE] = checkEnemyRangeClose
    condChecker[DIRECTIVE_COND_ENEMY_RANGE_OUT_OF_SIGHT] = checkEnemyRangeOutOfSight
    condChecker[DIRECTIVE_COND_ENEMY_DYING] = checkEnemyDying
    condChecker[DIRECTIVE_COND_ENEMY_NOT_BLOOD_DRAINED] = checkEnemyNotBloodDrained
    condChecker[DIRECTIVE_COND_ENEMY_NOT_GREEN_POISONED] = checkEnemyNotGreenPoisoned
    condChecker[DIRECTIVE_COND_ENEMY_NOT_YELLOW_POISONED] = checkEnemyNotYellowPoisoned
    condChecker[DIRECTIVE_COND_ENEMY_NOT_DARKBLUE_POISONED] = checkEnemyNotDarkbluePoisoned
    condChecker[DIRECTIVE_COND_ENEMY_NOT_GREEN_STALKERED] = checkEnemyNotGreenStalkered
    condChecker[DIRECTIVE_COND_ENEMY_NOT_PARALYZED] = checkEnemyNotParalyzed
    condChecker[DIRECTIVE_COND_ENEMY_NOT_DOOMED] = checkEnemyNotDoomed
    condChecker[DIRECTIVE_COND_ENEMY_NOT_BLINDED] = checkEnemyNotBlinded
    condChecker[DIRECTIVE_COND_ENEMY_NOT_IN_DARKNESS] = checkEnemyNotInDarkness
    condChecker[DIRECTIVE_COND_ENEMY_NOT_SEDUCTION] = checkEnemyNotSeduction
    condChecker[DIRECTIVE_COND_IM_OK] = checkImOK
    condChecker[DIRECTIVE_COND_IM_DYING] = checkImDying
    condChecker[DIRECTIVE_COND_IM_DAMAGED] = checkImDamaged
    condChecker[DIRECTIVE_COND_IM_HIDING] = checkImHiding
    condChecker[DIRECTIVE_COND_IM_WOLF] = checkImWolf
    condChecker[DIRECTIVE_COND_IM_BAT] = checkImBat
    condChecker[DIRECTIVE_COND_IM_INVISIBLE] = checkImInvisible
    condChecker[DIRECTIVE_COND_IM_WALKING_WALL] = checkImWalkingWall
    condChecker[DIRECTIVE_COND_TIMING_BLOOD_DRAIN] = checkTimingBloodDrain
    condChecker[DIRECTIVE_COND_MASTER_SUMMON_TIMING] = checkMasterSummonTiming
    condChecker[DIRECTIVE_COND_MASTER_NOT_READY] = checkMasterNotReady
    condChecker[DIRECTIVE_COND_IM_IN_BAD_POSITION] = checkImInBadPosition
    condChecker[DIRECTIVE_COND_FIND_WEAK_ENEMY] = checkFindWeakEnemy
    condChecker[DIRECTIVE_COND_ENEMY_NOT_DEATH] = checkEnemyNotDeath
    condChecker[DIRECTIVE_COND_ENEMY_NOT_HALLUCINATION] = checkEnemyNotHallucination
    condChecker[DIRECTIVE_COND_TIMING_MASTER_BLOOD_DRAIN] = checkTimingMasterBloodDrain
    condChecker[DIRECTIVE_COND_TIMING_DUPLICATE_SELF] = checkTimingDuplicateSelf
    condChecker[DIRECTIVE_COND_ENEMY_RANGE_IN_MISSILE] = checkEnemyRangeInMissile
    condChecker[DIRECTIVE_COND_POSSIBLE_SUMMON_MONSTERS] = checkPossibleSummonMonsters
    condChecker[DIRECTIVE_COND_ENEMY_TILE_NOT_ACID_SWAMP] = checkEnemyTileNotAcidSwamp
    condChecker[DIRECTIVE_COND_ENEMY_ON_AIR] = checkEnemyOnAir
    condChecker[DIRECTIVE_COND_ENEMY_ON_SAFE_ZONE] = checkEnemyOnSafeZone
    condChecker[DIRECTIVE_COND_CAN_ATTACK_THROWING_AXE] = checkCanAttackThrowingAxe
}

func abs(x int) int {
    if x > 0 {
        return x
    }
    return -x
}

func max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func (ai *MonsterAI) checkCondition(condition int, pEnemy CreatureInterface) bool {
    return condChecker[condition](ai.Body, pEnemy)
}

func (ai *MonsterAI) checkDirective(pDirective *data.Directive, pEnemy CreatureInterface) bool {
    if pDirective == nil {
        return false
    }

    rValue := true

    for _, condition := range pDirective.Conditions {
        if ai.checkCondition(condition, pEnemy) == false {
            rValue = false
            break
        }
    }

    return rValue
}

func (ai *MonsterAI) getDamaged() bool {
    return ai.bDamaged
}

func (ai *MonsterAI) getWallCount() int {
    return ai.WallCount
}

func (ai *MonsterAI) moveWall(ex ZoneCoord_t, ey ZoneCoord_t, nx *ZoneCoord_t, ny *ZoneCoord_t, ndir *Dir_t, bLeft bool) bool {
    //TODO
    return false
}

func flee(pEnemy CreatureInterface) bool {
    // TODO
    return false
}

func (ai *MonsterAI) moveNormal(ex ZoneCoord_t, ey ZoneCoord_t, nx *ZoneCoord_t, ny *ZoneCoord_t, ndir *Dir_t) bool {
    bestDir := Dir_t(DIR_NONE)
    curDir := ai.Body.getDir()

    if ai.Body.getX() < ex {
        if ai.Body.getY() < ey {
            bestDir = RIGHTDOWN
        } else if ai.Body.getY() > ey {
            bestDir = RIGHTUP
        } else {
            bestDir = RIGHT
        }
    } else if ai.Body.getX() > ex {
        if ai.Body.getY() < ey {
            bestDir = LEFTDOWN
        } else if ai.Body.getY() > ey {
            bestDir = LEFTUP
        } else {
            bestDir = LEFT
        }
    } else {
        if ai.Body.getY() < ey {
            bestDir = DOWN
        } else if ai.Body.getY() > ey {
            bestDir = UP
        } else {
            bestDir = DIR_NONE
        }
    }

    diffLevel := 0
    diff := 0
    found := false
    *ndir = bestDir
    *nx = ZoneCoord_t(int(ai.Body.getX()) + dirMoveMask[*ndir].X)
    *ny = ZoneCoord_t(int(ai.Body.getY()) + dirMoveMask[*ndir].Y)

    var bBlocked [DIR_MAX]bool
    DIR_MAX_1 := Dir_t(DIR_MAX - 1)

    bCanMove := ai.Body.canMove(*nx, *ny)

    if !bCanMove {
        bBlocked[*ndir] = true
        if !ai.Body.isBlockedByCreature(*nx, *ny) &&
            !ai.Body.isFlag(EFFECT_CLASS_TRANSFORM_TO_BAT) &&
            !ai.Body.isFlag(EFFECT_CLASS_HIDE) {
            ai.setMoveRule(MOVE_RULE_RIGHTWALL)
            ai.BlockedDir = bestDir
        }

        for {
            diffLevel++
            if found || diffLevel > 8 {
                break
            }

            if (diffLevel & 0x00000001) == 0 {
                diff = diffLevel >> 1
                dir1 := (ai.Body.getDir() + Dir_t(diff)) & DIR_MAX_1
                dir2 := (ai.Body.getDir() + DIR_MAX - Dir_t(diff)) & DIR_MAX_1

                if (abs(int(*ndir)+int(DIR_MAX)-int(dir1)) & int(DIR_MAX_1)) < (abs(int(*ndir)+int(DIR_MAX)-int(dir2)) & int(DIR_MAX_1)) {
                    *ndir = dir1
                } else {
                    diff = -diff
                    *ndir = dir2
                }
            } else {
                *ndir = Dir_t(int(ai.Body.getDir()) + int(DIR_MAX) - diff)
            }

            *ndir &= DIR_MAX_1

            *nx = ZoneCoord_t(int(ai.Body.getX()) + dirMoveMask[*ndir].X)
            *ny = ZoneCoord_t(int(ai.Body.getY()) + dirMoveMask[*ndir].Y)

            if ai.Body.canMove(*nx, *ny) {
                /*
                	if (ai.Body.isMaster())
                	{
                		Tile& tile = ai.Body.getZone().getTile(nx, ny);
                		if (tile.hasEffect())
                		{
                			if (tile.getEffect(Effect::EFFECT_CLASS_ACID_SWAMP)!=NULL
                				&& tile.getEffect(Effect::EFFECT_CLASS_GROUND_ATTACK)!=NULL)
                			{
                				found = true;
                			}
                		}
                	}
                	else
                */
                found = true
            }

            bBlocked[*ndir] = true
        }
    } else {
        found = true
    }

    if found && ai.MoveRule != MOVE_RULE_NORMAL {
        leftWall := bBlocked[(*ndir+2)&DIR_MAX_1]
        rightWall := bBlocked[(*ndir+DIR_MAX-2)&DIR_MAX_1]

        if leftWall && rightWall {
            if *ndir > curDir && *ndir < curDir+4 || curDir > 4 && (*ndir > curDir || *ndir < curDir-4) {
                ai.setMoveRule(MOVE_RULE_RIGHTWALL)
            } else {
                ai.setMoveRule(MOVE_RULE_LEFTWALL)
            }
        } else if leftWall {
            ai.setMoveRule(MOVE_RULE_LEFTWALL)

        } else if rightWall {
            ai.setMoveRule(MOVE_RULE_RIGHTWALL)

        } else {
            ai.setMoveRule(MOVE_RULE_NORMAL)
        }
    } else {
        ai.setMoveRule(MOVE_RULE_NORMAL)
    }

    return found
}

func (ai *MonsterAI) moveCoord(ex ZoneCoord_t, ey ZoneCoord_t) bool {
    pZone := ai.Body.getZone()

    var (
        nx   ZoneCoord_t
        ny   ZoneCoord_t
        ndir Dir_t
    )

    found := false
    switch ai.MoveRule {
    case MOVE_RULE_NORMAL:
        found = ai.moveNormal(ex, ey, &nx, &ny, &ndir)
    case MOVE_RULE_LEFTWALL:
        found = ai.moveWall(ex, ey, &nx, &ny, &ndir, true)
    case MOVE_RULE_RIGHTWALL:
        found = ai.moveWall(ex, ey, &nx, &ny, &ndir, false)
    }

    if found && (pZone.getZoneLevel(nx, ny)&ZoneLevel_t(SAFE_ZONE)) == 0 {
        pZone.moveCreature(ai.Body, nx, ny, ndir)
    }

    ai.LastAction = LAST_ACTION_MOVE

    return true
}

func (ai *MonsterAI) move(pEnemy CreatureInterface) bool {
    pZone := ai.Body.getZone()
    enemy := pEnemy.CreatureInstance()
    enemyX := enemy.X
    enemyY := enemy.Y
    myX := ai.Body.X
    myY := ai.Body.Y
    xOffset := enemyX - myX
    yOffset := enemyY - myY
    ex := enemy.X
    ey := enemy.Y

    ////////////////////////////////////////////////////////////
    // (enemyX, enemyY)
    //
    //                  myX, myY
    //
    //                           (ex, ey)
    //
    ////////////////////////////////////////////////////////////
    xOffset2 := xOffset << 1
    yOffset2 := yOffset << 1

    if enemyX+xOffset2 < 0 {
        ex = 0
    } else if enemyX+xOffset2 > pZone.getWidth() {
        ex = pZone.getWidth()
    } else {
        ex = enemyX + xOffset2
    }

    if enemyY+yOffset2 < 0 {
        ey = 0
    } else if enemyY+yOffset2 > pZone.getHeight() {
        ey = pZone.getHeight()
    } else {
        ey = enemyY + yOffset2
    }

    ai.setMoveRule(MOVE_RULE_NORMAL)
    return ai.moveCoord(ex, ey)
}

func (ai *MonsterAI) setMoveRule(rule MoveRule) {
    ai.MoveRule = rule
    ai.WallCount = 0
}

func (ai *MonsterAI) setDelay(currentTime time.Time) {
    delay := ai.Body.getDelay()

    if ai.Body.isFlag(EFFECT_CLASS_TRANSFORM_TO_BAT) {
        delay = 200
    } else if ai.Body.isFlag(EFFECT_CLASS_TRANSFORM_TO_WOLF) {
        delay = 300
    }

    modifier := rand.Intn(41) - 20

    delay = delay * Turn_t(1000)
    delay = delay + delay*Turn_t(modifier)/100

    nexttime := delay / 1000000 * time.Second
    nexttime.Add((delay % 1000000) * time.Microsecond)

    nexttime = nexttime + ai.Body.getAccuDelay()

    ai.Body.clearAccuDelay()

    if ai.Body.isFlag(EFFECT_CLASS_ICE_FIELD_TO_CREATURE) ||
        ai.Body.isFlag(EFFECT_CLASS_JABBING_VEIN) {
        ai.Body.setNextTurn(currentTime.Add(2 * nexttime))
    } else {
        ai.Body.setNextTurn(currentTime.Add(nexttime))
    }
}

func (ai *MonsterAI) setAttackDelay(currentTime time.Time) {
    delay := ai.Body.getAttackDelay()
    modifier = rand.Intn(21)

    delay := delay * 1000
    delay = delay + delay*modifier/100

    nexttime := delay / 1000000 * time.Second
    nexttime.Add(delay % 1000000 * time.Microsecond)
    nexttime.Add(ai.Body.getAccuDelay())

    ai.Body.clearAccuDelay()

    if ai.Body.isFlag(EFFECT_CLASS_ICE_OF_SOUL_STONE) {
        ai.Body.setNextTurn(currentTime.Add(2 * nexttime))
    } else {
        ai.Body.setNextTurn(currentTime.Add(nexttime))
    }
}

func (ai *MonsterAI) approach(pEnemy CreatureInterface) bool {
    return ai.move(pEnemy)
}

func (ai *MonsterAI) useSkill(pEnemy CreatureInterface, SkillType SkillType_t, ratio int) int {
    // enemy := pEnemy.CreatureInstance()
    // ex := enemy.X
    // ey := enemy.Y
    // dist := ai.Body.getDistance(ex, ey)

    if rand.Intn(100) >= ratio {
        // return SKILL_FAILED_RATIO
        return -1
    }

    if ai.Body.isFlag(EFFECT_CLASS_HIDE) {
        SkillType = SKILL_UN_BURROW
    } else if ai.Body.isFlag(EFFECT_CLASS_TRANSFORM_TO_BAT) {
        SkillType = SKILL_UN_TRANSFORM
    } else if ai.Body.isFlag(EFFECT_CLASS_INVISIBILITY) {
        SkillType = SKILL_UN_INVISIBILITY
    }

    skill, ok := skillTable[SkillType]
    if !ok {
        log.Errorf("技能%d的handler没有实现!!", SkillType)
        return 0
    }
    handler := skill.(SkillToObjectInterface)

    // 移动计算，以闭包形式发到agent的goroutine中运行
    if agent, ok := pEnemy.(*Agent); ok {
        closure := func() {
            handler.ExecuteToObject(ai.Body, pEnemy)
        }
        agent.computation <- closure
    } else {
        log.Errorln("怪物打怪物还没实现")
    }

    ai.LastAction = LAST_ACTION_SKILL
    return 0
}

func isValidZoneCoord(zone *Zone, x, y ZoneCoord_t) bool {
    // TODO
    return true
}

func (ai *MonsterAI) Deal(pEnemy CreatureInterface, currentTime time.Time) {
    for _, pDirective := range ai.DirectiveSet.Directives {
        if ai.checkDirective(pDirective, pEnemy) {
            switch pDirective.Action {
            case DIRECTIVE_ACTION_APPROACH:
                log.Debugln("动作是approach")
                ai.approach(pEnemy)
            case DIRECTIVE_ACTION_FLEE:
                log.Debugln("动作是逃跑")
                if !flee(pEnemy) {
                    ai.setMoveRule(MOVE_RULE_NORMAL)
                    rValue := ai.useSkill(pEnemy, SKILL_ATTACK_MELEE, 100)
                    if rValue != 0 {
                        ai.approach(pEnemy)
                    }
                }
            case DIRECTIVE_ACTION_USE_SKILL:
                log.Debugln("动作是放技能")
                if ai.Body.isFlag(EFFECT_CLASS_BLOCK_HEAD) || ai.Body.isFlag(EFFECT_CLASS_TENDRIL) {
                    continue
                }
                parameter := pDirective.Parameter
                ratio := pDirective.Ratio
                rValue := ai.useSkill(pEnemy, SkillType_t(parameter), ratio)
                if rValue != 0 {
                    break
                }
                ai.setMoveRule(MOVE_RULE_NORMAL)
            case DIRECTIVE_ACTION_FORGET:
                // log.Debugln("动作是忘记敌人")
                if len(ai.Body.Enemies) != 0 {
                    // ai.Body.Enemies = ai.Body.Enemies[1:]
                }
                ai.setMoveRule(MOVE_RULE_NORMAL)
            case DIRECTIVE_ACTION_CHANGE_ENEMY:
                log.Debugln("动作是切换敌人")
                ratio := pDirective.Parameter
                if rand.Intn(100) >= ratio {
                    break
                }

                // pNewEnemy := ai.Body.Enemies[0]
                // 从ObjectID得到Creature
                // if pNewEnemy != nil {
                // pEnemy = pNewEnemy
                // } else {
                if len(ai.Body.Enemies) == 0 {
                    // ai.Body.addEnemy(pEnemy)
                }
                // }
            case DIRECTIVE_ACTION_MOVE_RANDOM:
                log.Debugln("动作是随机移动")
                ratio := pDirective.Parameter
                if rand.Intn(100) >= ratio {
                    break
                }

                // var (
                // x   ZoneCoord_t
                // y   ZoneCoord_t
                // p   TPOINT
                // )

                // x = ai.Body.X
                // y = ai.Body.Y

                // p = getSafeTile(ai.Body.Zone, x, y)

                // if p.x == -1 {
                // break
                // }

                // x1 := p.x
                //								 y1 := p.y
                //								 if x != x1 || y != y1 {
                //										 move(x1, y1)
                //								 }
            case DIRECTIVE_ACTION_WAIT:
                log.Debugln("动作是等待")
                delay := 2 * time.Second
                ai.Body.addAccuDelay(delay)
            case DIRECTIVE_ACTION_FAST_FLEE:
                log.Debugln("动作是快速逃跑")
                result := false
                myX := ai.Body.X
                myY := ai.Body.Y
                nmX := pEnemy.CreatureInstance().X
                nmY := pEnemy.CreatureInstance().Y
                diffX := myX - nmX
                diffY := myY - nmY
                ratio := 5.0 / ZoneCoord_t(abs(int(diffX))+abs(int(diffY)))
                newX := (myX + diffX*ratio)
                newY := (myY + diffY*ratio)

                if isValidZoneCoord(&ai.Body.Scene.Zone, newX, newY) {
                    result = ai.Body.Scene.moveFastMonster(ai.Body, myX, myY, newX, newY, SKILL_RAPID_GLIDING)
                }
                if !result {
                    break
                }
            case DIRECTIVE_ACTION_SAY:
                log.Debugln("动作是说话")
                // parameter = pDirective.Parameter();
                // GCSay gcSay
                // gcSay.ObjectID = ai.Body.ObjectID
                // gcSay.setMessage(g_pStringPool.getString(parameter ));
                // gcSay.setColor(0x00ffffff);
                // ai.Body.getZone().broadcastPacket(ai.Body.getX(), ai.Body.getY(),
                // &gcSay);
            }
        }
    }

    switch ai.LastAction {
    case LAST_ACTION_NONE:
        ai.setDelay(currentTime)
    case LAST_ACTION_MOVE:
        ai.setDelay(currentTime)
    case LAST_ACTION_SKILL:
        ai.setAttackDelay(currentTime)
    }

    ai.LastAction = LAST_ACTION_NONE

    if (rand.Int() & 0x0000007F) > 64 { //%100 > 50)
        if ai.bDamaged {
            ai.bDamaged = false
            ai.Panic = ai.PanicMax
            ai.Courage = ai.CourageMax
        }
    }
}

func checkEnemyRangeMelee(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }

    enemy := pEnemy.CreatureInstance()
    dist := pMonster.getDistance(enemy.X, enemy.Y)
    if dist <= pMonster.getMeleeRange() {
        return true
    }
    return false
}

func checkEnemyRangeMissile(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    e := pEnemy.CreatureInstance()
    dist := pMonster.getDistance(e.X, e.Y)
    if dist > pMonster.getMeleeRange() && dist <= int(pMonster.getSight()) {
        return true
    }
    return false
}

func checkEnemyRangeInMissile(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    e := pEnemy.CreatureInstance()
    dist := pMonster.getDistance(e.X, e.Y)
    if dist <= int(pMonster.getSight()) {
        return true
    }
    return false
}

func checkEnemyRangeClose(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    e := pEnemy.CreatureInstance()
    dist := pMonster.getDistance(e.X, e.Y)
    if dist <= 1 {
        return true
    }
    return false
}

func checkEnemyRangeOutOfSight(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    dist := pMonster.getDistance(pEnemy.CreatureInstance().X, pEnemy.CreatureInstance().Y)
    if dist > int(pMonster.getSight()) {
        // log.Debugf("dist=%d, sight=%d\n", dist, pMonster.getSight())
        return true
    }
    return false
}

func checkEnemyDying(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }

    EnemyCurHP := pEnemy.getHP(ATTR_CURRENT)
    EnemyMaxHP := pEnemy.getHP(ATTR_MAX)

    if EnemyCurHP*5 < EnemyMaxHP {
        return true
    }
    return false
}

func checkEnemyNotBloodDrained(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    enemy, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }

    if !enemy.IsFlag(EFFECT_CLASS_BLOOD_DRAIN) {
        return true
    }
    return false
}

func checkEnemyNotGreenPoisoned(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    vampire, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }
    if !vampire.IsFlag(EFFECT_CLASS_POISON) {
        return true
    }
    return false
}

func checkEnemyNotYellowPoisoned(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    vampire, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }

    if !vampire.IsFlag(EFFECT_CLASS_YELLOW_POISON_TO_CREATURE) {
        return true
    }
    return false
}

func checkEnemyNotDarkbluePoisoned(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    vampire, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }
    if !vampire.IsFlag(EFFECT_CLASS_DARKBLUE_POISON) {
        return true
    }
    return false
}

func checkEnemyNotGreenStalkered(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    vampire, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }
    if !vampire.IsFlag(EFFECT_CLASS_GREEN_STALKER) {
        return true
    }
    return false
}

func checkEnemyNotParalyzed(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    if !pEnemy.CreatureInstance().isFlag(EFFECT_CLASS_PARALYZE) {
        return true
    }
    return false
}

func checkEnemyNotDoomed(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    if !pEnemy.CreatureInstance().isFlag(EFFECT_CLASS_DOOM) {
        return true
    }
    return false
}

func checkEnemyNotBlinded(pMonster *Monster, pEnemy CreatureInterface) bool {
    return false
}

func checkEnemyNotInDarkness(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    vampire, ok := pEnemy.(*Vampire)
    if !ok {
        return false
    }
    if !vampire.isFlag(EFFECT_CLASS_DARKNESS) {
        return true
    }
    return false
}

func checkEnemyNotSeduction(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    if !pEnemy.CreatureInstance().isFlag(EFFECT_CLASS_SEDUCTION) {
        return true
    }
    return false
}

func checkEnemyNotDeath(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }
    if !pEnemy.CreatureInstance().isFlag(EFFECT_CLASS_DEATH) {
        return true
    }
    return false
}

func checkImOK(pMonster *Monster, pEnemy CreatureInterface) bool {
    CurHP := pMonster.HP[ATTR_CURRENT]
    MaxHP := pMonster.HP[ATTR_MAX]
    if CurHP*3 > MaxHP {
        return true
    }
    return false
}

func checkImDying(pMonster *Monster, pEnemy CreatureInterface) bool {
    CurHP := pMonster.HP[ATTR_CURRENT]
    MaxHP := pMonster.HP[ATTR_MAX]
    if CurHP*4 < MaxHP {
        return true
    }
    return false
}

func checkImDamaged(pMonster *Monster, pEnemy CreatureInterface) bool {
    return pMonster.Brain.getDamaged()
}

func checkImHiding(pMonster *Monster, pEnemy CreatureInterface) bool {

    return pMonster.CreatureInstance().isFlag(EFFECT_CLASS_HIDE)
}

func checkImWolf(pMonster *Monster, pEnemy CreatureInterface) bool {

    return pMonster.CreatureInstance().isFlag(EFFECT_CLASS_TRANSFORM_TO_WOLF)
}

func checkImBat(pMonster *Monster, pEnemy CreatureInterface) bool {

    return pMonster.CreatureInstance().isFlag(EFFECT_CLASS_TRANSFORM_TO_BAT)
}

func checkImInvisible(pMonster *Monster, pEnemy CreatureInterface) bool {

    return pMonster.CreatureInstance().isFlag(EFFECT_CLASS_INVISIBILITY)
}

func checkImWalkingWall(pMonster *Monster, pEnemy CreatureInterface) bool {

    pAI := pMonster.Brain
    if pAI.getWallCount() > 3 &&
        (pAI.MoveRule == MOVE_RULE_LEFTWALL || pAI.MoveRule == MOVE_RULE_RIGHTWALL) {
        return true
    }

    return false
}

func checkTimingBloodDrain(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil || pEnemy.CreatureClass() == CREATURE_CLASS_NPC {
        return false
    }

    enemy := pEnemy.CreatureInstance()
    if enemy.isFlag(EFFECT_CLASS_BLOOD_DRAIN) ||
        enemy.isFlag(EFFECT_CLASS_NO_DAMAGE) ||
        enemy.isFlag(EFFECT_CLASS_IMMUNE_TO_BLOOD_DRAIN) {
        return false
    }

    dist := pMonster.getDistance(enemy.X, enemy.Y)
    if dist > pMonster.getMeleeRange() {
        return false
    }

    // EnemyCurHP := enemy.HP[ATTR_CURRENT]
    // EnemyMaxHP := enemy.HP[ATTR_CURRENT]
    //
    // if EnemyCurHP*5 >= EnemyMaxHP {
    //     return false
    // }

    return true

}

func checkTimingMasterBloodDrain(pMonster *Monster, pEnemy CreatureInterface) bool {

    if pEnemy == nil {
        return false
    }

    // masterHPPercent := pMonster.HP[ATTR_CURRENT] * 100 / pMonster.HP[ATTR_MAX]

    // startHPPercent := g_pVariableManager.getMasterBloodDrainStartHP()
    // startBDPercent := g_pVariableManager.getMasterBloodDrainStartBD()
    // endHPPercent := g_pVariableManager.getMasterBloodDrainEndHP()
    // endBDPercent := g_pVariableManager.getMasterBloodDrainEndBD()
    //
    // if masterHPPercent >= startHPPercent {
    //     return false
    // }

    // ratio := rand.Intn(100)

    // if masterHPPercent <= endHPPercent {
    //     return ratio < endBDPercent
    // }

    // maxBDPercent := max(startBDPercent, endBDPercent)
    // gapHPPercent := startHPPercent - endHPPercent
    // gapBDPercent := abs(endBDPercent - startBDPercent)
    //
    // permitRatio = maxBDPercent -
    // gapBDPercent*(masterHPPercent-endHPPercent)/gapHPPercent
    // return ratio < permitRatio

    // TODO
    return false
}

func checkMasterSummonTiming(pMonster *Monster, pEnemy CreatureInterface) bool {
    if !pMonster.isMaster {
        return false
    }

    pZone := pMonster.getZone()

    if !pZone.isMasterLair() {
        return false
    }

    // pMasterLairManager := pZone.getMasterLairManager()
    // bSummonTiming := !pMasterLairManager.isMasterReady() && pZone.getMonsterManager().getSize() <= 1

    // return bSummonTiming
    // TODO
    return false
}

//----------------------------------------------------------------------
//
// bool checkMasterNotReady(pMonster *Monster, pEnemy CreatureInterface) bool
//
//----------------------------------------------------------------------
func checkMasterNotReady(pMonster *Monster, pEnemy CreatureInterface) bool {

    if !pMonster.isMaster {
        return false
    }

    pZone := pMonster.getZone()

    if !pZone.isMasterLair() {
        return false
    }

    // pMasterLairManager := pZone.getMasterLairManager()
    // return !pMasterLairManager.isMasterReady()
    // TODO
    return false
}

//----------------------------------------------------------------------
//
// bool checkImInBadPosition(pMonster *Monster, pEnemy CreatureInterface) bool
//
//----------------------------------------------------------------------
// ÇöÀç À§Ä¡°¡ ¾È ÁÁÀº °÷ÀÎ°¡?
//
// pMonster°¡ ÀÖ´Â Å¸ÀÏ¿¡ AcidSwampÀÌ »Ñ·ÁÁ® ÀÖ´Â °æ¿ì
//----------------------------------------------------------------------
func checkImInBadPosition(pMonster *Monster, pEnemy CreatureInterface) bool {
    pZone := pMonster.getZone()

    // enemy := pEnemy.CreatureInstance()
    rTile := pZone.getTile(pMonster.X, pMonster.Y)

    if !pMonster.isFlag(EFFECT_CLASS_NO_DAMAGE) &&
        !pMonster.isFlag(EFFECT_CLASS_IMMUNE_TO_ACID) &&
        rTile.getEffect(EFFECT_CLASS_ACID_SWAMP) != nil {
        return true
    }

    if rTile.getEffect(EFFECT_CLASS_BLOODY_WALL) != nil ||
        rTile.getEffect(EFFECT_CLASS_GROUND_ATTACK) != nil {
        return true
    }

    return false
}

func checkFindWeakEnemy(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    pZone := pMonster.getZone()

    // strongValue := getStrongValue(pEnemy)
    //
    // if strongValue == 0 {
    //     return false
    // }

    var pWeakestCreature CreatureInterface

    cx := pMonster.X
    cy := pMonster.Y
    var (
        ix   int
        iy   int
        endx int
        endy int
    )

    sight := pMonster.getSight()

    endx = min(int(pZone.getWidth()-1), int(cx+ZoneCoord_t(sight)+1))
    endy = min(int(pZone.getHeight()-1), int(cy+ZoneCoord_t(sight)+1))

    for ix = max(0, int(cx-ZoneCoord_t(sight)-1)); ix <= endx; ix++ {
        for iy = max(0, int(cy-ZoneCoord_t(sight)-1)); iy <= endy; iy++ {
            rTile := pZone.Tile(ix, iy)

            if rTile.hasCreature() {
                // const list<Object*> & objectList = rTile.getObjectList();
                // for (list<Object*>::const_iterator itr = objectList.begin() ;
                // 	itr != objectList.end() && (*itr).getObjectPriority() <= OBJECT_PRIORITY_BURROWING_CREATURE;
                // 	itr++)
                // {
                // 	Creature* pCreature = dynamic_cast<Creature*>(*itr);
                // 	Assert(pCreature != nil);
                //
                // 	// pMonster, pEnemy°¡ ¾Æ´Ï¶ó¸é..
                // 	// Player¶ó¸é °ø°ÝÇÏ°Ô µÈ´Ù.
                // 	if (pCreature != pMonster
                // 		&& pCreature != pEnemy
                // 		&& pCreature.isPC()
                // 		&& pMonster.isRealEnemy(pCreature))
                // 	{
                // 		int checkStrongValue = getStrongValue(pCreature);
                //
                // 		// ´õ ¾àÇÑ³ÑÀ» pWeakestCreature·Î ÀÓ¸í~ÇÑ´Ù
                // 		if (checkStrongValue < strongValue)
                // 		{
                // 			pWeakestCreature = pCreature;
                // 			strongValue = checkStrongValue;
                // 		}
                // 	}
                // }
            }
        }
    }

    // Á© ¾àÇÑ³ÑÀ» Ã£Àº °æ¿ì..
    if pWeakestCreature != nil {
        pMonster.deleteAllEnemy()
        pMonster.addEnemy(pWeakestCreature)

        return true
    }

    return false
}

func checkEnemyNotHallucination(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }

    if !pEnemy.CreatureInstance().isFlag(EFFECT_CLASS_HALLUCINATION) {
        return true
    }
    return false
}

func checkTimingDuplicateSelf(pMonster *Monster, pEnemy CreatureInterface) bool {
    if !pMonster.isMaster {
        return false
    }

    // pZone := pMonster.getZone()

    currentHP := pMonster.HP[ATTR_CURRENT]
    maxHP := pMonster.HP[ATTR_MAX]

    if currentHP*100/maxHP > 70 {
        return false
    }

    // if pZone.getMonsterManager().getSize() > 12 {
    //     return false
    // }

    return true
}

func checkPossibleSummonMonsters(pMonster *Monster, pEnemy CreatureInterface) bool {
    return pMonster.hasNextMonsterSummonInfo()
}

func checkEnemyTileNotAcidSwamp(pMonster *Monster, pEnemy CreatureInterface) bool {
    pZone := pMonster.getZone()

    rTile := pZone.getTile(pEnemy.CreatureInstance().X, pEnemy.CreatureInstance().Y)
    if rTile.getEffect(EFFECT_CLASS_ACID_SWAMP) == nil {
        return true
    }

    return false
}

func checkEnemyOnAir(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    if pEnemy.CreatureInstance().MoveMode == MOVE_MODE_FLYING {
        return true
    }
    return false
}

func checkEnemyOnSafeZone(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }
    return (pMonster.getZone().getZoneLevel(pEnemy.CreatureInstance().X, pEnemy.CreatureInstance().Y) & ZoneLevel_t(SAFE_ZONE)) != 0
}

func checkCanAttackThrowingAxe(pMonster *Monster, pEnemy CreatureInterface) bool {
    if pEnemy == nil {
        return false
    }

    // enemy := pEnemy.CreatureInstance()
    // dir := getDirection(pMonster.getX(), pMonster.getY(), enemy.X, enemy.Y)
    // X := ZoneCoord_t(int(pMonster.getX()) + dirMoveMask[dir].X*7)
    // Y := ZoneCoord_t(int(pMonster.getY()) + dirMoveMask[dir].Y*7)
    // distance := getDistance(X, Y, enemy.X, enemy.Y)

    // if distance > 2 {
    //     return false
    // }
    return true
}

func getDirection(x1, y1, x2, y2 ZoneCoord_t) Dir_t {
    // TODO
    return DIR_NONE
}
