package data

import (
	"testing"
	"bytes"
)

func TestPCOusterInfoDump(t *testing.T) {
	info := PCOusterInfo {
		PCType: 'O',
		ObjectID: 1,
		Name: "test", 
		Level: 150,
	
		HairColor: 101,
	
		Alignment:7500,
		STR: [3]uint16{10, 10, 10},
		DEX: [3]uint16{25, 25, 25},
		INT: [3]uint16{10, 10, 10},
	
		HP: [2]uint16{315, 315},
		MP: [2]uint16{186, 111},
	
		Rank:50,
		RankExp:10700,
		Exp:125,
	
		Fame:500,
		Gold:68,
		Sight :13,
		Bonus :9999,
		SkillBonus :9999,
	
		Competence:1,
		GuildID:66,

		GuildMemberRank:4,
		AdvancementLevel:100,
	}
	
	buf := &bytes.Buffer{}
	info.Dump(buf)
	
	right := []byte{1,0,0,0,4,116,101,115,116,150,0,101,0,0,76,29,0,0,10,0,10,0,10,0,25,0,25,0,25,0,10,0,10,0,10,0,59,1,59,1,186,0,111,0,50,204,41,0,0,125,0,0,0,244,1,0,0,68,0,0,0,13,15,39,15,39,0,0,1,66,0,0,4,0,0,0,0,100,0,0,0,0}
	if !bytes.Equal(buf.Bytes(), right) {
		t.Failed()
	}
}

