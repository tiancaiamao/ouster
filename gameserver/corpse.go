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

// Corpse继承Item对象，实现ItemInterface接口
func (c Corpse) ItemClass() ItemClass {
    return ITEM_CLASS_CORPSE
}

type SlayerCorpse struct {
    Corpse
}

type OusterCorpse struct {
    Corpse
}

type VampireCorpse struct {
    Corpse
}

type MonsterCorpse struct {
	Corpse
}