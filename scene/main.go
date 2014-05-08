package scene

import (
	"github.com/tiancaiamao/ouster/data"
)

var (
	maps map[string]*Zone
)

func Initialize() {
	maps = make(map[string]*Zone)
	maps["limbo_lair_se"] = New(data.LimboLairSE)

	for _, m := range maps {
		m.Go()
	}
}

func Query(mapName string) *Zone {
	return maps[mapName]
}
