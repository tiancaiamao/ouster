package aoi

type Entity interface {
	X() uint8
	Y() uint8
	Id() uint32
}

type Callback func(watcher Entity, marker Entity)
type Aoi interface {
	Add(x uint8, y uint8, id uint32) Entity
	Del(e Entity)
	Update(e Entity, x uint8, y uint8)
	Message(cb Callback)
}
