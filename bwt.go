package bwt

import (
	"errors"
	"index/suffixarray"
	"reflect"
	"sort"
)

// CheckEndSymbol is a global variable for checking end symbol before Burrows–Wheeler transform
var CheckEndSymbol = true

// ErrEndSymbolExisted means you should choose another EndSymbol
var ErrEndSymbolExisted = errors.New("bwt: end-symbol existed in string")

// ErrEmptySeq means a empty sequence is given
var ErrEmptySeq = errors.New("bwt: empty sequence")

// Transform returns Burrows–Wheeler transform of a byte slice.
// See https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform
func Transform(s []byte, es byte) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptySeq
	}
	if CheckEndSymbol {
		for _, c := range s {
			if c == es {
				return nil, ErrEndSymbolExisted
			}
		}
	}
	sa := SuffixArray(s)
	bwt, err := FromSuffixArray(s, sa, es)
	return bwt, err
}

// InverseTransform reverses the bwt to original byte slice. Not optimized yet.
func InverseTransform(t []byte, es byte) []byte {
	n := len(t)
	lines := make([][]byte, n)
	for i := 0; i < n; i++ {
		lines[i] = make([]byte, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			lines[j][n-1-i] = t[j]
		}
		sort.Sort(SliceOfByteSlice(lines))
	}

	s := make([]byte, n-1)
	for _, line := range lines {
		if line[n-1] == es {
			s = line[0 : n-1]
			break
		}
	}
	return s
}

// SuffixArray returns the suffix array of s.
// This function is the performance bottleneck of bwt and bwt/fmi package, with O(nlogn).
func SuffixArray(s []byte) []int {
	// sa := make([]int, len(s)+1)
	// sa[0] = len(s)

	// for i := 0; i < len(s); i++ {
	// 	sa[i+1] = i
	// }
	// sort.Slice(sa[1:], func(i, j int) bool {
	// 	return bytes.Compare(s[sa[i+1]:], s[sa[j+1]:]) < 0
	// })
	// return sa

	// https://github.com/shenwei356/bwt/issues/3 .
	// nearly copy from https://github.com/crazyleg/burrow-wheelers-golang/blob/master/pkg/bwtgolang/suffixarrayBWT.go#L8
	// It's 4X faster!
	//
	// benchmark                 old ns/op     new ns/op     delta
	// BenchmarkTransform-16     1339346       310706        -76.80%
	//
	// benchmark                 old allocs     new allocs     delta
	// BenchmarkTransform-16     4              5              +25.00%
	//
	// benchmark                 old bytes     new bytes     delta
	// BenchmarkTransform-16     55362         79962         +44.43%
	_sa := suffixarray.New(s)
	tmp := reflect.ValueOf(_sa).Elem().FieldByName("sa").FieldByIndex([]int{0})
	var sa []int = make([]int, len(s)+1)
	sa[0] = len(s)
	for i := 0; i < len(s); i++ {
		sa[i+1] = int(tmp.Index(i).Int())
	}
	return sa
}

// ErrInvalidSuffixArray means length of sa is not equal to 1+len(s)
var ErrInvalidSuffixArray = errors.New("bwt: invalid suffix array")

// FromSuffixArray compute BWT from sa
func FromSuffixArray(s []byte, sa []int, es byte) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptySeq
	}
	if len(s)+1 != len(sa) || sa[0] != len(s) {
		return nil, ErrInvalidSuffixArray
	}
	bwt := make([]byte, len(sa))
	bwt[0] = s[len(s)-1]
	for i := 1; i < len(sa); i++ {
		if sa[i] == 0 {
			bwt[i] = es
		} else {
			bwt[i] = s[sa[i]-1]
		}
	}
	return bwt, nil
}
