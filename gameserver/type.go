package scene

type ObjectIDType uint32

const (
	ObjectIDMaskNPC = 1 << 31
)

func (id ObjectIDType) Monster() bool {
	return (id & ObjectIDMaskNPC) != 0
}

func (id ObjectIDType) Player() bool {
	return (id & ObjectIDMaskNPC) == 0
}

func (id ObjectIDType) Index() uint32 {
	return uint32(id &^ ObjectIDMaskNPC)
}
