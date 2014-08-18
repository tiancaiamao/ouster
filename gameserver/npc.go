package main


type NPC struct {
	Creature
}

func (npc *NPC) CreatureClass() CreatureClass {
	return CREATURE_CLASS_NPC
}