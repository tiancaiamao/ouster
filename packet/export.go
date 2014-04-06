package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tiancaiamao/ouster"
	"io"
	"log"
	"reflect"
)

// Read a packet from a io.Reader, the format of a packet consistent of
// a BigEndian uint32 indicate the packet size, a Msgpack uint16 id indicate the
// packet type, and a Msgpack array which is the packet content
func Read(conn io.Reader) (interface{}, error) {
	var header [4]byte
	_, err := io.ReadFull(conn, header[:])
	if err != nil {
		return nil, ouster.NewError(err.Error())
	}

	log.Println("read packet head:", header)

	size := binary.BigEndian.Uint32(header[:])
	log.Println("read a packet, header size:", size)
	data := make([]byte, size)

	_, err = io.ReadFull(conn, data)
	if err != nil {
		log.Println("ReadFull error")
		return nil, ouster.NewError(err.Error())
	}

	log.Println("read packet data:", data)

	p, err := parse(data)
	if err != nil {
		log.Println("parse packet error")
		return nil, ouster.NewError(err.Error())
	}

	return p, nil
}

// data = id + struct of the packet
func parse(data []byte) (interface{}, error) {
	buf := bytes.NewBuffer(data)
	dec := NewDecoder(buf)
	var pkt Packet
	err := dec.Decode(&pkt)
	if err != nil {
		return nil, ouster.NewError(err.Error())
	}
	log.Println("read a packet:", pkt.Id, pkt.Obj)

	// translate map to XXXPacket
	if reflect.TypeOf(pkt.Obj).Kind() == reflect.Map {
		tp := PacketMap[pkt.Id]
		if !reflect.TypeOf(pkt.Obj).ConvertibleTo(tp) {
			return nil, ouster.NewError("in consistent PacketId and content's Type")
		}
		ret := reflect.ValueOf(pkt.Obj).Convert(tp).Interface()
		// log.Println("convert to ", tp.Name())
		return ret, nil
	} else {
		return pkt.Obj, nil
	}
}

func Write(conn io.Writer, id uint16, obj interface{}) error {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	err := enc.Encode(id, obj)
	if err != nil {
		return ouster.NewError("encode error: " + err.Error())
	}

	len := make([]byte, 4) //TODO
	binary.BigEndian.PutUint32(len, uint32(buf.Len()))

	_, err = conn.Write(len)
	if err != nil {
		return ouster.NewError("write to io.Writer error: " + err.Error())
	}

	_, err = io.Copy(conn, buf)
	if err != nil {
		return ouster.NewError(fmt.Sprintf("write %d error", id))
	}
	// log.Println("send a:", id, buf)
	return nil
}
