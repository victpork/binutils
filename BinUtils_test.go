package binutils

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMarshallString(t *testing.T) {
	testStr := "HelloWorld"
	expected := []byte{0, 0, 0, 10, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	result := Marshall(&testStr)
	if !bytes.Equal(result, expected) {
		t.Errorf("Term does not match, expected '%v' now '%v'", expected, result)
	}
}

func TestMarshallStruct(t *testing.T) {
	type sshptyRequest struct {
		Term    string
		Width   uint32
		Height  uint32
		PWidth  uint32
		PHeight uint32
	}

	original := sshptyRequest{Term: "xterm", Width: 80, Height: 24}
	expected := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0,
		0, 0}
	result := Marshall(&original)
	if !bytes.Equal(result, expected) {
		t.Errorf("Term does not match, expected '%v' now '%v'", expected, result)
	}
}

func TestSSHPtyRequest(t *testing.T) {

	type sshptyRequest struct {
		Term     string
		Width    uint32
		Height   uint32
		PWidth   uint32
		PHeight  uint32
		TermMode []byte
	}

	var result sshptyRequest
	expected := sshptyRequest{Term: "xterm", Width: 80, Height: 24, TermMode: []byte{3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}}
	b := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 21, 3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}
	Unmarshall(b, &result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Term does not match, expected '%v' now '%v'", expected, result)
	}
}

func TestStructWithArrayFields(t *testing.T) {
	type arrayStruct struct {
		FieldA [3]byte
		FieldB [3]uint16
		FieldC [3]string
	}

	b := []byte{10, 17, 10, 0, 0, 0, 1, 1, 2, 0, 0, 0, 5, 72, 101, 108, 108, 111, 0, 0, 0, 4, 72, 101,
		108, 108, 0, 0, 0, 3, 72, 101, 108}
	var result arrayStruct
	expected := arrayStruct{
		FieldA: [3]byte{10, 17, 10},
		FieldB: [3]uint16{0, 1, 258},
		FieldC: [3]string{"Hello", "Hell", "Hel"},
	}
	Unmarshall(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestString(t *testing.T) {
	var result string
	b := []byte{0, 0, 0, 5, 72, 101, 108, 108, 111}

	Unmarshall(b, &result)
	if result != "Hello" {
		t.Errorf("Result not expected: Expected: %v Actual %v", "Hello", result)
	}
}

func TestAnonNestedStruct(t *testing.T) {
	type nestedStruct struct {
		FieldA uint16
		FieldB struct {
			FieldB1 uint16
			FieldB2 uint16
		}
		FieldC uint16
	}

	var result nestedStruct
	expected := nestedStruct{
		FieldA: 17,
		FieldB: struct {
			FieldB1 uint16
			FieldB2 uint16
		}{
			FieldB1: 13,
			FieldB2: 15},
		FieldC: 21,
	}
	b := []byte{0, 17, 0, 13, 0, 15, 0, 21}
	Unmarshall(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestArrayOfStruct(t *testing.T) {
	type innerStruct struct {
		FieldA1 int16
		FieldA2 bool
		FieldA3 int32
	}
	type arrayOfStruct struct {
		FieldA [3]innerStruct
	}

	expected := arrayOfStruct{FieldA: [3]innerStruct{
		innerStruct{FieldA1: 258, FieldA2: false, FieldA3: 50595078},
		innerStruct{FieldA1: 2828, FieldA2: true, FieldA3: 219025168},
		innerStruct{FieldA1: 5398, FieldA2: false, FieldA3: 387455258}}}
	b := []byte{1, 2, 0, 3, 4, 5, 6, 11, 12, 1, 13, 14, 15, 16, 21, 22, 0, 23, 24, 25, 26}
	var result arrayOfStruct
	Unmarshall(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}

}

func TestNestedStruct(t *testing.T) {
	type innerStruct struct {
		FieldA string
		FieldB uint16
	}
	type nestedStruct struct {
		FieldA uint16
		FieldB innerStruct
		FieldC uint16
	}

	var result nestedStruct
	expected := nestedStruct{FieldA: 17, FieldB: innerStruct{FieldA: "Hello", FieldB: 19}, FieldC: 50}
	b := []byte{0, 17, 0, 0, 0, 5, 72, 101, 108, 108, 111, 0, 19, 0, 50}
	Unmarshall(b, &result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

func TestBasicArray(t *testing.T) {
	//t.Skip()
	var result []uint16
	b := []byte{0, 0, 0, 1, 0, 2}
	expected := []uint16{uint16(0), uint16(1), uint16(2)}
	Unmarshall(b, &result)
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
		}
	}
}
