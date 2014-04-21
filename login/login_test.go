package login

import (
	"github.com/tiancaiamao/ouster/config"
	"github.com/tiancaiamao/ouster/packet"
	"net"
	"testing"
)

func TestLogin(t *testing.T) {
	l, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	//------- client ----------
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1"+config.ServerPort)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		packet.Write(conn, packet.PLogin, packet.LoginPacket{
			Username: "genius",
			Password: "0101001",
		})

		info, err := packet.Read(conn)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := info.(packet.CharactorInfoPacket); !ok {
			t.Fatal("need a CharactorInfoPacket")
		}

		packet.Write(conn, packet.PSelectCharactor, packet.SelectCharactorPacket{
			Which: 0,
		})

		loginOk, err := packet.Read(conn)
		if _, ok := loginOk.(packet.LoginOkPacket); !ok {
			t.Fatal("need a LoginOkPacket")
		}
	}()
	//----------------------

	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}

	go func(c net.Conn) {
		Login(c)
	}(conn)
}
