package binutils

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func readNextString(b []byte, p *int) (res string) {
	strlen := int(binary.BigEndian.Uint32(b[*p : *p+4]))
	res = string(b[*p+4 : *p+strlen+4])
	*p += strlen + 4
	return
}

// Unmarshall takes a pointer to struct and try to fill the fields
// with data from the byte stream. s must be a pointer to the variable
func Unmarshall(b []byte, s interface{}) {
	pos := 0

	switch s.(type) {
	case *string:
		s := s.(*string)
		*s = readNextString(b, &pos)
	case *uint8, *uint16, *uint32, *uint64, *int8, *int16, *int32, *int64, *[]uint8, *[]uint16, *[]uint32, *[]uint64:
		re := bytes.NewReader(b)
		err := binary.Read(re, binary.BigEndian, s)
		if err != nil {
			panic("Could not read data from binary stream")
		}
	default:
		sType := reflect.ValueOf(s).Elem()
		if sType.Kind() == reflect.Struct {
			recurseReadStruct(b, sType, &pos)
		} else {
			panic("Unsupported type " + sType.Kind().String())
		}
	}
}

// Marshall takes a pointer to struct and try to iterate the fields
// and convert data format into the byte array. Returns the concat byte array
func Marshall(s interface{}) (b []byte) {
	switch v := s.(type) {
	case *string:
		lenByteArr := make([]byte, 4)
		binary.BigEndian.PutUint32(lenByteArr, uint32(len(*v)))
		b = append(lenByteArr, *v...)
	case *uint8, *uint16, *uint32, *uint64, *int8, *int16, *int32, *int64, *[]uint8, *[]uint16, *[]uint32, *[]uint64:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, s)
		if err != nil {
			panic("Could not read data from binary stream")
		}
		b = buf.Bytes()
	default:
		uValue := reflect.ValueOf(s).Elem()
		if uValue.Kind() == reflect.Struct {
			b = make([]byte, sizeOfStruct(uValue))
			b = recurseWriteStruct(uValue, b)
		} else {
			panic("Unsupported type " + uValue.Kind().String())
		}
	}
	return
}

// sizeOfStruct returns the number of bytes required to represent the struct(variable)
func sizeOfStruct(v reflect.Value) (size int) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Bool, reflect.Uint8, reflect.Int8:
				size++
			case reflect.Uint16, reflect.Int16:
				size += 2
			case reflect.Uint32, reflect.Int32:
				size += 4
			case reflect.Uint64, reflect.Int64:
				size += 8
			case reflect.Struct:
				size += sizeOfStruct(field.Elem())
			case reflect.String:
				size += 4 + field.Len()
			case reflect.Array:
				size += sizeOfStruct(field.Elem()) * field.Len()
			case reflect.Slice:
				size += 4 + sizeOfStruct(field.Elem())*field.Len()
			}
		}
	}
	return
}

func recurseWriteStruct(v reflect.Value, b []byte) (res []byte) {
	// TODO: Optimize memory usage

	copy(res, b)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Bool:
				var bit byte
				if v.Bool() {
					bit = 1
				}
				res = append(res, bit)
			case reflect.Uint8:
				res = append(res, uint8(field.Uint()))
			case reflect.Uint16:
				buf := make([]byte, 2)
				binary.BigEndian.PutUint16(buf, uint16(field.Uint()))
				res = append(res, buf...)
			case reflect.Uint32:
				buf := make([]byte, 4)
				binary.BigEndian.PutUint32(buf, uint32(field.Uint()))
				res = append(res, buf...)
			case reflect.Uint64:
				buf := make([]byte, 8)
				binary.BigEndian.PutUint64(buf, field.Uint())
				res = append(res, buf...)
			case reflect.String:
				lenByteArr := make([]byte, 4)
				binary.BigEndian.PutUint32(lenByteArr, uint32(field.Len()))
				res = append(res, lenByteArr...)
				res = append(res, []byte(field.String())...)
			case reflect.Struct:
				res = append(recurseWriteStruct(field, res))
			}
		}
	}
	return
}

func readArray(b []byte, baseType reflect.Type, size int, p *int) (v reflect.Value) {
	v = reflect.New(reflect.ArrayOf(size, baseType))
	valueInterface := v.Interface()
	switch baseType.Kind() {
	case reflect.String:
		for i := 0; i < size; i++ {
			v.Elem().Index(i).Set(reflect.ValueOf(readNextString(b, p)))
		}
	case reflect.Struct:
		//Todo: support array of struct
		for i := 0; i < size; i++ {
			recurseReadStruct(b, v.Elem().Index(i), p)
		}
	default:
		re := bytes.NewReader(b[*p:])
		err := binary.Read(re, binary.BigEndian, valueInterface)

		if err != nil {
			panic("Could not read data from binary stream")
		}
	}
	v = reflect.Indirect(v)
	return
}

func recurseReadStruct(b []byte, v reflect.Value, p *int) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Bool:
				field.SetBool(b[*p] != 0)
				*p++
			case reflect.Int8:
				field.SetInt(int64(b[*p]))
				*p++
			case reflect.Uint8:
				field.SetUint(uint64(b[*p]))
				*p++
			case reflect.Int16:
				field.SetInt(int64(binary.BigEndian.Uint16(b[*p:])))
				*p += 2
			case reflect.Uint16:
				field.SetUint(uint64(binary.BigEndian.Uint16(b[*p:])))
				*p += 2
			case reflect.Int32:
				field.SetInt(int64(binary.BigEndian.Uint32(b[*p:])))
				*p += 4
			case reflect.Uint32:
				field.SetUint(uint64(binary.BigEndian.Uint32(b[*p:])))
				*p += 4
			case reflect.Int64:
				field.SetInt(int64(binary.BigEndian.Uint64(b[*p:])))
				*p += 8
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
				recurseReadStruct(b, field, p)
			case reflect.Slice:
				capacity := int(binary.BigEndian.Uint32(b[*p : *p+4]))
				*p += 4
				// Todo: This line is so slow...
				field.Set(readArray(b, field.Type().Elem(), capacity, p).Slice(0, capacity))
				*p += capacity * int(field.Type().Elem().Size())
			default:
				panic("Unsupported type " + field.Kind().String())
			}
		}
	}
}
