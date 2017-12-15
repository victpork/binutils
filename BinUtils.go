// Package binutils provides ways of encoding and decoding traditional C/C++
// data structs for use in binary network protocols.
//
// Example:
//  import (
//      ...
//      "github.com/mkishere/binutils"
//  )
//
//  type TestStruct struct {
//      FieldA uint16
//      FieldB uint32
//      FieldC string
//      FieldD [3]uint64
//      FieldE struct {
//          FieldE1 int16
//          FieldE2 int32
//          FieldE3 string
//      }
//  }
//
//  ...
//  var t TestStruct
//  binData := getByteArrayFromSomewhere()
//  binutils.Unmarshal(binData, &t)
//
//  fmt.printf("A:%v B:%v C:%v", t.FieldA, t.FieldB, t.FieldC)
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

// Unmarshal takes a pointer to struct and try to fill the fields
// with data from the byte stream. s must be a pointer to the variable
func Unmarshal(b []byte, s interface{}) {
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

// Marshal takes a pointer to struct and try to iterate the fields
// and convert data format into the byte array. Returns the concat byte array
func Marshal(s interface{}) (b []byte) {
	switch v := s.(type) {
	case *string:
		b = make([]byte, 4+len(*v))
		binary.BigEndian.PutUint32(b, uint32(len(*v)))
		copy(b[4:], []byte(*v))
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
			b = make([]byte, sizeOfType(s))
			recurseWriteStruct(uValue, b[:0])
		} else {
			panic("Unsupported type " + uValue.Kind().String())
		}
	}
	return
}

// sizeOfType returns the number of bytes required to represent the variable
func sizeOfType(i interface{}) (size int) {
	switch i.(type) {
	case *bool, *int8, *uint8, bool, int8, uint8:
		size = 1
	case []bool:
		size = len(i.([]bool))
	case *uint16, *int16, uint16, int16:
		size = 2
	case []uint16, []int16:
		size = 2 * len(i.([]uint16))
	case *uint32, *int32, uint32, int32:
		size = 4
	case *uint64, *int64, uint64, int64:
		size = 8
	case *string, string:
		size = 4 + len(i.(string))
	default:
		v := reflect.ValueOf(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Slice:
			if v.Len() > 0 {
				size += 4
				for i := 0; i < v.Len(); i++ {
					size += sizeOfType(v.Index(i).Interface())
				}
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				size += sizeOfType(field.Interface())
			}
		case reflect.Array:
			underlyingType := v.Type().Elem()
			if underlyingType.Kind() == reflect.Struct {
				for i := 0; i < v.Len(); i++ {
					size += sizeOfType(v.Index(i).Interface())
				}
			} else if underlyingType.Kind() == reflect.Array || underlyingType.Kind() == reflect.Slice {
				panic("Does not support multidimension array/slices")
			} else {
				size = sizeOfType(v.Index(0).Interface()) * v.Len()
			}
		}
	}
	return
}

func recurseWriteStruct(v reflect.Value, b []byte) int {
	for i := 0; i < v.NumField(); i++ {
		pos := len(b)
		field := v.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Bool:
				var bit byte
				if v.Bool() {
					bit = 1
				}
				b = append(b, bit)
			case reflect.Uint8:
				b = append(b, uint8(field.Uint()))
			case reflect.Uint16:
				binary.BigEndian.PutUint16(b[pos:], uint16(field.Uint()))
			case reflect.Uint32:
				b = b[:pos+4]
				binary.BigEndian.PutUint32(b[pos:], uint32(field.Uint()))
			case reflect.Uint64:
				b = b[:pos+8]
				binary.BigEndian.PutUint64(b[pos:], field.Uint())
			case reflect.String:
				strLen := field.Len()
				b = b[:pos+4+strLen]
				binary.BigEndian.PutUint32(b[pos:], uint32(strLen))
				copy(b[pos+4:], []byte(field.String()))
			case reflect.Struct:
				pos += recurseWriteStruct(field, b[pos:])
				b = b[:pos]
			case reflect.Slice:
				b = b[:pos+4]
				binary.BigEndian.PutUint32(b[pos:pos+4], uint32(field.Len()))
				pos += 4
				fallthrough
			case reflect.Array:
				pos += writeArrayToByte(b[pos:pos], field)
				b = b[pos:pos]
			}
		}
	}
	return len(b)
}

func writeArrayToByte(b []byte, array reflect.Value) (pos int) {
	size := array.Len()
	switch array.Type().Elem().Kind() {
	case reflect.Struct:
		for i := 0; i < size; i++ {
			pos += recurseWriteStruct(array.Index(i), b[pos:])
			b = b[:pos]
		}
	case reflect.Array, reflect.Slice:
		panic("Does not support multidimensional array/slice")
	default:
		pArray := array.Interface()
		buf := bytes.NewBuffer(make([]byte, size)[:0])
		err := binary.Write(buf, binary.BigEndian, pArray)
		if err != nil {
			panic("Could not write to byte stream")
		}
		pos = buf.Len()
		b = b[:pos]
		copy(b, buf.Bytes())
	}
	return
}

func readArrayFromByte(b []byte, baseType reflect.Type, size int, p *int) (v reflect.Value) {
	v = reflect.New(reflect.ArrayOf(size, baseType))
	valueInterface := v.Interface()
	switch baseType.Kind() {
	case reflect.String:
		for i := 0; i < size; i++ {
			v.Elem().Index(i).Set(reflect.ValueOf(readNextString(b, p)))
		}
	case reflect.Struct:
		for i := 0; i < size; i++ {
			recurseReadStruct(b, v.Elem().Index(i), p)
		}
	case reflect.Array, reflect.Slice:
		panic("Does not support multidimensional array/slice")
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
					field.Set(readArrayFromByte(b, field.Type().Elem(), capacity, p))
				}
				*p += capacity * int(field.Type().Elem().Size())
			case reflect.Struct:
				recurseReadStruct(b, field, p)
			case reflect.Slice:
				capacity := int(binary.BigEndian.Uint32(b[*p : *p+4]))
				*p += 4
				// Todo: This line is so slow...
				field.Set(readArrayFromByte(b, field.Type().Elem(), capacity, p).Slice(0, capacity))
				*p += capacity * int(field.Type().Elem().Size())
			default:
				panic("Unsupported type " + field.Kind().String())
			}
		}
	}
}
