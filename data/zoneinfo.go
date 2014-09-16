package data

import (
    "encoding/json"
    "github.com/tiancaiamao/ouster/config"
    // "github.com/tiancaiamao/ouster/log"
    . "github.com/tiancaiamao/ouster/util"
    "os"
    "path"
)

type ZoneInfoItem struct {
    ZoneID   ZoneID_t
    ZoneInfo ZoneInfo
}

type Monster struct {
    MonsterType MonsterType_t
    Count       int
}

type ZoneInfo struct {
    EventMonsterList []Monster
    MonsterList      []Monster
}

var ZoneInfoTable map[ZoneID_t]ZoneInfo

func init() {
    var array []ZoneInfoItem

    filePath := path.Join(config.DataFilePath, "zoneinfo.json")
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }

    dec := json.NewDecoder(file)
    err = dec.Decode(&array)
    if err != nil {
        panic(err)
    }

    ZoneInfoTable = make(map[ZoneID_t]ZoneInfo)
    for _, v := range array {
        ZoneInfoTable[v.ZoneID] = v.ZoneInfo
    }
}
