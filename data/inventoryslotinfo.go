package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type InventorySlotInfo struct {
    PCItemInfo

    InvenX CoordInven_t
    InvenY CoordInven_t
}

func (slot *InventorySlotInfo) Read(reader io.Reader) error {
    slot.PCItemInfo.Read(reader)

    binary.Read(reader, binary.LittleEndian, &slot.InvenX)
    binary.Read(reader, binary.LittleEndian, &slot.InvenY)
    return nil
}

func (slot *InventorySlotInfo) Write(writer io.Writer) error {
    slot.PCItemInfo.Write(writer)

    binary.Write(writer, binary.LittleEndian, slot.InvenX)
    binary.Write(writer, binary.LittleEndian, slot.InvenY)
    return nil
}

func (info *InventorySlotInfo) Size() uint32 {
    return 2 + info.PCItemInfo.Size()
}
