package data

import (
    "bytes"
    "io"
    "reflect"
    "testing"
)

type Dumper interface {
    Write(io.Writer) error
    Size() uint32
}

func TestSize(t *testing.T) {
    arr := []Dumper{
        &GearInfo{},
        &BloodBibleSignInfo{},
        &EffectInfo{},
        &ExtraInfo{},
        &GearSlotInfo{},
        &InventoryInfo{},
        &InventorySlotInfo{},
        &NPCInfo{},
        // &PCInfo{},
        &PCItemInfo{},
        VampireSkillInfo{},
        OusterSkillInfo{},
        // SlayerSkillInfo{},
        &SubItemInfo{},
        &NicknameInfo{},
    }

    buf := &bytes.Buffer{}
    for _, v := range arr {
        buf.Reset()
        v.Write(buf)
        if uint32(buf.Len()) != v.Size() {
            t.Errorf("%s的Size不对:期待%d 实际%d\n", reflect.TypeOf(v), buf.Len(), v.Size())
        }
    }
}
