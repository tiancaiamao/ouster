package data

import (
    "bytes"
    "github.com/tiancaiamao/ouster/util"
    "testing"
)

func TestGearInfoWrite(t *testing.T) {
    info := &GearInfo{
        GearSlotInfoList: []GearSlotInfo{
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0x76de, IClass: 0x23, ItemType: 0x0, OptionType: []util.OptionType_t{}, Durability: 0xbb8, Silver: 0x0, Grade: 4, EnchantLevel: 0x0, ItemNum: 0x1, MainColor: 0x0, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0x76df00, IClass: 0x0, ItemType: 0x16, OptionType: []util.OptionType_t{}, Durability: 0xdac00, Silver: 0x0, Grade: 1024, EnchantLevel: 0x0, ItemNum: 0x0, MainColor: 0x0, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0x76e00002, IClass: 0x0, ItemType: 0x400, OptionType: []util.OptionType_t{0x0, 0x0}, Durability: 0x1, Silver: 0x0, Grade: -1, EnchantLevel: 0x0, ItemNum: 0x14, MainColor: 0x0, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x4},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0x76e100, IClass: 0x0, ItemType: 0xf, OptionType: []util.OptionType_t{}, Durability: 0x157c00, Silver: 0x0, Grade: 1024, EnchantLevel: 0x0, ItemNum: 0x0, MainColor: 0x1, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0x76e20302, IClass: 0x0, ItemType: 0xe00, OptionType: []util.OptionType_t{}, Durability: 0x11940000, Silver: 0x0, Grade: 262144, EnchantLevel: 0x0, ItemNum: 0x0, MainColor: 0x100, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0xe3030400, IClass: 0x76, ItemType: 0x0, OptionType: []util.OptionType_t{0x0}, Durability: 0x10000, Silver: 0x0, Grade: -65536, EnchantLevel: 0xff, ItemNum: 0xff, MainColor: 0x900, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
            GearSlotInfo{PCItemInfo: PCItemInfo{ObjectID: 0xe4040900, IClass: 0x76, ItemType: 0x0, OptionType: []util.OptionType_t{0x5}, Durability: 0x10000, Silver: 0x0, Grade: -65536, EnchantLevel: 0xff, ItemNum: 0xff, MainColor: 0x900, SubItemInfoList: []SubItemInfo{}}, SlotID: 0x0},
        },
    }
    buf := &bytes.Buffer{}
    info.Write(buf)
    right := []byte{7, 222, 118, 0, 0, 35, 0, 0, 0, 184, 11, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 223, 118, 0, 0, 22, 0, 0, 0, 172, 13, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 224, 118, 0, 0, 4, 2, 0, 0, 1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 20, 0, 0, 0, 4, 0, 225, 118, 0, 0, 15, 0, 0, 0, 124, 21, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 2, 3, 226, 118, 0, 0, 14, 0, 0, 0, 148, 17, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1, 0, 0, 0, 4, 3, 227, 118, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 9, 0, 0, 0, 9, 4, 228, 118, 0, 0, 1, 5, 0, 0, 1, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 9, 0, 0}
    if !bytes.Equal(right, buf.Bytes()) {
        t.Errorf("GearInfo error")
    }
}
