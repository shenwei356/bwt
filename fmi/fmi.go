package fmi

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/shenwei356/bwt"
	"github.com/shenwei356/util/byteutil"
	"github.com/shenwei356/util/struct/stack"
)

// FMIndex is Burrows-Wheeler Index
type FMIndex struct {
	EndSymbol byte
	// Burrows-Wheeler Transform
	BWT []byte
	// Alphabet in the BWT
	Alphabet []byte
	// Count of Letters in Alphabet
	CountOfLetters map[byte]int
	// Matrix M where each row is a rotation of the text,
	// and the rows have been sorted lexicographically.
	M [][]byte
	// C[c] is a table that, for each character c in the alphabet,
	// contains the number of occurrences of lexically smaller characters
	// in the text.
	C map[byte]int
	// Occ(c, k) is the number of occurrences of character c in the
	// prefix L[1..k], k is 0-based
	Occ map[byte][]int
	// SuffixArray, only used when call Locate
	SuffixArray []int
}

// NewFMIndex is constructor of FMIndex
func NewFMIndex() *FMIndex {
	fmi := new(FMIndex)
	fmi.EndSymbol = '$'
	return fmi
}

// Transform return Burrows-Wheeler-Transform of s
func (fmi *FMIndex) Transform(s []byte) ([]byte, error) {
	bwt, rotations, err := bwt.Transform(s, fmi.EndSymbol)
	if err != nil {
		return nil, err
	}
	fmi.BWT = bwt
	fmi.M = rotations
	fmi.CountOfLetters = byteutil.CountOfByte(fmi.BWT)
	delete(fmi.CountOfLetters, fmi.EndSymbol)
	fmi.Alphabet = byteutil.AlphabetFromCountOfByte(fmi.CountOfLetters)
	fmi.C = ComputeC(fmi.M, fmi.Alphabet)
	fmi.Occ = ComputeOccurrence(fmi.BWT, fmi.Alphabet)
	return bwt, nil
}

// TransformForLocate compute SuffixArray in addition to Transform
func (fmi *FMIndex) TransformForLocate(s []byte) ([]byte, error) {
	b, err := fmi.Transform(s)
	if err != nil {
		return nil, err
	}
	fmi.SuffixArray = bwt.SuffixArray(s)
	return b, nil
}

// Last2First mapping
func (fmi *FMIndex) Last2First(i int) int {
	c := fmi.BWT[i]
	return fmi.C[c] + fmi.Occ[c][i]
}

func (fmi *FMIndex) nextLetterInAlphabet(c byte) byte {
	var nextLetter byte
	for i, letter := range fmi.Alphabet {
		if letter == c {
			if i < len(fmi.Alphabet)-1 {
				nextLetter = fmi.Alphabet[i+1]
			} else {
				nextLetter = fmi.Alphabet[i]
			}
			break
		}
	}
	return nextLetter
}

// ErrSuffixArrayIsNil is xxx
var ErrSuffixArrayIsNil = errors.New("bwt/fmi: SuffixArray is nil, you should call TransformForLocate instead of Transform")

type sMatch struct {
	query      []byte
	start, end int
	mismatches int
}

