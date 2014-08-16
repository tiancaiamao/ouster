package main

type CorpseType uint8

const (
    SLAYER_CORPSE = iota
    VAMPIRE_CORPSE
    NPC_CORPSE
    MONSTER_CORPSE
    OUSTERS_CORPSE
)

type Corpse struct {
    Item

    Treasures     []*Item
    TreasureCount uint8

    X   ZoneCoord_t
    Y   ZoneCoord_t
}
