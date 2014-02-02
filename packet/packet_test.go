package packet

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"reflect"
)

func TestPacketWriter(t *testing.T) {
	p := Writer()
	a := byte(0xFF)
	b := uint16(0xFF00)
	c := uint32(0xFF0000)
	d := uint32(0xFF000000)
	e := uint64(0xFF00000000000000)
	f32 := float32(1.0)
	f64 := float64(2.0)

	p.WriteBool(true)
	p.WriteByte(a)
	p.WriteU16(b)
	p.WriteU24(c)
	p.WriteU32(d)
	p.WriteU64(e)
	p.WriteFloat32(f32)
	p.WriteFloat64(f64)

	p.WriteString("hello world")
	p.WriteBytes([]byte("hello world"))

	reader := Reader(p.Data())

	BOOL, _ := reader.ReadBool()
	if BOOL != true {
		t.Error("packet readbool mismatch")
	}

	tmp, _ := reader.ReadByte()
	if a != tmp {
		t.Error("packet readbyte mismatch")
	}

	tmp1, _ := reader.ReadU16()
	if b != tmp1 {
		t.Error("packet readu16 mismatch")
	}

	tmp2, _ := reader.ReadU24()
	if c != tmp2 {
		t.Error("packet readu24 mismatch")
	}

	tmp3, _ := reader.ReadU32()
	if d != tmp3 {
		t.Error("packet readu32 mismatch")
	}

	tmp4, _ := reader.ReadU64()
	if e != tmp4 {
		t.Error("packet readu64 mismatch")
	}

	tmp5, _ := reader.ReadFloat32()
	if f32 != tmp5 {
		t.Error("packet readf32 mismatch")
	}

	tmp6, _ := reader.ReadFloat64()
	if f64 != tmp6 {
		t.Error("packet readf32 mismatch")
	}

	tmp100, _ := reader.ReadString()

	if "hello world" != tmp100 {
		t.Error("packet read string mistmatch")
	}

	tmp101, _ := reader.ReadBytes()

	fmt.Println(tmp101)
	if tmp101[0] != 'h' {
		t.Error("packet read bytes mistmatch")
	}

	_, err := reader.ReadByte()

	if err == nil {
		t.Error("overflow check failed")
	}
}

func BenchmarkPacketWriter(b *testing.B) {
	a := byte(0xFF)
	rand.Seed(time.Now().Unix())

	for i := 0; i < b.N; i++ {
		p := Writer()
		n := rand.Intn(128)
		for j := 0; j < n; j++ {
			p.WriteByte(a)
		}
	}
}

func TestUnpack(t *testing.T) {
	ci := CharactorInfoPacket{
		Name:  "test",
		Class: "class",
		Level: 63,
	}

	buf := Pack(PCharactorInfo, ci, Writer())
	t.Log(buf)
	result, err := Parse(buf)
	if err != nil {
		t.Fatal(err)
	}

	raw, ok := result.(*CharactorInfoPacket)
	if !ok {
		t.Fatal("not a charactorinfo packet")
	}
	if raw.Level != 63 {
		t.Fatal("parse error")
	}
}

func TestParse(t *testing.T) {
	// ci := LoginPacket{
	// 	Username: "genius",
	// 	Password: "0101001",
	// }

	// buf := Pack(PLogin, ci, Writer())
	// t.Log(buf)

	input := []byte{0, 0, 0, 19, 0, 1, 0, 6, 103, 101, 110, 105, 117, 115, 0, 7, 48, 49, 48, 49, 48, 48, 49}
	itf, err := Parse(input[4:])
	if err != nil {
		t.Fatal("parse error", err)
	}
	if _, ok := itf.(*LoginPacket); !ok {
		t.Fatal("not right packet")
	}
	t.Log(reflect.TypeOf(itf).Name())
}
