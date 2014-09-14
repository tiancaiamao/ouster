package main

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

// TODO
func (manager *NPCManager) addCreature(*NPC) {
}

// TODO
func (m *NPCManager) heartbeat() {
}
