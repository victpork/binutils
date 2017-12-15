[![Build Status](https://api.travis-ci.org/mkishere/binutils.svg)](http://travis-ci.org/mkishere/binutils) [![GoDoc](https://godoc.org/github.com/mkishere/binutils?status.svg)](http://godoc.org/github.com/mkishere/binutils)


# BinUtils - turns network protocol binary stream into Go struct
BinUtils helps to convert `[]byte` array captured from binary protocol connections into struct, just like what you have in languages like C/C++.

### Download and install
In command prompt:
```
go get github.com/mkishere/binutils
```
### Run test cases
```
go test
```

### Usage
In your code:
```go
import (
    ...
    "github.com/mkishere/binutils"
)

type TestStruct struct {
    FieldA uint16
    FieldB uint32
    FieldC string
    FieldD [3]uint64
    FieldE struct {
        FieldE1 int16
        FieldE2 int32
        FieldE3 string
    }
}

...
var t TestStruct
binData := getByteArrayFromSomewhere()
binutils.Unmarshal(binData, &t)

fmt.printf("A:%v B:%v C:%v", t.FieldA, t.FieldB, t.FieldC)
```

Marshalling is also easy:
```go
import (
    ...
    "github.com/mkishere/binutils"
)

var t SomeStruct
bindata := Marshal(&t) //bindata is []byte
```
### Design
BinUtils uses reflection to read in the structure of the datatype. It maybe slower, but saves you writing struct metadata along the code.

The library uses standard network binary encoding: For numeric variables all are taken as Big Endian. For strings the length of the string will be represented as a 4-byte uint32 and ASCII code of actual string follows. E.g. `"Hello"` would become `[]byte{0, 0, 0, 5, 72, 101, 108, 108, 111}`. Slices are represented similar to `string`, with 4-byte `uint32` _size_ as header and elements comes after.

The library is for Go to communicate with pre-existing binary protocol (e.g. SSH), and efficiency is not in mind. If you are working on brand-new client/server protocol, work with a modern binary protocol like [Protocol Buffers](https://github.com/google/protobuf) or [gob](https://golang.org/pkg/encoding/gob/).

### Support datatypes
Currently the library supports `bool`, `int`/`uint` (8/16/32/64) and `string` as datatypes for the field, plus nested struct (both named and anonymous), fixed sized array with above data types. When defining struct, please refrain from using machine dependant types `int`/`uint`. Use lenght-specific types instead e.g `uint32`

