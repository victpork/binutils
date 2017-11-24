package binutils

import (
	"testing"
)

func TestBinUtil(t *testing.T) {
	type sshptyRequest struct {
		Term    string
		Width   uint32
		Height  uint32
		PWidth  uint32
		PHeight uint32
	}
	var req sshptyRequest
	expect := sshptyRequest{Term: "xterm", Width: 80, Height: 24}
	b := []byte{0, 0, 0, 5, 120, 116, 101, 114, 109, 0, 0, 0, 80, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 21, 3, 0, 0, 0, 127, 42, 0, 0, 0, 1, 128, 0, 0, 150, 0, 129, 0, 0, 150, 0, 0}
	binstruct := CreateByteStream(b)
	binstruct.ReadAsStruct(&req)
	if req != expect {
		t.Errorf("Term does not match, expected '%v' now '%v'", expect, req)
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
	expect := nestedStruct{
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
	binstruct := CreateByteStream(b)
	binstruct.ReadAsStruct(&result)
	if result != expect {
		t.Errorf("Result not expected: Expected: %v Actual %v", expect, result)
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
	binstruct := CreateByteStream(b)
	binstruct.ReadAsStruct(&result)
	if result != expected {
		t.Errorf("Result not expected: Expected: %v Actual %v", expected, result)
	}
}

// This test case will fail misarably
func TestArray(t *testing.T) {
	type array [3]uint16
	var result array
	b := []byte{0, 0, 0, 1, 0, 2}
	binstruct := CreateByteStream(b)
	binstruct.ReadAsStruct(&result)

}
