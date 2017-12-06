package binutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func readNextString(b []byte, p *int) (res string) {
	strlen := int(binary.BigEndian.Uint32(b[*p : *p+4]))
	res = string(b[*p+4 : *p+strlen+4])
	*p += strlen + 4
	return
}

// Unmarashall takes a pointer to struct and try to fill the fields
// with data from the byte stream. s must be a pointer to the variable
func Unmarashall(b []byte, s interface{}) {
	pos := 0
	sType := reflect.ValueOf(s).Elem()

	if sType.Kind() == reflect.Struct {
		recurseStruct(b, sType, &pos)
	} else if sType.Kind() == reflect.String {
		s := s.(*string)
		*s = readNextString(b, &pos)
	} else {
		re := bytes.NewReader(b)
		err := binary.Read(re, binary.BigEndian, s)
		if err != nil {
			panic("Expecting struct instead of other types")
		}
	}
}

func readArray(b []byte, baseType reflect.Type, size int, p *int) (v reflect.Value) {
	v = reflect.New(reflect.ArrayOf(size, baseType))
	valueInterface := v.Interface()
	switch baseType.Kind() {
	case reflect.String:
		for i := 0; i < size; i++ {
			v.Elem().Index(i).Set(reflect.ValueOf(readNextString(b, p)))
		}
		v = reflect.Indirect(v)
	case reflect.Struct:
	default:
		re := bytes.NewReader(b[*p:])
		err := binary.Read(re, binary.BigEndian, valueInterface)
		v = reflect.Indirect(v)
		if err != nil {
			panic(fmt.Sprintf("Does not support type: %v", err))
		}
	}
	return
}

func recurseStruct(b []byte, v reflect.Value, p *int) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Bool:
				field.SetBool(b[*p] != 0)
				*p++
			case reflect.Uint16:
				field.SetUint(uint64(binary.BigEndian.Uint16(b[*p:])))
				*p += 2
			case reflect.Uint32:
				field.SetUint(uint64(binary.BigEndian.Uint32(b[*p:])))
				*p += 4
			case reflect.Uint64:
				field.SetUint(binary.BigEndian.Uint64(b[*p:]))
				*p += 8
			case reflect.String:
				field.SetString(readNextString(b, p))
			case reflect.Array:
				capacity := field.Cap()
				uType := field.Type().Elem().Kind()
				if uType == reflect.Uint8 {
					tmpSlice := b[*p : *p+capacity]
					reflect.Copy(field, reflect.ValueOf(tmpSlice))
				} else {
					field.Set(readArray(b, field.Type().Elem(), capacity, p))
				}
				*p += capacity * int(field.Type().Elem().Size())
			case reflect.Struct:
				recurseStruct(b, field, p)
			default:
				panic("Does not support other types!")
			}
		}
	}
}
