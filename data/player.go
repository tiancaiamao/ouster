package data

type PlayerClass uint8
const (
	_ = iota
	BRUTE = iota
	Knight
)
type Player struct {
	Name string
	Class PlayerClass
	Level uint8
	Map string
	XP int
	HP int
	MP int
	Carried []int
	Pos Point  // NOTE:this property belongs to Map
}
