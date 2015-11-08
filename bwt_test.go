package bwt

import "testing"

func TestTransformAndInverseTransform(t *testing.T) {
	s := "abracadabra"
	trans := "ard$rcaaaabb"
	tr, _, err := Transform([]byte(s), '$')
	if err != nil {
		t.Error(err)
	}
	if string(tr) != trans {
		t.Error("Test failed: Transform")
	}
	if string(InverseTransform([]byte(trans), '$')) != s {
		t.Error("Test failed: InverseTransform")
	}
}
