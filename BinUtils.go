package binutils

import (
	"encoding/binary"
	"reflect"
)

// ByteStream is the container for the byte stream to be converted
type ByteStream struct {
	data []byte
	pos  int
}

// CreateByteStream is the constructor for the ByteStream struct
// takes the byte array as the only argument
func CreateByteStream(b []byte) ByteStream {
	return ByteStream{data: b}
}

// ReadNextUINT16 reads the next 2 bytes
// In the ByteStream as uint64
func (b *ByteStream) ReadNextUINT16() uint16 {
	defer func() { b.pos += 2 }()
	return binary.BigEndian.Uint16(b.data[b.pos:])
}

// ReadNextUINT32 reads the next 4 bytes in the stream as uint32
func (b *ByteStream) ReadNextUINT32() uint32 {
	defer func() { b.pos += 4 }()
	return binary.BigEndian.Uint32(b.data[b.pos:])
}

// ReadNextUINT64 reads the next 8 bytes in the stream as uint64
func (b *ByteStream) ReadNextUINT64() uint64 {
	defer func() { b.pos += 8 }()
	return binary.BigEndian.Uint64(b.data[b.pos:])
}

// ReadNextString reads next field as string, assuming the first 2 bytes is the length of the string
func (b *ByteStream) ReadNextString() string {
	strlen := int(binary.BigEndian.Uint32(b.data[b.pos : b.pos+4]))
	defer func() { b.pos += strlen + 4 }()
	return string(b.data[b.pos+4 : b.pos+strlen+4])
}

// ReadAsStruct takes a pointer to struct and try to fill the fields
// with data from the byte stream
func (b *ByteStream) ReadAsStruct(s interface{}) {
	stMeta := reflect.ValueOf(s).Elem()

	for i := 0; i < stMeta.NumField(); i++ {
		field := stMeta.Field(i)
		if field.CanSet() {
			switch field.Type().String() {
			case "uint16":
				field.SetUint(uint64(b.ReadNextUINT16()))
			case "uint32":
				field.SetUint(uint64(b.ReadNextUINT32()))
			case "uint64":
				field.SetUint(b.ReadNextUINT64())
			case "string":
				field.SetString(b.ReadNextString())
			case "[]byte":
				capacity := field.Cap()
				field.SetBytes(b.data[b.pos : b.pos+capacity])
				b.pos += capacity
			default:
				panic("Does not support other types!")
			}
		}
	}
}
