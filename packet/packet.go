package packet

import (
	"fmt"
	"github.com/tiancaiamao/ouster"
	"github.com/ugorji/go/codec"
	"io"
	"reflect"
)

var mh codec.MsgpackHandle

type PacketError string

func (err PacketError) Error() string {
	return string(err)
}

type Packet struct {
	Id  uint16
	Obj interface{}
}

type Decoder struct {
	*codec.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{codec.NewDecoder(r, &mh)}
}

func (dec *Decoder) Decode(pkt *Packet) error {
	err := dec.Decoder.Decode(&pkt.Id)
	if err != nil || pkt.Id >= PMax {
		return PacketError(fmt.Sprintf("Decode: decode packet id error: %s", err.Error()))
	}

	ti, ok := PacketMap[pkt.Id]
	if !ok {
		return PacketError("Decode: invalid packet id")
	}

	// TODO: use a object pool here to avoid frequence allocation
	v := reflect.New(ti)
	err = dec.Decoder.Decode(v)
	pkt.Obj = v.Elem().Interface()
	return err
}

type Encoder struct {
	*codec.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{codec.NewEncoder(w, &mh)}
}

func (enc *Encoder) Encode(id uint16, obj interface{}) error {
	if id >= PMax {
		return PacketError("Encode: packet id out ranged")
	}

	err := enc.Encoder.Encode(id)
	if err != nil {
		return ouster.NewError("Encode id error:" + err.Error())
	}

	ti := PacketMap[id]
	if ti != reflect.TypeOf(obj) {
		return PacketError("Encode: inconsistent of packet's id and Obj: id is " + ti.String() + " , but obj is " + reflect.TypeOf(obj).String())
	}

	err = enc.Encoder.Encode(obj)
	if err != nil {
		return ouster.NewError("Encode obj error:" + err.Error())
	}
	return nil
}
