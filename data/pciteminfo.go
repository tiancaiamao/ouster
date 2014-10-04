package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type PCItemInfo struct {
    ObjectID        ObjectID_t
    IClass          byte
    ItemType        ItemType_t
    OptionType      []OptionType_t
    Durability      Durability_t
    Silver          Silver_t
    Grade           Grade_t
    EnchantLevel    EnchantLevel_t
    ItemNum         ItemNum_t
    MainColor       uint16
    SubItemInfoList []SubItemInfo
}

func (info *PCItemInfo) Read(reader io.Reader) {
    binary.Read(reader, binary.LittleEndian, &info.ObjectID)
    binary.Read(reader, binary.LittleEndian, &info.IClass)
    binary.Read(reader, binary.LittleEndian, &info.ItemType)

    var optionSize uint8
    binary.Read(reader, binary.LittleEndian, &optionSize)
    info.OptionType = make([]OptionType_t, optionSize)
    for i := 0; i < int(optionSize); i++ {
        binary.Read(reader, binary.LittleEndian, &info.OptionType[i])
    }

    binary.Read(reader, binary.LittleEndian, &info.Durability)
    binary.Read(reader, binary.LittleEndian, &info.Silver)
    binary.Read(reader, binary.LittleEndian, &info.Grade)
    binary.Read(reader, binary.LittleEndian, &info.EnchantLevel)
    binary.Read(reader, binary.LittleEndian, &info.ItemNum)
    binary.Read(reader, binary.LittleEndian, &info.MainColor)

    var num uint8
    binary.Read(reader, binary.LittleEndian, &num)
    info.SubItemInfoList = make([]SubItemInfo, num)
    for i := 0; i < int(num); i++ {
        info.SubItemInfoList[i].Read(reader)
    }
}

func (info *PCItemInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.ObjectID)
    binary.Write(writer, binary.LittleEndian, info.IClass)
    binary.Write(writer, binary.LittleEndian, info.ItemType)

    optionSize := uint8(len(info.OptionType))
    binary.Write(writer, binary.LittleEndian, optionSize)
    for i := 0; i < int(optionSize); i++ {
        binary.Write(writer, binary.LittleEndian, &info.OptionType[i])
    }

    binary.Write(writer, binary.LittleEndian, info.Durability)
    binary.Write(writer, binary.LittleEndian, info.Silver)
    binary.Write(writer, binary.LittleEndian, info.Grade)
    binary.Write(writer, binary.LittleEndian, info.EnchantLevel)
    binary.Write(writer, binary.LittleEndian, info.ItemNum)
    binary.Write(writer, binary.LittleEndian, info.MainColor)

    var num uint8
    binary.Write(writer, binary.LittleEndian, num)
    for i := 0; i < int(num); i++ {
        info.SubItemInfoList[i].Write(writer)
    }
    return nil
}
