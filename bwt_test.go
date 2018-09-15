package bwt

import (
	"math/rand"
	"testing"
)

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

var cases [][]byte

func init() {
	rand.Seed(1)
	alphabet := "ACGT"
	n := len(alphabet)
	scales := []float32{1e3, 1e5}
	cases = make([][]byte, len(scales))
	for i, scale := range scales {
		l := rand.Float32() * scale * 10
		buf := make([]byte, int(l))
		for j := 0; j < int(l); j++ {
			buf[j] = alphabet[rand.Intn(n)]
		}
		cases[i] = buf
	}
}

func BenchmarkTransform(t *testing.B) {
	for _, s := range cases {
		_, _, err := Transform(s, '$')
		if err != nil {
			t.Error(err)
			return
		}
	}
}
