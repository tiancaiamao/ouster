package scene

import (
	"github.com/tiancaiamao/ouster/data"
)

var (
	maps map[string]*Map
)

func Initialize() {
	_maps := []*data.Map{
		&data.AncientTemple,
		&data.FrontierOutpost,
	}
	_names := []string{
		"ancient_temple",
		"frontier_outpost",
	}

	maps = make(map[string]*Map)
	for i:=0; i<len(_maps); i++ {
		name := _names[i]
		mapData := _maps[i]
		m := New(mapData)		
		m.Go()
		
		maps[name] = m
	}
}

func Query(mapName string) *Map {
	return maps[mapName]
}
