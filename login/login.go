package login

import (
	"encoding/json"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster"
	"io/ioutil"
	"net"
	"log"
	"time"
)

type LoginError string

func (l LoginError) Error() string {
	return "login error:" + string(l)
}

// TODO: this should be written as a state machine, provide api something like \n
// session :=login.New()	// create a login new session, which is actually the state\n
// session.Login(conn)
func Login(conn net.Conn) (*data.Player, error) {
	err := conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
	if err != nil {
		return nil, err
	}

	//------------------
	p, err := packet.Read(conn)
	if err != nil {
		log.Println("readPacket error")
		return nil, ouster.NewError(err.Error())
	}
	pkt, ok := p.(packet.LoginPacket)
	if !ok {
		log.Println("expect a LoginPacke")
		return nil, LoginError("expect a LoginPacket")
	}

	log.Println("get a LoginPacket...run here")
	// check username and ignore password ...check whether this user exist...and so on

	charactor, err := loadUser(pkt.Username)
	if err != nil {
		return nil, LoginError("not a valid user")
	}

	err = packet.Write(conn, packet.PCharactorInfo, charactor)
	if err != nil {
		return nil, ouster.NewError(err.Error())
	}

	//------------------
	p, err = packet.Read(conn)
	if err != nil {
		return nil, LoginError(err.Error())
	}
	_, ok = p.(packet.SelectCharactorPacket)
	if !ok {
		return nil, LoginError("expect a SelectCharactorPacket")
	}
	log.Println("run here get a SelectCharactorPacket")

	// load player info ...
	player, err := loadCharactor(config.DataDir+"/player/Delrek")
	if err != nil {
		return nil, ouster.NewError(err.Error())
	}

	err = packet.Write(conn, packet.PLoginOk, packet.LoginOkPacket{})
	if err != nil {
		return nil, LoginError("write LoginOkPacket error")
	}

	return player, nil
}

func loadUser(name string) (packet.CharactorInfoPacket, error) {
	return packet.CharactorInfoPacket{
		Name:  "test",
		Class: "Brute",
		Level: 1,
	}, nil
}

func loadCharactor(filePath string) (*data.Player, error) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, LoginError("no such charactor!")
	}
	
	var player data.Player
	err = json.Unmarshal(buf, &player)
	if err != nil {
		return nil, LoginError("error charactor info!")
	}
	return &player, nil
}
