package main

import (
    "encoding/json"
    // "fmt"
    "github.com/tiancaiamao/ouster/data"
    . "github.com/tiancaiamao/ouster/util"
    "os"
    "strconv"
)

type Item struct {
    MType          string
    SType          string
    EName          string
    Level          string
    STR            string
    DEX            string
    INTE           string
    BSize          string
    Fame           string
    Exp            string
    MColor         string
    SColor         string
    Align          string
    AOrder         string
    Moral          string
    Delay          string
    Sight          string
    MeleeRange     string
    MissileRange   string
    RegenPortal    string
    RegenInvisible string
    RegenBat       string
    UnburrowChance string
    MMode          string
    AIType         string
    Master         string
    ClanType       string
    Chief          string
    NormalRegin    string
    HasTreasure    string
    MonsterClass   string
    SkullType      string
}

func main() {
    file, err := os.Open("MonsterInfo.json")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var items []Item
    dec := json.NewDecoder(file)
    err = dec.Decode(&items)
    if err != nil {
        panic(err)
    }

    infos := make([]data.MonsterInfo, len(items))
    for i := 0; i < len(items); i++ {
        item := &items[i]
        info := &infos[i]

        var tmp int
        tmp, _ = strconv.Atoi(item.MType)
        info.MonsterType = MonsterType_t(tmp)

        info.Name = item.EName

        tmp, _ = strconv.Atoi(item.Level)
        info.Level = Level_t(tmp)

        tmp, _ = strconv.Atoi(item.STR)
        info.STR = Attr_t(tmp)

        tmp, _ = strconv.Atoi(item.DEX)
        info.DEX = Attr_t(tmp)

        tmp, _ = strconv.Atoi(item.INTE)
        info.INTE = Attr_t(tmp)

        tmp, _ = strconv.Atoi(item.BSize)
        info.BodySize = tmp

        tmp, _ = strconv.Atoi(item.Exp)
        info.Exp = Exp_t(tmp)

        tmp, _ = strconv.Atoi(item.MColor)
        info.MColor = Color_t(tmp)

        tmp, _ = strconv.Atoi(item.SColor)
        info.SColor = Color_t(tmp)

        tmp, _ = strconv.Atoi(item.Sight)
        info.Sight = Sight_t(tmp)

        switch item.MMode {
        case "WALK":
            info.MoveMode = MOVE_MODE_WALKING
        case "FLY":
            info.MoveMode = MOVE_MODE_FLYING
        }

        tmp, _ = strconv.Atoi(item.MeleeRange)
        info.MeleeRange = tmp

        tmp, _ = strconv.Atoi(item.MissileRange)
        info.MissileRange = tmp

        tmp, _ = strconv.Atoi(item.AIType)
        info.AIType = tmp

        tmp, _ = strconv.Atoi(item.UnburrowChance)
        info.UnburrowChance = tmp

        tmp, _ = strconv.Atoi(item.SType)
        info.SType = tmp

        tmp, _ = strconv.Atoi(item.Fame)
        info.Fame = Fame_t(tmp)

        tmp, _ = strconv.Atoi(item.Align)
        info.Align = tmp

        tmp, _ = strconv.Atoi(item.AOrder)
        info.AOrder = tmp

        tmp, _ = strconv.Atoi(item.Moral)
        info.Moral = tmp

        tmp, _ = strconv.Atoi(item.Delay)
        info.Delay = tmp

        tmp, _ = strconv.Atoi(item.RegenPortal)
        info.RegenPortal = (tmp != 0)

        tmp, _ = strconv.Atoi(item.RegenInvisible)
        info.RegenInvisible = (tmp != 0)

        tmp, _ = strconv.Atoi(item.RegenBat)
        info.RegenBat = (tmp != 0)

        tmp, _ = strconv.Atoi(item.Master)
        info.Master = (tmp != 0)

        tmp, _ = strconv.Atoi(item.Chief)
        info.Chief = (tmp != 0)

        tmp, _ = strconv.Atoi(item.NormalRegin)
        info.NormalRegin = (tmp != 0)

        tmp, _ = strconv.Atoi(item.ClanType)
        info.ClanType = tmp

        tmp, _ = strconv.Atoi(item.MonsterClass)
        info.MonsterClass = tmp

        tmp, _ = strconv.Atoi(item.SkullType)
        info.SkullType = tmp
    }

    out, err := os.Create("monsterinfo__.json")
    if err != nil {
        panic(err)
    }
    defer out.Close()
    enc := json.NewEncoder(out)
    err = enc.Encode(infos)
    if err != nil {
        panic(err)
    }
    return
}
