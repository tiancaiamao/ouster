package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type NPCInfo struct {
    Name  string
    NPCID NPCID_t
    X     ZoneCoord_t
    Y     ZoneCoord_t
}

func (info *NPCInfo) Read(reader io.Reader) error {
    var szName uint8
    var buf [256]byte
    binary.Read(reader, binary.LittleEndian, &szName)
    if szName > 0 {
        _, err := reader.Read(buf[:szName])
        if err != nil {
            return err
        }

        info.Name = string(buf[:szName])
        binary.Read(reader, binary.LittleEndian, &info.NPCID)
        binary.Read(reader, binary.LittleEndian, &info.X)
        binary.Read(reader, binary.LittleEndian, &info.Y)
    }
    return nil
}

func (info *NPCInfo) Write(writer io.Writer) error {
    binary.Write(writer, binary.LittleEndian, uint8(len(info.Name)))
    if len(info.Name) > 0 {
        io.WriteString(writer, info.Name)
        binary.Write(writer, binary.LittleEndian, info.NPCID)
        binary.Write(writer, binary.LittleEndian, info.X)
        binary.Write(writer, binary.LittleEndian, info.Y)
    }
    return nil
}
