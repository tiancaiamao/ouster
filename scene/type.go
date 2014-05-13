package scene

type ObjectIDType uint32

const (
	ObjectIDMaskNPC = 1 << 31
)

func (id ObjectIDType) Monster() bool {
	return (id & ObjectIDMaskNPC) != 0
}

func (id ObjectIDType) Index() int {
	return int(id &^ ObjectIDMaskNPC)
}
