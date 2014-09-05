package main

import (
    . "github.com/tiancaiamao/ouster/util"
)

type ZoneInfo struct {
    IsPKZone       bool   `json:"isPkZone"`
    IsNoPortalZone bool   `json:"isNoPortalZone"`
    IsMasterLair   bool   `json:"isMasterLair"`
    IsHolyLand     bool   `json:"isHolyLand"`
    SMPFileName    string `json:"smpFileName"`
    SSIFileName    string `json:"ssiFileName"`
}

// 用于维护一个从ZoneID到ZoneInfo的映射关系
var gZoneInfoManager ZoneInfoManager

type ZoneInfoManager map[ZoneID_t]*ZoneInfo

func (zm ZoneInfoManager) GetZoneInfo(id ZoneID_t) *ZoneInfo {
    return zm[id]
}
