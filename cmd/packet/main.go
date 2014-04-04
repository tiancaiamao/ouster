package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster/packet"
	"io"
	"os"
	"reflect"
)

type Test struct {
	a int
}

type FPoint struct {
	X float32
	Y float32
}

func main() {
	for _, v := range packet.PacketMap {
		// fmt.Println(k, "----", v)
		Marshal(v, os.Stdout)
		Unmarshal(v, os.Stdout)
	}
	Unmarshal(reflect.TypeOf(FPoint{}), os.Stdout)
	return
}

func Unmarshal(obj reflect.Type, writer io.Writer) {
	fmt.Fprintf(writer, `bool Unpack_%s(msgpack_object obj, struct %s *st) {
	if(obj.type != MSGPACK_OBJECT_ARRAY) {
		return false;
	}
`, obj.Name(), obj.Name())

	num := obj.NumField()
	fmt.Fprintf(writer, `	if(obj.via.array.size != %d) {
		return false;
	}
`, num)

	for i := 0; i < num; i++ {
		fild := obj.Field(i)
		fmt.Fprintf(writer, `	if(!Unpack_%s(obj.via.array.ptr[%d], &st->%s))
		return false;
`, fild.Type.Name(), i, fild.Name)
	}

	fmt.Fprintf(writer, `
}

`)
}

func Marshal(obj reflect.Type, writer io.Writer) {
	fmt.Fprintf(writer, `void Pack_%s(struct %s *st, msgpack_packer *packer) {
`, obj.Name(), obj.Name())

	num := obj.NumField()
	fmt.Fprintf(writer, `	msgpack_pack_array(packer, %d);
`, num)

	for i := 0; i < num; i++ {
		fild := obj.Field(i)
		fmt.Fprintf(writer, `	Pack_%s(&st->%s, packer);
`, fild.Type.Name(), fild.Name)
	}

	fmt.Fprintf(writer, `}

`)
}
