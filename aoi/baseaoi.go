package aoi

type AoiObj struct {
	Point
	AoiMode
}

type BaseAoi struct {
}

func (aoi *BaseAoi) Update(id uint32, mode AoiMode, pos Point) {
	
}

func (aoi *BaseAoi) Message(cb AoiCallback) {
}
