package packet

import (
	"bytes"
	"fmt"
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

func (enc *Encoder) Encode(pkt *Packet) error {
	if pkt.Id >= PMax {
		return PacketError("Encode: packet id out ranged")
	}

	err := enc.Encoder.Encode(pkt.Id)
	if err != nil {
		return err
	}

	ti := PacketMap[pkt.Id]
	if ti != reflect.TypeOf(pkt.Obj) {
		return PacketError("Encode: inconsistent of packet's id and Obj\n")
	}

	err = enc.Encoder.Encode(pkt.Obj)
	return err
}

func Marshal(pkt *Packet) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	err := enc.Encode(pkt)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, pkt *Packet) error {
	dec := NewDecoder(bytes.NewBuffer(data))
	return dec.Decode(pkt)
}
