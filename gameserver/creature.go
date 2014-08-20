package main

type CreatureClass int

const (
    CREATURE_CLASS_SLAYER  CreatureClass = iota // PC Slayer
    CREATURE_CLASS_VAMPIRE                      // PC Vampire
    CREATURE_CLASS_NPC                          // NPC
    CREATURE_CLASS_MONSTER                      // NPC Slayer, NPC Vampire
    CREATURE_CLASS_OUSTER                       // PC Ousters
    CREATURE_CLASS_MAX
)

type MoveMode uint8

const (
    MOVE_MODE_WALKING MoveMode = iota
    MOVE_MODE_FLYING
    MOVE_MODE_BURROWING
    MOVE_MODE_MAX
)

type CreatureInterface interface {
    ObjectInterface
    CreatureClass() CreatureClass
    CreatureInstance() *Creature
    IsAbleToMove() bool
}

type Creature struct {
    Object
    MoveMode MoveMode
    X        ZoneCoord_t
    Y        ZoneCoord_t
    Dir      Dir_t

    ViewportWidth       ZoneCoord_t
    ViewportUpperHeight ZoneCoord_t
    ViewportLowerHeight ZoneCoord_t
    Resist              [MAGIC_DOMAIN_MAX]Resist_t

    Flag BitSet

    Sight Sight_t
}

func (c Creature) ObjectClass() ObjectClass {
    return OBJECT_CLASS_CREATURE
}

func (c Creature) CreatureInstance() *Creature {
    return &c
}

func (c Creature) heartbeat() {
	
}

// TODO
func (c Creature) IsAbleToMove() bool {
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

func canSee(watcher CreatureInterface, marker CreatureInterface) bool {
    return true //TODO
}
