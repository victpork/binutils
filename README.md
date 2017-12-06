# BinUtils - turns network protocol binary stream into Go struct
BinUtils helps to convert `[]byte` array captured from binary protocl connections into struct, just like what you have in languages like C/C++

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
}

...
var t TestStruct
binData := getByteArrayFromSomewhere()
binutils.Unmarshall(binData, &t)

fmt.printf("A:%v B:%v C:%v", t.FieldA, t.FieldB, t.FieldC)
```
### Design
BinUtils uses reflection to read in the structure of the datatype. It maybe slower, but saves you writing struct metadata along the code.

### Limitations
Currently the library supports `bool`, `uint16`, `uint32`, `uint64` and `string` as datatypes for the field, plus nested struct (both named and anonymous), fixed sized array with above data types. Support for slices, array of struct is coming soon!