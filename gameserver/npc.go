package main

import (
    . "github.com/tiancaiamao/ouster/util"
)

type NPC struct {
    Creature
}

func (npc *NPC) CreatureClass() CreatureClass {
    return CREATURE_CLASS_NPC
}

type NPCManager struct {
}

func NewNPCManager() *NPCManager {
    return &NPCManager{}
}

func (npc *NPC) getProtection() Protection_t {
    return 0
}

// TODO
func (manager *NPCManager) addCreature(*NPC) {
}

// TODO
func (m *NPCManager) heartbeat() {
}
