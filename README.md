# BinUtils - turns network protocol binary stream into Go struct
BinUtils helps to convert `[]byte` array captured from binary protocl connections into struct, just like what you have in languages like C++

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
bStream := CreateByteStream(binData)
bStream.ReadAsStruct(&t)

fmt.printf("A:%v B:%v C:%v", t.FieldA, t.FieldB, t.FieldC)
```

You can also reads a single variable:
```go
str := bStream.ReadNextString()
```

### Limitations
Currently the library supports `uint16`, `uint32`, `uint64` and `string` as datatypes for the field, plus nested struct (both named and anonymous). Support for slice(array) is coming soon!