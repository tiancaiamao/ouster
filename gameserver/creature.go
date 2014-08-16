package main

type CreatureClass uint8

const (
    CREATURE_CLASS_SLAYER  = iota // PC Slayer
    CREATURE_CLASS_VAMPIRE        // PC Vampire
    CREATURE_CLASS_NPC            // NPC
    CREATURE_CLASS_MONSTER        // NPC Slayer, NPC Vampire
    CREATURE_CLASS_OUSTERS        // PC Ousters
    CREATURE_CLASS_MAX
)

type Creature struct {
    Object

    ViewportWidth       ZoneCoord_t
    ViewportUpperHeight ZoneCoord_t
    ViewportLowerHeight ZoneCoord_t
    Resist              [MAGIC_DOMAIN_MAX]Resist_t
    Dir                 Dir_t
    Sight               Sight_t
    CClass              CreatureClass
}
