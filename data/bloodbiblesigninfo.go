package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type BloodBibleSignInfo struct {
    OpenNum  uint32
    SignList []ItemType_t
}

func (info *BloodBibleSignInfo) Size() uint32 {
    var sz uint32
    sz = 5 + uint32(len(info.SignList))*2
    return sz
}

func (info *BloodBibleSignInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, info.OpenNum)
    num := uint8(len(info.SignList))
    binary.Write(writer, binary.LittleEndian, num)
    for i := 0; i < len(info.SignList); i++ {
        binary.Write(writer, binary.LittleEndian, info.SignList[i])
    }
    return nil
}

func (info *BloodBibleSignInfo) Read(reader io.Reader) {
    binary.Read(reader, binary.LittleEndian, &info.OpenNum)
    var num uint8
    binary.Read(reader, binary.LittleEndian, &num)
    info.SignList = make([]ItemType_t, num)
    for i := 0; i < int(num); i++ {
        binary.Read(reader, binary.LittleEndian, &info.SignList[i])
    }
    return
}
