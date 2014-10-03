package main

import (
    . "github.com/tiancaiamao/ouster/util"
)

type CreatureClass int

const (
    CREATURE_CLASS_SLAYER  CreatureClass = iota // PC Slayer
    CREATURE_CLASS_VAMPIRE                      // PC Vampire
    CREATURE_CLASS_NPC                          // NPC
    CREATURE_CLASS_MONSTER                      // NPC Slayer, NPC Vampire
    CREATURE_CLASS_OUSTER                       // PC Ousters
    CREATURE_CLASS_MAX
)

type CreatureInterface interface {
    ObjectInterface
    CreatureClass() CreatureClass
    CreatureInstance() *Creature

    getProtection() Protection_t
    getHP(int) HP_t
    IsAbleToMove() bool
}

type Creature struct {
    Object

    Scene    *Scene
    MoveMode MoveMode
    X        ZoneCoord_t
    Y        ZoneCoord_t
    Dir      Dir_t

    Level Level_t

    ViewportWidth       ZoneCoord_t
    ViewportUpperHeight ZoneCoord_t
    ViewportLowerHeight ZoneCoord_t
    Resist              [MAGIC_DOMAIN_MAX]Resist_t

    Flag BitSet

    Sight Sight_t
}

func (c *Creature) Init() {
    c.Flag = NewBitSet(EFFECT_CLASS_MAX)
}

func (c *Creature) ObjectClass() ObjectClass {
    return OBJECT_CLASS_CREATURE
}

func (c *Creature) CreatureInstance() *Creature {
    return c
}

func (c *Creature) heartbeat() {

}

// TODO
func (c *Creature) IsAbleToMove() bool {
    if c.Flag.IsFlag(EFFECT_CLASS_COMA) ||
        c.Flag.IsFlag(EFFECT_CLASS_PARALYZE) ||
        c.Flag.IsFlag(EFFECT_CLASS_ETERNITY_PAUSE) ||
        c.Flag.IsFlag(EFFECT_CLASS_CASKET) ||
        c.Flag.IsFlag(EFFECT_CLASS_CAUSE_CRITICAL_WOUNDS) ||
        c.Flag.IsFlag(EFFECT_CLASS_SOUL_CHAIN) ||
        c.Flag.IsFlag(EFFECT_CLASS_LOVE_CHAIN) ||
        c.Flag.IsFlag(EFFECT_CLASS_GUN_SHOT_GUIDANCE_AIM) ||
        c.Flag.IsFlag(EFFECT_CLASS_SLEEP) ||
        c.Flag.IsFlag(EFFECT_CLASS_ARMAGEDDON) ||
        c.Flag.IsFlag(EFFECT_CLASS_POISON_MESH) ||
        c.Flag.IsFlag(EFFECT_CLASS_TENDRIL) ||
        c.Flag.IsFlag(EFFECT_CLASS_TRAPPED) ||
        c.Flag.IsFlag(EFFECT_CLASS_INSTALL_TURRET) ||
        c.Flag.IsFlag(EFFECT_CLASS_EXPLOSION_WATER) {
        return false
    }

    return true
}

func (c *Creature) IsFlag(effect uint) bool {
    return c.Flag.IsFlag(effect)
}

func (c *Creature) isFlag(effect uint) bool {
    return c.Flag.IsFlag(effect)
}

// TODO
func (c *Creature) removeFlag(effect EffectClass) {

}

func (c *Creature) setFlag(ec EffectClass) {

}

func canSee(watcher CreatureInterface, marker CreatureInterface) bool {
    return true //TODO
}

func (c *Creature) canMove(nx ZoneCoord_t, ny ZoneCoord_t) bool {
    if c.Flag.IsFlag(EFFECT_CLASS_POISON_MESH) ||
        c.Flag.IsFlag(EFFECT_CLASS_TENDRIL) ||
        c.Flag.IsFlag(EFFECT_CLASS_BLOODY_WALL_BLOCKED) ||
        c.Flag.IsFlag(EFFECT_CLASS_CASKET) ||
        !isValidZoneCoord(&c.Scene.Zone, nx, ny) {
        return false
    }

    rTile := c.Scene.getTile(nx, ny)

    if rTile.isBlocked(c.MoveMode) ||
        rTile.hasEffect() && (rTile.getEffect(EFFECT_CLASS_BLOODY_WALL_BLOCKED) != nil ||
            rTile.getEffect(EFFECT_CLASS_SANCTUARY) != nil) {
        return false
    }

    rNewTile := c.Scene.getTile(c.X, c.Y)

    if rNewTile.getEffect(EFFECT_CLASS_SANCTUARY) != nil {
        return false
    }

    return true
}

func (c *Creature) isBlockedByCreature(nx ZoneCoord_t, ny ZoneCoord_t) bool {
    if !isValidZoneCoord(&c.Scene.Zone, nx, ny) ||
        !c.Scene.getTile(nx, ny).HasCreature(c.MoveMode) {
        return false
    }
    return true
}

func (c *Creature) addEffect(effect EffectInterface) {

}

// TODO
func verifyDistance(c1, c2 CreatureInterface) bool {
    return true
}