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
