package cell

import (
    "fmt"
    "github.com/tiancaiamao/ouster/aoi"
    "testing"
)

func Test(t *testing.T) {
    cellaoi := New(256, 256, 32, 32)
    o1 := cellaoi.Add(0, 0, 1)
    o2 := cellaoi.Add(100, 0, 2)
    o3 := cellaoi.Add(0, 100, 3)
    cellaoi.Add(100, 100, 4)

    t.Log(o1.Id(), o2.Id(), o3.Id())

    for i := 0; i < 100; i++ {
        cellaoi.Update(o1, o1.X()+1, o1.Y()+1)
        cellaoi.Update(o2, o2.X()-1, o2.Y()+1)
        cellaoi.Update(o3, o3.X()+1, o3.Y()-1)

        t.Logf("---------the %dth round-----------", i)
        cellaoi.Message(func(watcher aoi.Entity, marker aoi.Entity) {
            t.Logf("watcher %d (%d %d)=> marker %d (%d %d)",
                watcher.Id(),
                watcher.X(),
                watcher.Y(),
                marker.Id(),
                marker.X(),
                marker.Y())
        })
    }
    fmt.Println("finished")
}
