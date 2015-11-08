package bwt

import (
	"errors"
	"sort"

	"github.com/shenwei356/util/byteutil"
)

// ErrEndSymbolExisted means you should choose another EndSymbol
var ErrEndSymbolExisted = errors.New("EndSymbol existed in string, please Choose another one")

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
