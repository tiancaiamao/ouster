package packet

import (
	"errors"
	"reflect"
	// "fmt"
)

//----------------------------------------------- write-out struct fields with packet writer.
func Pack(tos uint16, tbl interface{}, writer *Packet) []byte {
	if writer == nil {
		writer = Writer()
	}

	v := reflect.ValueOf(tbl)

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		v = v.Elem()
	}
	count := v.NumField()

	// write code
	if tos != 0 {
		writer.WriteU16(uint16(tos))
	}

	for i := 0; i < count; i++ {
		f := v.Field(i)
		if _is_primitive(f) {
			_write_primitive(f, writer)
		} else {
			switch f.Type().Kind() {
			case reflect.Slice, reflect.Array:
				writer.WriteU16(uint16(f.Len()))
				for a := 0; a < f.Len(); a++ {
					if _is_primitive(f.Index(a)) {
						_write_primitive(f.Index(a), writer)
					} else {
						elem := f.Index(a).Interface()
						Pack(0, elem, writer)
					}
				}
			}
		}
	}

	return writer.Data()
}

func Unpack(tbl interface{}, reader *Packet) (interface{}, error) {
	if reader == nil {
		return nil, errors.New("error parameter: unpack receive a nil reader!")
	}

	// fmt.Println("type of input:", reflect.TypeOf(tbl).Name())
	v := reflect.ValueOf(tbl)

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		v = v.Elem()
	}
	count := v.NumField()
	// fmt.Println("count of input is:", count)

	var err error
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if _is_primitive(f) {			
			err = _read_primitive(f, reader)
			if err != nil {
				return nil, err
			}
		} else {
			switch f.Type().Kind() {
			case reflect.Slice, reflect.Array:
				ui16, err := reader.ReadU16()
				for a := 0; a < int(ui16); a++ {
					if _is_primitive(f.Index(a)) {
						err = _read_primitive(f.Index(a), reader)
						if err != nil {
							return nil, err
						}
					} else {
						elem := f.Index(a).Interface()
						_, err := Unpack(elem, reader)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	
	return tbl, err
}

//----------------------------------------------- test whether the field is primitive type
func _is_primitive(f reflect.Value) bool {
	switch f.Type().Kind() {
	case reflect.Bool,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

//----------------------------------------------- write a primitive field
func _write_primitive(f reflect.Value, writer *Packet) {
	switch f.Type().Kind() {
	case reflect.Bool:
		writer.WriteBool(f.Interface().(bool))
	case reflect.Uint8:
		writer.WriteByte(f.Interface().(byte))
	case reflect.Uint16:
		writer.WriteU16(f.Interface().(uint16))
	case reflect.Uint32:
		writer.WriteU32(f.Interface().(uint32))
	case reflect.Uint64:
		writer.WriteU64(f.Interface().(uint64))

	case reflect.Int:
		writer.WriteU32(uint32(f.Interface().(int)))
	case reflect.Int8:
		writer.WriteByte(byte(f.Interface().(int8)))
	case reflect.Int16:
		writer.WriteU16(uint16(f.Interface().(int16)))
	case reflect.Int32:
		writer.WriteU32(uint32(f.Interface().(int32)))
	case reflect.Int64:
		writer.WriteU64(uint64(f.Interface().(int64)))

	case reflect.Float32:
		writer.WriteFloat32(f.Interface().(float32))

	case reflect.Float64:
		writer.WriteFloat64(f.Interface().(float64))

	case reflect.String:
		writer.WriteString(f.Interface().(string))
	}
}

func _read_primitive(f reflect.Value, reader *Packet) error {	
	switch f.Type().Kind() {
	case reflect.Bool:
		tmp, err := reader.ReadBool()
		f.SetBool(tmp)
		return err
	case reflect.Uint8:
		tmp, err := reader.ReadByte()
		f.SetUint(uint64(tmp))
		return err
	case reflect.Uint16:
		tmp, err := reader.ReadU16()
		f.SetUint(uint64(tmp))
		return err
	case reflect.Uint32:
		tmp, err := reader.ReadU32()
		f.SetUint(uint64(tmp))
		return err
	case reflect.Uint64:
		tmp, err := reader.ReadU64()
		f.SetUint(tmp)
		return err
	case reflect.Int:
		tmp, err := reader.ReadS32()
		f.SetInt(int64(tmp))
		return err
	case reflect.Int8:
		tmp, err := reader.ReadByte()
		f.SetInt(int64(tmp))
		return err
	case reflect.Int16:
		tmp, err := reader.ReadU16()
		f.SetInt(int64(tmp))
		return err
	case reflect.Int32:
		tmp, err := reader.ReadU32()
		f.SetUint(uint64(tmp))
		return err
	case reflect.Int64:
		tmp, err := reader.ReadU64()
		f.SetUint(uint64(tmp))
		return err
	case reflect.Float32:
		tmp, err := reader.ReadFloat32()
		f.SetFloat(float64(tmp))
		return err
	case reflect.Float64:
		tmp, err := reader.ReadFloat64()
		f.SetFloat(tmp)
		return err
	case reflect.String:
		tmp, err := reader.ReadString()
		f.SetString(tmp)
		return err
	}
	return errors.New("not support primitive in struct")
}
