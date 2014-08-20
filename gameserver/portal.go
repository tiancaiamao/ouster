package main

type PortalType_t byte

const (
    PORTAL_CLASS_NORMAL = iota
    PORTAL_CLASS_PRIVATE
    PORTAL_CLASS_MULTI
    PORTAL_CLASS_TRIGGERED
)

type PortalType byte

const (
    PORTAL_NORMAL PortalType = iota
    PORTAL_SLAYER
    PORTAL_VAMPIRE
    PORTAL_MULTI_TARGET
    PORTAL_PRIVATE
    PORTAL_GUILD
    PORTAL_BATTLE
    PORTAL_OUSTER
)