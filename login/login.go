package login

import (
	"encoding/binary"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"io"
	"net"
	"time"
	"log"
	"bytes"
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
		return nil, LoginError(err.Error())
	}
	
	log.Println("read packet data:", data)
	// packet解析
	p, err := packet.Parse(data)
	if err != nil {
		log.Println("parse packet error")
		return nil, LoginError(err.Error())
	}

	return p, nil
}

func Login(conn net.Conn) (*data.Player, error) {
	err := conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
	if err != nil {
		return nil, err
	}

	//------------------
	p, err := readPacket(conn)
	if err != nil {
		return nil, err
	}
	pkt, ok := p.(*packet.LoginPacket)
	if !ok {
		return nil, LoginError("expect a LoginPacket")
	}

	log.Println("get a LoginPacket...run here")
	// check username and ignore password ...
	charactor, err := loadUser(pkt.Username)

	writer := packet.Writer()
	writer.WriteU32(0)
	buf := packet.Pack(packet.PCharactorInfo, charactor, writer)
	binary.BigEndian.PutUint32(buf[:4], uint32(len(buf) - 4))
	
	// err = writeN(buf, conn)
	io.Copy(conn, bytes.NewReader(buf))
	if err != nil {
		return nil, LoginError("write LoginOkPacket error")
	}
	log.Println("send a CharactorInfoPacket:", buf)

	//------------------
	p, err = readPacket(conn)
	if err != nil {
		return nil, LoginError(err.Error())
	}
	selec, ok := p.(packet.SelectCharactorPacket)
	if !ok {
		return nil, LoginError("expect a SelectCharactorPacket")
	}
	
	// load player info ...
	player, err := loadCharactor(selec.Name)
	if err != nil {
		return nil, LoginError(err.Error())
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

func loadCharactor(charactor string) (*data.Player, error) {
	if charactor == "test" {
		return &data.Player{
			Name:  "test",
			Class: data.PlayerClass(data.BRUTE),
			Level: 1,
		}, nil
	} else {
		return nil, LoginError("no such charactor!")
	}
}
