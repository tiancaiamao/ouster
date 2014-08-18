package main

type TileFlags uint16

const (
    TILE_GROUND_BLOCKED TileFlags = iota
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
    Objects []ObjectInterface
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

func (tile *Tile) HasCreature(m MoveMode) bool {
    return (tile.Flags & (1 << (TileFlags(m) + TILE_WALKING_CREATURE))) != 0
}

func (tile *Tile) DeleteCreature(id ObjectID_t) {
    for i := 0; i < len(tile.Objects); i++ {
        if tile.Objects[i].ObjectInstance().ObjectID == id {
            copy(tile.Objects[i:], tile.Objects[i+1:])
            return
        }
    }
}

func (tile *Tile) AddCreature(creature CreatureInterface) {
    tile.Objects = append(tile.Objects, creature)
}
