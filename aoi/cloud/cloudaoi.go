// 解决的根本问题就是，哪个点，周围有哪些东西。啊个点有哪些东西，最好由player维护。
// 这个接口是只能由Map调的，不能用Player直接掉。
// 间接的思路是，不维护东西链表。
package cloud

// import "fmt"

type FPoint struct {
	X float32
	Y float32
}

type object struct {
	// id uint32
	mode    AoiMode
	pos     FPoint
	last    FPoint
	version uint32
}

type key struct {
	smaller uint32
	larger  uint32
}

type value struct {
	obj1 *object
	obj2 *object
}

type kv struct {
	id  uint32
	obj *object
}

type Aoi struct {
	all     map[uint32]*object
	move    []kv
	change  []kv
	version uint32
	hotpair map[key]value
}

func New() *Aoi {
	aoi := new(Aoi)
	aoi.all = make(map[uint32]*object)
	aoi.move = make([]kv, 0, 100)
	aoi.change = make([]kv, 0, 100)
	aoi.hotpair = make(map[key]value)
	return aoi
}

func (aoi *Aoi) Update(id uint32, mode AoiMode, pos FPoint) {
	if (mode & ModeDrop) != 0 {
		delete(aoi.all, id)
		return
	}

	obj, ok := aoi.all[id]
	if !ok {
		obj = new(object)
		aoi.all[id] = obj
	}

	changed := modeChanged(obj, mode)
	if changed {
		aoi.change = append(aoi.change, kv{id, obj})
		obj.mode = mode
		obj.last = pos
	} else {
		moved := posChanged(obj, pos)
		if moved {
			// fmt.Printf("object moved:%d (%f %f)\n", id, pos.X, pos.Y)
			aoi.move = append(aoi.move, kv{id, obj})
			obj.last = pos
		}
	}
	obj.pos = pos
	obj.version = aoi.version
}

func posChanged(obj *object, pos FPoint) bool {
	var ret bool
	if 4*distance2(obj.last, pos) >= RADIS*RADIS {
		ret = true
	}
	return ret
}

func (aoi *Aoi) addHotPair(id1 uint32, obj1 *object, id2 uint32, obj2 *object) {
	if id1 < id2 {
		aoi.hotpair[key{id1, id2}] = value{obj1, obj2}
	} else if id1 > id2 {
		aoi.hotpair[key{id2, id1}] = value{obj2, obj1}
	}
}

type AoiCallback func(uint32, uint32)

func send(id1 uint32, obj1 *object, id2 uint32, obj2 *object, cb AoiCallback) {
	if (obj1.mode&ModeWatcher) != 0 && (obj2.mode&ModeMarker) != 0 {
		cb(id1, id2)
	}
	if (obj2.mode&ModeWatcher) != 0 && (obj1.mode&ModeMarker) != 0 {
		cb(id2, id1)
	}
}

func (aoi *Aoi) Message(cb AoiCallback) {
	// fmt.Println("enter Message-------------")
	// fmt.Println("move:", aoi.move)
	// fmt.Println("change:", aoi.change)

	// compare change and all to generate hotpair and send message
	for _, kv := range aoi.change {
		for id, obj := range aoi.all {
			if kv.id != id {
				dist := distance2(kv.obj.pos, obj.pos)
				if dist <= RADIS*RADIS {
					send(kv.id, kv.obj, id, obj, cb)
				} else if dist > RADIS*RADIS && dist < 4*RADIS*RADIS {
					aoi.addHotPair(kv.id, kv.obj, id, obj)
				}
			}
		}
	}

	// compare move and all to generate hotpair
	for _, kv := range aoi.move {
		for id, obj := range aoi.all {
			if kv.id != id {
				dist := distance2(kv.obj.pos, obj.pos)
				if dist > RADIS*RADIS && dist < 4*RADIS*RADIS {
					aoi.addHotPair(kv.id, kv.obj, id, obj)
				}
			}
		}
	}

	// fmt.Println("hotpair:", aoi.hotpair)

	for k, v := range aoi.hotpair {
		// filter those changed in update function
		if v.obj1.version == aoi.version || v.obj2.version == aoi.version {
			if distance2(v.obj1.pos, v.obj2.pos) < RADIS*RADIS {
				send(k.smaller, v.obj1, k.larger, v.obj2, cb)
				delete(aoi.hotpair, k)
			}
		}
	}

	aoi.move = aoi.move[:0]
	aoi.change = aoi.change[:0]
	aoi.version++
}
