package bwt

import (
	"errors"
	"sort"

	"github.com/shenwei356/util/byteutil"
)

// ErrEndSymbolExisted means you should choose another EndSymbol
var ErrEndSymbolExisted = errors.New("bwt: end-symbol existed in string")

// Transform returns Burrows–Wheeler transform  of a byte slice.
// See https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform
func Transform(s []byte, es byte) ([]byte, [][]byte, error) {
	count := byteutil.CountOfByte(s)
	if _, ok := count[es]; ok {
		return nil, nil, ErrEndSymbolExisted
	}
	s = append(s, es)
	n := len(s)

	rotations := make([][]byte, n)
	i := 0
	for j := 0; j < n; j++ {
		rotations[i] = append(s[n-j:], s[0:n-j]...)
		i++
	}
	sort.Sort(byteutil.SliceOfByteSlice(rotations))

	bwt := make([]byte, n)
	i = 0
	for _, t := range rotations {
		bwt[i] = t[n-1]
		i++
	}
	return bwt, rotations, nil
}

// InverseTransform reverses the bwt to original byte slice
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
		sort.Sort(byteutil.SliceOfByteSlice(lines))
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

// SuffixArray returns the suffix array of s
func SuffixArray(s []byte) []int {
	n := len(s)
	suffixMap := make(map[string]int, n)
	for i := 0; i < n; i++ {
		suffixMap[string(s[i:])] = i
	}
	suffixes := make([]string, n)
	i := 0
	for suffix := range suffixMap {
		suffixes[i] = suffix
		i++
	}
	indice := make([]int, n+1)
	indice[0] = n
	i = 1
	sort.Strings(suffixes)
	for _, suffix := range suffixes {
		indice[i] = suffixMap[suffix]
		i++
	}
	return indice
}

// ErrInvalidSuffixArray means length of sa is not equal to 1+len(s)
var ErrInvalidSuffixArray = errors.New("bwt: invalid suffix array")

// FromSuffixArray compute BWT from sa
func FromSuffixArray(s []byte, sa []int, es byte) ([]byte, error) {
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