// Locate locates the pattern
func (fmi *FMIndex) Locate(query []byte, mismatches int) ([]int, error) {
	locations := []int{}
	locationsMap := make(map[int]struct{})
	letters := byteutil.Alphabet(query)

	for _, letter := range letters { // query having illegal letter
		if _, ok := fmi.CountOfLetters[letter]; !ok {
			return locations, nil
		}
	}

	if fmi.SuffixArray == nil {
		return nil, ErrSuffixArrayIsNil
	}

	n := len(fmi.BWT)
	var matches stack.Stack

	// start and end are 0-based
	matches.Put(sMatch{query: query, start: 0, end: n - 1, mismatches: mismatches})
	// fmt.Printf("====%s====\n", query)
	// fmt.Println(fmi)
	var match sMatch
	var last byte
	var start, end int
	var m int
	for !matches.Empty() {
		match = matches.Pop().(sMatch)
		query = match.query[0 : len(match.query)-1]
		last = match.query[len(match.query)-1]
		if match.mismatches == 0 {
			letters = []byte{last}
		} else {
			letters = fmi.Alphabet
		}

		// fmt.Printf("\n%s, %s, %c\n", match.query, query, last)
		// fmt.Printf("query: %s, last: %c\n", query, last)
		for _, c := range letters {
			// fmt.Printf("letter: %c, start: %d, end: %d, mismatches: %d\n", c, match.start, match.end, match.mismatches)
			if match.start == 0 {
				start = fmi.C[c] + 0
			} else {
				start = fmi.C[c] + fmi.Occ[c][match.start-1]
			}
			end = fmi.C[c] + fmi.Occ[c][match.end] - 1
			//fmt.Printf("    s: %d, e: %d\n", start, end)

			if start <= end {
				if len(query) == 0 {
					for _, i := range fmi.SuffixArray[start : end+1] {
						// fmt.Printf("    >>> found: %d\n", i)
						locationsMap[i] = struct{}{}
					}
				} else {
					m = match.mismatches
					if c != last {
						if match.mismatches > 1 {
							m = match.mismatches - 1
						} else {
							m = 0
						}
					}

					// fmt.Printf("    >>> candidate: query: %s, start: %d, end: %d, m: %d\n", query, start, end, m)
					matches.Put(sMatch{query: query, start: start, end: end, mismatches: m})
				}
			}
		}
	}
	i := 0
	locations = make([]int, len(locationsMap))
	for loc := range locationsMap {
		locations[i] = loc
		i++
	}
	sort.Ints(locations)
	return locations, nil
}

func (fmi *FMIndex) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("EndSymbol: %c\n", fmi.EndSymbol))
	buffer.WriteString(fmt.Sprintf("BWT: %s\n", string(fmi.BWT)))
	buffer.WriteString(fmt.Sprintf("Alphabet: %s\n", string(fmi.Alphabet)))
	buffer.WriteString("M:\n")
	for i, r := range fmi.M {
		buffer.WriteString(fmt.Sprintf("%5d: %s\n", i, string(r)))
	}
	buffer.WriteString("C:\n")
	for _, letter := range fmi.Alphabet {
		buffer.WriteString(fmt.Sprintf("  %c: %d\n", letter, fmi.C[letter]))
	}
	buffer.WriteString("Occ:\n")
	buffer.WriteString(fmt.Sprintf("  BWT[%s]\n", strings.Join(strings.Split(string(fmi.BWT), ""), " ")))
	for _, letter := range fmi.Alphabet {
		buffer.WriteString(fmt.Sprintf("  %c: %v\n", letter, fmi.Occ[letter]))
	}

	buffer.WriteString("SA:\n")
	buffer.WriteString(fmt.Sprintf("  %d\n", fmi.SuffixArray))

	return buffer.String()
}

// ComputeC computes C.
// C[c] is a table that, for each character c in the alphabet,
// contains the number of occurrences of lexically smaller characters
// in the text.
func ComputeC(m [][]byte, alphabet []byte) map[byte]int {
	if alphabet == nil {
		byteutil.Alphabet(m[0])
	}
	C := make(map[byte]int, len(alphabet))
	count := 0
	for _, r := range m {
		c := r[0]
		if _, ok := C[c]; !ok {
			C[c] = count
		}
		count++
	}
	return C
}

// ComputeOccurrence returns occurrence information.
// Occ(c, k) is the number of occurrences of character c in the prefix L[1..k]
func ComputeOccurrence(bwt []byte, letters []byte) map[byte][]int {
	if letters == nil {
		letters = byteutil.Alphabet(bwt)
	}
	occ := make(map[byte][]int, len(letters)-1)
	for _, letter := range letters {
		occ[letter] = []int{0}
	}
	occ[bwt[0]] = []int{1}
	for _, letter := range bwt[1:] {
		for k := range occ {
			if k == letter {
				occ[k] = append(occ[k], occ[k][len(occ[k])-1]+1)
			} else {
				occ[k] = append(occ[k], occ[k][len(occ[k])-1])
			}
		}
	}
	return occ
}
