package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type GearSlotInfo struct {
    PCItemInfo

    SlotID SlotID_t
}

func (info *GearSlotInfo) Read(reader io.Reader) error {
    info.PCItemInfo.Read(reader)
    err := binary.Read(reader, binary.LittleEndian, &info.SlotID)
    return err
}

func (info *GearSlotInfo) Write(writer io.Writer) error {
    err := info.PCItemInfo.Write(writer)
    if err != nil {
        return err
    }

    err = binary.Write(writer, binary.LittleEndian, info.SlotID)
    return err
}
