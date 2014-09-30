package main

import (
    "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
)

type TileFlags uint16

const (
    TILE_GROUND_BLOCKED = iota
    TILE_AIR_BLOCKED
    TILE_UNDERGROUND_BLOCKED
    TILE_WALKING_CREATURE
    TILE_FLYING_CREATURE
    TILE_BURROWING_CREATURE
    TILE_ITEM
    TILE_OBSTACLE
    TILE_EFFECT
    TILE_BUILDING
    TILE_PORTAL
    TILE_TERRAIN
)

type Tile struct {
    Flags   uint16
    Option  uint16
    Objects map[ObjectID_t]ObjectInterface
    Sector  *Sector
}

const SECTOR_SIZE = 13

type Sector struct {
    Objects       map[ObjectID_t]Object
    NearbySectors [9]*Sector
}

func (tile *Tile) SetBlocked(mode MoveMode) {
    tile.Flags |= (1 << mode)
}

// TODO
func (tile *Tile) HasPortal() bool {
    return false
}

func (tile *Tile) getEffect(uint16) EffectInterface {
    // TODO
    return nil
}

// TODO
func (time *Tile) deleteEffect(ObjectID_t) {
}

// TODO
func (tile *Tile) hasEffect() bool {
    return false
}

func (tile *Tile) IsGroundBlocked() bool {
    return (tile.Flags & (1 << TILE_GROUND_BLOCKED)) != 0
}

func (tile *Tile) IsAirBlocked() bool {
    return (tile.Flags & (1 << TILE_AIR_BLOCKED)) != 0
}

func (tile *Tile) IsUndergroundBlocked() bool {
    return (tile.Flags & (1 << TILE_UNDERGROUND_BLOCKED)) != 0
}

func (tile *Tile) IsBlocked(m MoveMode) bool {
    return (tile.Flags & (1 << m)) != 0
}

func (tile *Tile) isBlocked(m MoveMode) bool {
    return (tile.Flags & (1 << m)) != 0
}

func (tile *Tile) HasCreature(m MoveMode) bool {
    return (tile.Flags & (1 << (TileFlags(m) + TILE_WALKING_CREATURE))) != 0
}

func (tile *Tile) hasCreature() bool {
    return tile.Flags != 0
}

// TODO
func (tile *Tile) hasPortal() bool {
    return false
}

func (tile *Tile) DeleteCreature(id ObjectID_t) {
    var object ObjectInterface
    object, ok := tile.Objects[id]
    if !ok {
        log.Warnf("object not exist in tile! id=%d tile=%#v\n", id, tile)
        return
    }

    delete(tile.Objects, id)

    var creature *Creature
    switch raw := object.(type) {
    case *Agent:
        creature = raw.CreatureInstance()
    case *Monster:
        creature = raw.CreatureInstance()
    default:
        log.Errorln(raw)
        panic("不对")
    }
    tile.Flags &^= (1 << (TILE_WALKING_CREATURE + creature.MoveMode))
    tile.Flags &^= (1 << (TILE_GROUND_BLOCKED + creature.MoveMode))
}

func (tile *Tile) AddCreature(creature CreatureInterface) {
    inst := creature.CreatureInstance()
    if tile.HasCreature(inst.MoveMode) {
        panic("重复加入到tile")
    }

    tile.Objects[inst.ObjectID] = creature

    tile.Flags |= (1 << (TILE_WALKING_CREATURE + inst.MoveMode))
    tile.Flags |= (1 << (TILE_GROUND_BLOCKED + inst.MoveMode))
}

func (tile *Tile) GetCreature(mode MoveMode) CreatureInterface {
    // TODO
    return nil
}

func (tile *Tile) AddEffect(effect EffectInterface) {
    // TODO
}
