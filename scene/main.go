package scene

import (
	"github.com/tiancaiamao/ouster/data"
)

var (
	maps *Map
)

func Initialize() {
	mapData := &data.Test
	m := New(mapData)
	maps = m
	m.Go()
}

func Query(mapName string) *Map{
	return maps
}
