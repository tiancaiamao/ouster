package cell

import "github.com/tiancaiamao/ouster/aoi"

type Entity struct {
    next *Entity
    prev *Entity

    id      uint32
    x       uint8
    y       uint8
    version uint8
}

func (e *Entity) X() uint8 {
    return e.x
}

func (e *Entity) Y() uint8 {
    return e.y
}

func (e *Entity) Id() uint32 {
    return e.id
}

func (e *Entity) remove() {
    e.prev.next = e.next
    if e.next != nil {
        e.next.prev = e.prev
    }
}

type Sector struct {
    head   Entity
    nearby [8]*Sector
}

func (s *Sector) insert(e *Entity) {
    e.next = s.head.next
    if s.head.next != nil {
        s.head.next.prev = e
    }
    s.head.next = e
    e.prev = &s.head
}

type CellAoi struct {
    mapWidth   uint16
    mapHeight  uint16
    cellWidth  uint8
    cellHeight uint8
    cells      [][]Sector
    version    uint8
}

var idx2dir [8][2]int = [8][2]int{
    [2]int{-1, -1},
    [2]int{0, -1},
    [2]int{1, -1},
    [2]int{1, 0},
    [2]int{1, 1},
    [2]int{0, 1},
    [2]int{-1, 1},
    [2]int{-1, 0},
}

func New(mapWidth, mapHeight uint16, cellWidth, cellHeight uint8) *CellAoi {
    ret := &CellAoi{
        mapWidth:   mapWidth,
        mapHeight:  mapHeight,
        cellWidth:  cellWidth,
        cellHeight: cellHeight,
    }

    width := (mapWidth+uint16(cellWidth))/uint16(cellWidth) - 1
    height := (mapHeight+uint16(cellHeight))/uint16(cellHeight) - 1
    ret.cells = make([][]Sector, width)
    for i := 0; i < int(width); i++ {
        ret.cells[i] = make([]Sector, height)
    }

    for x := 0; x < int(width); x++ {
        for y := 0; y < int(height); y++ {
            sector := &ret.cells[x][y]

            for i := 0; i < 8; i++ {
                dir := idx2dir[i]
                xx := x + dir[0]
                yy := y + dir[1]

                if xx >= 0 && xx < int(width) && yy >= 0 && yy < int(height) {
                    sector.nearby[i] = &ret.cells[xx][yy]
                }
            }
        }
    }

    return ret
}

func (aoi *CellAoi) getCell(x uint8, y uint8) *Sector {
    xIndex := x / aoi.cellWidth
    yIndex := y / aoi.cellHeight

    return &aoi.cells[xIndex][yIndex]
}

func (aoi *CellAoi) Add(x uint8, y uint8, id uint32) aoi.Entity {
    sector := aoi.getCell(x, y)
    if sector == nil {
        return nil
    }

    // TODO: object pool to reuse memory
    ret := new(Entity)
    ret.x = x
    ret.y = y
    ret.id = id
    sector.insert(ret)
    return ret
}

func (aoi *CellAoi) Del(entity aoi.Entity) {
    e := entity.(*Entity)
    e.remove()
}

func (aoi *CellAoi) Message(cb aoi.Callback) {
    aoi.version++

    width := (aoi.mapWidth+uint16(aoi.cellWidth))/uint16(aoi.cellWidth) - 1
    height := (aoi.mapHeight+uint16(aoi.cellHeight))/uint16(aoi.cellHeight) - 1
    for x := 0; x < int(width); x++ {
        for y := 0; y < int(height); y++ {
            sector := &aoi.cells[x][y]

            for e := sector.head.next; e != nil; e = e.next {
                if e.version == aoi.version {
                    for ee := sector.head.next; ee != nil; ee = ee.next {
                        if e != ee {
                            cb(e, ee)
                        }
                    }
                } else {
                    e.version = aoi.version
                }
            }
        }
    }
}

func (aoi *CellAoi) Update(entity aoi.Entity, x uint8, y uint8) {
    e := entity.(*Entity)
    e.version++
    oldSector := aoi.getCell(e.x, e.y)
    newSector := aoi.getCell(x, y)

    // remove from oldSector and insert into newSector
    if oldSector != newSector {
        e.remove()
        newSector.insert(e)
    }

    e.x = x
    e.y = y
}

func (aoi *CellAoi) Nearby(x uint8, y uint8, cb aoi.Callback) {
    sector := aoi.getCell(x, y)

    for e := sector.head.next; e != nil; e = e.next {
        cb(nil, e)
    }
}
