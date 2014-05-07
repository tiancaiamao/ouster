package scene

import (
	"github.com/tiancaiamao/ouster/data/darkeden"
	"os"
)

var (
	maps map[string]*Zone
)

func Initialize() {
	dir := "/Users/genius/project/vs/data/"
	smp, _ := os.Open(dir + "limbo_lair_se.smp")
	mapData, _ := darkeden.Load(smp)

	m := New(mapData)
	maps["test"] = m
	m.Go()
}

func Query(mapName string) *Zone {
	return maps[mapName]
}
