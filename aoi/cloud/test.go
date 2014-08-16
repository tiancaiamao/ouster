package cloud

import "fmt"

type AoiMode uint8

const (
    ModeWatcher = (1 << iota)
    ModeMarker
    ModeDrop
)

type MyAoi struct {
    all           map[uint32]*object
    watcherStatic map[uint32]*object
    markerStatic  map[uint32]*object
    watcherMove   map[uint32]*object
    markerMove    map[uint32]*object
    // moveQueue     []uint32
    version uint32
}

func NewMy() *MyAoi {
    aoi := new(MyAoi)
    aoi.all = make(map[uint32]*object)
    aoi.watcherStatic = make(map[uint32]*object)
    aoi.markerStatic = make(map[uint32]*object)
    aoi.watcherMove = make(map[uint32]*object)
    aoi.markerMove = make(map[uint32]*object)
    // aoi.moveQueue = make([]uint32, 0, 100)
    return aoi
}

const RADIS float32 = 20.0

func (aoi *MyAoi) Update(id uint32, mode AoiMode, pos FPoint) {
    if (mode & ModeDrop) != 0 {
        delete(aoi.all, id)
        delete(aoi.watcherStatic, id)
        delete(aoi.watcherMove, id)
        delete(aoi.markerStatic, id)
        delete(aoi.markerMove, id)
        return
    }

    obj, ok := aoi.all[id]
    if !ok {
        obj = new(object)
        aoi.all[id] = obj
    }

    moved := posChanged(obj, pos)
    changed := modeChanged(obj, mode)
    if changed {
        obj.mode = mode
        moved = true
    }

    if moved {
        fmt.Printf("moved: %d, pos=%f, %f\n", id, obj.pos.X, obj.pos.Y)
        obj.pos = pos
        obj.last = pos
        obj.version = aoi.version

        // aoi.moveQueue = append(aoi.moveQueue, id)

        if (obj.mode & ModeWatcher) != 0 {
            delete(aoi.watcherStatic, id)
            aoi.watcherMove[id] = obj
        }
        if (obj.mode & ModeMarker) != 0 {
            delete(aoi.markerStatic, id)
            aoi.markerMove[id] = obj
        }
    }
}

func modeChanged(obj *object, mode AoiMode) bool {
    var changed bool
    if obj.mode != mode {
        changed = true
    }
    return changed
}

func (aoi *MyAoi) Message(cb AoiCallback) {
    fmt.Println("enter Message-------------")
    fmt.Println("watchMove:", aoi.watcherMove)
    fmt.Println("markerMove", aoi.markerMove)
    fmt.Println("watcherStatic", aoi.watcherStatic)
    fmt.Println("markerStatic", aoi.markerStatic)

    // delete those NOT MOVED from move set and add them to static set
    for id, obj := range aoi.watcherMove {
        if obj.version != aoi.version {
            aoi.watcherStatic[id] = obj
            delete(aoi.watcherMove, id)
            // fmt.Println("// delete those NOT MOVED from move set and add them to static set")
        }
    }
    for id, obj := range aoi.markerMove {
        if obj.version != aoi.version {
            aoi.markerStatic[id] = obj
            delete(aoi.markerMove, id)
        }
    }

    fmt.Println("after process-------------")
    fmt.Println("watchMove:", aoi.watcherMove)
    fmt.Println("markerMove", aoi.markerMove)
    fmt.Println("watcherStatic", aoi.watcherStatic)
    fmt.Println("markerStatic", aoi.markerStatic)

    // compare watcherMove and all marker
    for watcherId, watcherObj := range aoi.watcherMove {
        for markerId, markerObj := range aoi.markerMove {
            if watcherObj != markerObj {
                dist := distance2(watcherObj.pos, markerObj.pos)
                fmt.Printf("watcherMove:%d -- markerMove:%d, dist=%f\n", watcherId, markerId, dist)
                if dist < RADIS*RADIS {
                    cb(watcherId, markerId)
                }
            }
        }
        for markerId, markerObj := range aoi.markerStatic {
            if watcherObj != markerObj {
                dist := distance2(watcherObj.pos, markerObj.pos)
                fmt.Printf("watcherMove:%d -- markerStatic:%d, dist=%f\n", watcherId, markerId, dist)

                if dist < RADIS*RADIS {
                    cb(watcherId, markerId)
                }
            }
        }
    }

    // compare static watcher to markerMove
    for watcherId, watcherObj := range aoi.watcherStatic {
        for markerId, markerObj := range aoi.markerMove {
            if watcherObj != markerObj {
                dist := distance2(watcherObj.pos, markerObj.pos)
                fmt.Printf("watcherStatic:%d -- markerMove:%d, dist=%f\n", watcherId, markerId, dist)

                if dist < RADIS*RADIS {
                    cb(watcherId, markerId)
                }
            }
        }
    }

    // aoi.moveQueue = aoi.moveQueue[:0]
    aoi.version++
}

func distance2(p1 FPoint, p2 FPoint) float32 {
    return (p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y)
}
