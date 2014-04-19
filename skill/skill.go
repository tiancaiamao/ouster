package skill

import "github.com/tiancaiamao/ouster"

type SelfSkill interface {
	ExecuteSelf() Effect
}

type TargetSkill interface {
	ExecuteTarget(from, to ouster.Creature) (int, bool)
}

type RegionSkill interface {
	ExecuteRegion() Effect
}

type Effect struct {
	Damage  int // increase damage
	Dodge   int // increase dodge
	Defense int // increase Defence
	Harm    int // cause Harm
	ToHit   int // increase tohit
}

func Query(id int) interface{} {
	return nil
}
