package main

type ZoneInfo struct {
    IsPKZone       bool   `json:"isPkZone"`
    IsNoPortalZone bool   `json:"isNoPortalZone"`
    IsMasterLair   bool   `json:"isMasterLair"`
    IsHolyLand     bool   `json:"isHolyLand"`
    SMPFileName    string `json:"smpFileName"`
    SMIFileName    string `json:"smiFileName"`
}

var gZoneInfoManager ZoneInfoManager

type ZoneInfoManager map[ZoneID_t]*ZoneInfo

func (zm ZoneInfoManager) GetZoneInfo(id ZoneID_t) *ZoneInfo {
    return zm[id]
}
