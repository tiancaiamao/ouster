package aoi

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	fmt.Println("test begin...")
	aoi := New()
	var objs [4]FPoint
	objs[0].X = 0
	objs[0].Y = 0
	objs[1].X = 100
	objs[1].Y = 0
	objs[2].X = 0
	objs[2].Y = 100
	objs[3].X = 100
	objs[3].Y = 100

	for i := 0; i < 100; i++ {
		objs[0].X += 1
		objs[0].Y += 1
		objs[1].X -= 1
		objs[1].Y += 1
		objs[2].X += 1
		objs[2].Y -= 1
		// objs[3].Y -= 1
		// objs[3].Y -= 1

		for i := uint32(0); i < 4; i++ {
			if i == 1 {
				aoi.Update(i, ModeMarker, objs[i])
			} else {
				aoi.Update(i, ModeMarker|ModeWatcher, objs[i])
			}
		}

		aoi.Message(func(watcher uint32, marker uint32) {
			fmt.Println("run here...")
			t.Logf("watcher %d => marker %d", watcher, marker)
		})
	}
	fmt.Println("finished")
}
