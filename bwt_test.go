package bwt

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestTransformAndInverseTransform(t *testing.T) {
	s := "abracadabra"
	trans := "ard$rcaaaabb"
	tr, err := Transform([]byte(s), '$')
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

func TestSA(t *testing.T) {
	s := "mississippi"
	sa := SuffixArray([]byte(s))
	sa1 := []int{11, 10, 7, 4, 1, 0, 9, 8, 6, 3, 5, 2}
	// fmt.Printf("%s\nanswer: %v, result: %v", s, sa1, sa)
	if len(sa) != len(sa1) {
		t.Error(fmt.Errorf("sa error. answer: %v, result: %v", sa1, sa))
		return
	}
	for i := range sa {
		if sa[i] != sa1[i] {
			t.Error(fmt.Errorf("sa error. answer: %v, result: %v", sa1, sa))
			return
		}
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

var result []byte

func BenchmarkTransform(t *testing.B) {
	var r []byte
	var err error
	for i := 0; i < t.N; i++ {
		r, err = Transform(cases[0], '$')
		if err != nil {
			t.Error(err)
			return
		}
	}
	result = r
}
