package login

import (
	"encoding/binary"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/packet"
	"io"
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

	size := binary.BigEndian.Uint16(header[:])
	data := make([]byte, size)

	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, LoginError(err.Error())
	}

	// packet解析
	p, err := packet.Parse(data)
	if err != nil {
		return nil, LoginError(err.Error())
	}

	return p, nil
}

func writeN(buf []byte, conn net.Conn) error {
	written := 0
	for written != len(buf) {
		n, err := conn.Write(buf[:written])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
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
	pkt, ok := p.(packet.LoginPacket)
	if !ok {
		return nil, LoginError("expect a LoginPacket")
	}

	// check username and ignore password ...
	charactor, err := loadUser(pkt.Username)

	buf := packet.Pack(-1, charactor, packet.Writer())
	err = writeN(buf, conn)
	if err != nil {
		return nil, LoginError("write LoginOkPacket error")
	}

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
