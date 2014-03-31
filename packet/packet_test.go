package packet

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestPrimitive(t *testing.T) {
	m := make(map[string]int)
	m["test"] = 234
	var testCase = [...]interface{}{
		8, -1, uint8(3), int32(43), int64(1435435345), float32(32.342), float64(-3321.34), [...]int{2, 3, 1}, m,
	}

	for i, v := range testCase {
		PacketMap[PTest] = reflect.TypeOf(v)

		buf := &bytes.Buffer{}
		enc := NewEncoder(buf)
		err := enc.Encode(PTest, v)
		if err != nil {
			t.Fatal(fmt.Sprintf("the %dth testCase failed...%s", i, err.Error()))
		}

		//		t.Log(buf.Bytes())
		var pkt Packet
		pkt.Id = 2345
		pkt.Obj = nil
		dec := NewDecoder(buf)
		err = dec.Decode(&pkt)
		if err != nil {
			t.Fatal(fmt.Sprintf("decode %dth testCase failed...%s", i, err.Error()))
		}
		if pkt.Id != PTest {
			t.Fatal("id error")
		}
		if !reflect.DeepEqual(v, pkt.Obj) {
			t.Fatal("Obj error")
		}
	}
}

func TestStruct(t *testing.T) {
	var pkt Packet
	pkt.Id = PLogin
	obj := LoginPacket{
		Username: "genius",
		Password: "0101001",
	}

	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	err := enc.Encode(PLogin, obj)
	if err != nil {
		t.Fatal(err)
	}

	pkt.Obj = nil
	dec := NewDecoder(buf)
	err = dec.Decode(&pkt)
	if err != nil {
		t.Fatal("Decode obj error:", err)
	}
	v, ok := pkt.Obj.(LoginPacket)
	if !ok {
		t.Fatal("Decode type error")
	}
	if v.Username != "genius" || v.Password != "0101001" {
		t.Fatal("username or password error")
	}
}

func TestExport(t *testing.T) {
	buf := &bytes.Buffer{}
	login := LoginPacket{
		Username: "genius",
		Password: "0101001",
	}
	err := Write(buf, PLogin, login)
	if err != nil {
		t.Fatal(err)
	}

	pkt, err := Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pkt)
}

func TestType(t *testing.T) {
	buf := &bytes.Buffer{}
	for k, v := range PacketMap {
		if k != PMax {
			err := Write(buf, k, reflect.Zero(v).Interface())
			if err != nil {
				t.Fatal(err)
			}	
		}		
	}
	
	for k, v := range PacketMap {
		if k != PMax {
			pkt, err := Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			if reflect.TypeOf(pkt) != v {
				t.Fatal("insistent for ", k)
			}			
		}
	}
}