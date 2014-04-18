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

type Fpoint struct {
	X float32
	Y float32
}



func main() {
	for _, v := range packet.PacketMap {
		Marshal(v, os.Stdout)
		Unmarshal(v, os.Stdout)

	}
	// Marshal(reflect.TypeOf(packet.SMove{}), os.Stdout)
	// Unmarshal(reflect.TypeOf(packet.SMove{}), os.Stdout)
		
	Marshal(reflect.TypeOf(Fpoint{}), os.Stdout)
	Unmarshal(reflect.TypeOf(Fpoint{}), os.Stdout)
	// PacketId()
	return
}

func Unmarshal(obj reflect.Type, writer io.Writer) {
	fmt.Fprintf(writer, `bool Unpack_%s(msgpack_object obj, struct %s *st) {`, obj.Name(), obj.Name())

	switch obj.Kind() {
	case reflect.Struct:
		fmt.Fprintf(writer, `
	if(obj.type != MSGPACK_OBJECT_ARRAY) {
		return false;
	}
`)

		num := obj.NumField()
		fmt.Fprintf(writer, `	if(obj.via.array.size != %d) {
		return false;
	}
`, num)

		for i := 0; i < num; i++ {
			fild := obj.Field(i)
			kind := fild.Type.Kind()
			if kind == reflect.Array || kind == reflect.Slice {
				fmt.Fprintf(writer, `	if(obj.via.array.ptr[%d].type != MSGPACK_OBJECT_ARRAY)
		return false;
`, i)

				fmt.Fprintf(writer, `	for(int i=0; i<obj.via.array.ptr[%d].via.array.size; i++) {
		Unpack_%s(obj.via.array.ptr[i], &st->Array[i]);			
	}`, i, fild.Type.Elem().Name())
			} else {
				fmt.Fprintf(writer, `	if(!Unpack_%s(obj.via.array.ptr[%d], &st->%s))
		return false;
`, fild.Type.Name(), i, fild.Name)
			}
		}
	case reflect.Map:
		fmt.Fprintf(writer, `
	if(obj.type != MSGPACK_OBJECT_MAP) {
		return false;
	}
	
	`)
	default:
		panic("unsupport!!!")
	}

	fmt.Fprintf(writer, `
	return true;
}

`)
}

func Marshal(obj reflect.Type, writer io.Writer) {

	fmt.Fprintf(writer, `void Pack_%s(struct %s *st, msgpack_packer *packer) {
`, obj.Name(), obj.Name())

	if obj.Kind() == reflect.Struct {
		num := obj.NumField()
		fmt.Fprintf(writer, `	msgpack_pack_array(packer, %d);
`, num)

		for i := 0; i < num; i++ {
			fild := obj.Field(i)
			if fild.Type.Kind() == reflect.Slice || fild.Type.Kind() == reflect.Array {
				fmt.Fprintf(writer, `	for (int i=0; i<st->size; i++) {
		Pack_%s(&st->Array[i], packer);
	}
`, fild.Type.Elem().Name())
			} else {
				fmt.Fprintf(writer, `	Pack_%s(&st->%s, packer);
`, fild.Type.Name(), fild.Name)
			}
		}
	}

	fmt.Fprintf(writer, `}

`)
}

func PacketId() {
	for k, v := range packet.PacketMap {
		// XXXPacket -> PXXX
		str := "P" + v.Name()[:len(v.Name())-6]
		fmt.Printf("#define %s %d\n", str, k)
	}
}
