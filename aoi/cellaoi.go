package aoi

type Point struct {
	X uint32
	Y uint32
}

type CellAoi struct {
	mapWidth   uint32
	mapHeight  uint32
	cellWidth  uint32
	cellHeight uint32
	cells      [][]uint32
	id2Pos     map[uint32]Point
}

func NewCellAoi(mapWidth uint32, mapHeight, uint32, cellWidth uint32, cellHeight uint32) *CellAoi {
	ret := &CellAoi{
		mapWidth:   mapWidth,
		mapHeight:  mapHeight,
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
	}
	// width := (mapWidth+cellWidth)/cellWidth - 1
	// height := (mapHeight+cellHeight)/cellHeight - 1
	// ret.cells = make([][]uint32, width * height)
	// ret.id2Pos = make(map[uint32]Point)
	return ret
}

// func getCell(x uint32, y uint32) []uint32 {
// 	return
// }

func (aoi *CellAoi) Update(id uint32, mode AoiMode, pos Point) {
	// prevPos, ok := id2Pos[id]
	// if !ok {

	// } else {
	// 	old := getCell(prevPos)

	// 	for nearby := range old {

	// 	}

	// 	in := getCell(pos.X, pos.Y)
	// 	for nearby := range in {

	// 	}

	// }
}

func (aoi *CellAoi) Message(cb AoiCallback) {

}
