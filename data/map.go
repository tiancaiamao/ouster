package data

// import ("fmt")

type EnemyGroupType uint32
const (
	_ = iota

)

type Rect struct {
	X float32
	Y float32
	W float32
	H float32
}

type Point struct {
	X int
	Y int
}

type EnemyGroup struct {
	Type      EnemyType
	Location  Rect //Location area for enemygroup
	LevelMin  int  //Defines the level range of enemies in group.
	LevelMax  int
	NumberMin int //Defines the range of enemies in group. If only one number is given, it's the exact amount.
	NumberMax int
	Chance    float32 // Percentage of chance
}

type LayerType uint32
const (
	_ = iota
	BACKGROUND LayerType = iota
	OBJECT
	COLLISION
)
type Layer struct {
	Type   LayerType
	// Format string // must be dec
	Data   []uint16
}

type EnemyType uint32
const (
	_ = iota
	CURSED_GRAVE EnemyType = iota
	SKELETAL_ARCHER
	UNDEAD 
	WYVERN
)

type Enemy struct {
	Type              EnemyType
	Location          Point  //Location of enemy
	Direction         int    //Direction of enemy
	WayPoint          Point  //Enemy waypoint
	WanderArea        Rect   // Wander area for enemy.
	RequiresStatus    string // Status required for enemy load
	RequiresNotStatus string // Status required to be missing for enemy load
}

// type Enemy struct {
// 	typ string
// 	pos FPoint
// 	direction int
// 	wayPoints []FPoint
// 	wander bool
// 	wanderArea Rect
// 	heroAlly bool
// 	summonPowerIndex int
// 	summoner *StatBlock
// }

type NpcType uint32
type Npc struct {
	Type           NpcType
	RequiresStatus string //Status required for NPC load. There can be multiple states, separated by comma
	RequiresNot    string //Status required to be missing for NPC load. There can be multiple states, separated by comma
	Location       Point  //  Location of NPC
}

type MapType uint32
type Map struct {
	Title       string //Title of map
	Width       uint16 //Height of map
	Height      uint16 //Width of map
	// Tileset     string //Tileset use for map
	Location    Point  //Spawn point location in map
	Layers      []Layer
	Enemies     []Enemy
	EnemyGroups []EnemyGroup
	Npcs        []Npc
	Name        string // Uniq identifier for the effect definition.
	Type        MapType
	Additive    bool
}

// func init() {
// 	fmt.Println("background is ", BACKGROUND)
// 	fmt.Println("object is ", OBJECT)
// }
