package login

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"github.com/tiancaiamao/ouster/config"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type LoginError string

func (l LoginError) Error() string {
	return "login error:" + string(l)
}

func readPacket(conn io.Reader) (interface{}, error) {
	var header [4]byte
	_, err := io.ReadFull(conn, header[:])
	if err != nil {
		return nil, LoginError(err.Error())
	}

	log.Println("read packet head:", header)

	size := binary.BigEndian.Uint32(header[:])
	log.Println("read a packet, header size:", size)
	data := make([]byte, size)

	_, err = io.ReadFull(conn, data)
	if err != nil {
		log.Println("ReadFull error")
		return nil, LoginError(err.Error())
	}

	log.Println("read packet data:", data)
	// packet解析
	p, err := packet.Parse(data)
	if err != nil {
		log.Println("parse packet error")
		return nil, LoginError(err.Error())
	}

	log.Println("run here.......................")

	return p, nil
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
		return nil, ouster.NewError(err)
	}
	pkt, ok := p.(*packet.LoginPacket)
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
		return nil, ouster.NewError(err)
	}
	log.Println("send a CharactorInfoPacket:", buf)

	//------------------
	p, err = packet.Read(conn)
	if err != nil {
		return nil, LoginError(err.Error())
	}
	_, ok = p.(*packet.SelectCharactorPacket)
	if !ok {
		return nil, LoginError("expect a SelectCharactorPacket")
	}
	log.Println("run here get a SelectCharactorPacket")

	// load player info ...
	player, err := loadCharactor(config.DataDir+"/player/Delrek")
	if err != nil {
		return nil, LoginError(err.Error())
	}

	err = packet.Write(conn, PLoginOk, packet.LoginOkPacket{})
	if err != nil {
		return nil, LoginError("write LoginOkPacket error")
	}
	log.Println("send a LoginOkPacket:", buf)

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
