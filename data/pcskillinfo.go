package data

import (
    "encoding/binary"
    . "github.com/tiancaiamao/ouster/util"
    "io"
)

type PCSkillInfo interface {
    Write(io.Writer)
    Size() uint32
}
type VampireSkillInfo struct {
    LearnNewSkill           bool
    SubVampireSkillInfoList []SubVampireSkillInfo
}

func (info VampireSkillInfo) Size() uint32 {
    var sz uint32
    sz = 2
    for i := 0; i < len(info.SubVampireSkillInfoList); i++ {
        sz += info.SubVampireSkillInfoList[i].Size()
    }
    return sz
}
func (info VampireSkillInfo) Write(writer io.Writer) {
    if info.LearnNewSkill {
        binary.Write(writer, binary.LittleEndian, uint8(1))
    } else {
        binary.Write(writer, binary.LittleEndian, uint8(0))
    }

    binary.Write(writer, binary.LittleEndian, uint8(len(info.SubVampireSkillInfoList)))
    for _, v := range info.SubVampireSkillInfoList {
        v.Write(writer)
    }
}

type SubVampireSkillInfo struct {
    SkillType   SkillType_t
    Interval    uint32
    CastingTime uint32
}

func (info *SubVampireSkillInfo) Size() uint32 {
    return 10
}
func (info SubVampireSkillInfo) Write(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, info.SkillType)
    binary.Write(writer, binary.LittleEndian, info.Interval)
    binary.Write(writer, binary.LittleEndian, info.CastingTime)
    return
}

type OusterSkillInfo struct {
    LearnNewSkill          bool
    SubOusterSkillInfoList []SubOusterSkillInfo
}

func (info OusterSkillInfo) Size() uint32 {
    var sz uint32
    sz = 2
    for i := 0; i < len(info.SubOusterSkillInfoList); i++ {
        sz += info.SubOusterSkillInfoList[i].Size()
    }
    return sz
}

func (info OusterSkillInfo) Write(writer io.Writer) {
    if info.LearnNewSkill {
        binary.Write(writer, binary.LittleEndian, uint8(1))
    } else {
        binary.Write(writer, binary.LittleEndian, uint8(0))
    }

    binary.Write(writer, binary.LittleEndian, uint8(len(info.SubOusterSkillInfoList)))
    for _, v := range info.SubOusterSkillInfoList {
        v.Write(writer)
    }
}

type SubOusterSkillInfo struct {
    SkillType   SkillType_t
    ExpLevel    ExpLevel_t
    Interval    uint32
    CastingTime uint32
}

func (info *SubOusterSkillInfo) Size() uint32 {
    return 12
}

func (info SubOusterSkillInfo) Write(writer io.Writer) {
    binary.Write(writer, binary.LittleEndian, info.SkillType)
    binary.Write(writer, binary.LittleEndian, info.ExpLevel)
    binary.Write(writer, binary.LittleEndian, info.Interval)
    binary.Write(writer, binary.LittleEndian, info.CastingTime)
    return
}
