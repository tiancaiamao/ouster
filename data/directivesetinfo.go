package data

import (
    "encoding/json"
    "github.com/tiancaiamao/ouster/config"
    "os"
    "path"
)

type Directive struct {
    Conditions []int
    Action     int
    Parameter  int
    Ratio      int
    Weight     int
}

type DirectiveSetItem struct {
    ID           int
    DirectiveSet DirectiveSet
}

type DirectiveSet struct {
    Directives     []*Directive
    DeadDirectives []*Directive
    Name           string

    bAttackAir   bool
    bSeeSafeZone bool
}

var DirectiveSetTable map[int]DirectiveSet

func init() {
    var array []DirectiveSetItem

    filePath := path.Join(config.DataFilePath, "directivesetinfo.json")
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }

    dec := json.NewDecoder(file)
    err = dec.Decode(&array)
    if err != nil {
        panic(err)
    }

    DirectiveSetTable = make(map[int]DirectiveSet)
    for _, v := range array {
        DirectiveSetTable[v.ID] = v.DirectiveSet
    }
}
