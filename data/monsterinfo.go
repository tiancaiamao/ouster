package data

import (
    "encoding/json"
    "github.com/tiancaiamao/ouster/config"
    . "github.com/tiancaiamao/ouster/util"
    "os"
    "path"
)

type MonsterInfoItem struct {
    MonsterType MonsterType_t
    MonsterInfo MonsterInfo
}

type MonsterInfo struct {
    Name         string
    Level        uint8
    STR          Attr_t
    DEX          Attr_t
    INTE         Attr_t
    BodySize     uint
    HP           HP_t
    Exp          Exp_t
    MColor       uint16
    SColor       uint16
    Sight        uint8
    MeleeRange   int
    MissileRange int
}

var MonsterInfoTable map[MonsterType_t]MonsterInfo

func init() {
    var array []MonsterInfoItem

    filePath := path.Join(config.DataFilePath, "monsterinfo.json")
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }

    dec := json.NewDecoder(file)
    err = dec.Decode(&array)
    if err != nil {
        panic(err)
    }

    MonsterInfoTable = make(map[MonsterType_t]MonsterInfo)
    for _, v := range array {
        MonsterInfoTable[v.MonsterType] = v.MonsterInfo
    }
}
