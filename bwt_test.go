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

func TestFromSuffixArray(t *testing.T) {
	s := "GATGCGAGAGATG"
	trans := "GGGGGGTCAA$TAA"

	sa := SuffixArray([]byte(s))
	B, err := FromSuffixArray([]byte(s), sa, '$')
	if err != nil {
		t.Error("Test failed: FromSuffixArray error")
	}
	if string(B) != trans {
		t.Error("Test failed: FromSuffixArray returns wrong result")
	}
}
