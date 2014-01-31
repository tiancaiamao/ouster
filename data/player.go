package data

type PlayerClass uint8
const (
	_ = iota
	BRUTE = iota
)
type Player struct {
	Name string
	Class PlayerClass
	Level uint8
	Map string
	HP int
	MP int
	Carried []int
	Pos Point  // NOTE:this property belongs to Map
}
