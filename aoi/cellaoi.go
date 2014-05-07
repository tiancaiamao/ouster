package aoi

import ()

type Entity struct {
	next *Entity
	prev *Entity

	x uint16
	y uint16
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
	cellWidth  uint16
	cellHeight uint16
	cells      [][]Sector
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

func NewCellAoi(mapWidth, mapHeight, cellWidth, cellHeight uint16) *CellAoi {
	ret := &CellAoi{
		mapWidth:   mapWidth,
		mapHeight:  mapHeight,
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
	}

	width := (mapWidth+cellWidth)/cellWidth - 1
	height := (mapHeight+cellHeight)/cellHeight - 1
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

	// for

	// ret.id2Pos = make(map[uint32]Point)
	return ret
}

func (aoi *CellAoi) getCell(x uint16, y uint16) *Sector {
	if x >= aoi.mapWidth || y >= aoi.mapHeight {
		return nil
	}

	xIndex := x / aoi.cellWidth
	yIndex := y / aoi.cellHeight

	return &aoi.cells[xIndex][yIndex]
}

func (aoi *CellAoi) Add(x uint16, y uint16) *Entity {
	sector := aoi.getCell(x, y)
	if sector == nil {
		return nil
	}

	// TODO: object pool to reuse memory
	ret := new(Entity)
	ret.x = x
	ret.y = y
	sector.insert(ret)
	return ret
}

func (aoi *CellAoi) Update(e *Entity, x uint16, y uint16) {
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

type Callback func(*Entity)

func (aoi *CellAoi) Nearby(x uint16, y uint16, cb Callback) {
	sector := aoi.getCell(x, y)

	for e := sector.head.next; e != nil; e = e.next {
		cb(e)
	}
}
