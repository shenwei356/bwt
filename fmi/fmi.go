package fmi

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/shenwei356/bwt"
	"github.com/shenwei356/util/byteutil"
	"github.com/shenwei356/util/struct/stack"
)

// FMIndex is Burrows-Wheeler Index
type FMIndex struct {
	// EndSymbol
	EndSymbol byte

	// SuffixArray
	SuffixArray []int

	// Burrows-Wheeler Transform
	BWT []byte

	// First column of BWM
	F []byte

	// Alphabet in the BWT
	Alphabet []byte

	// Count of Letters in Alphabet
	CountOfLetters map[byte]int

	// C[c] is a table that, for each character c in the alphabet,
	// contains the number of occurrences of lexically smaller characters
	// in the text.
	C map[byte]int

	// Occ(c, k) is the number of occurrences of character c in the
	// prefix L[1..k], k is 0-based
	Occ map[byte]*[]int
}

// NewFMIndex is constructor of FMIndex
func NewFMIndex() *FMIndex {
	fmi := new(FMIndex)
	fmi.EndSymbol = '$'
	return fmi
}

// Transform return Burrows-Wheeler-Transform of s
func (fmi *FMIndex) Transform(s []byte) ([]byte, error) {
	var err error

	sa := bwt.SuffixArray(s)
	fmi.SuffixArray = sa

	fmi.BWT, err = bwt.FromSuffixArray(s, fmi.SuffixArray, fmi.EndSymbol)
	if err != nil {
		return nil, err
	}

	F := make([]byte, len(s)+1)
	F[0] = fmi.EndSymbol
	for i := 1; i <= len(s); i++ {
		F[i] = s[sa[i]]
	}
	fmi.F = F

	fmi.CountOfLetters = byteutil.CountOfByte(fmi.BWT)
	delete(fmi.CountOfLetters, fmi.EndSymbol)

	fmi.Alphabet = byteutil.AlphabetFromCountOfByte(fmi.CountOfLetters)

	fmi.C = ComputeC(fmi.F, fmi.Alphabet)

	fmi.Occ = ComputeOccurrence(fmi.BWT, fmi.Alphabet)

	return fmi.BWT, nil
}

// Last2First mapping
func (fmi *FMIndex) Last2First(i int) int {
	c := fmi.BWT[i]
	return fmi.C[c] + (*fmi.Occ[c])[i]
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
				start = fmi.C[c] + (*fmi.Occ[c])[match.start-1]
			}
			end = fmi.C[c] + (*fmi.Occ[c])[match.end] - 1
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
	buffer.WriteString("F:\n")
	buffer.WriteString(string(fmi.F))
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
func ComputeC(L []byte, alphabet []byte) map[byte]int {
	if alphabet == nil {
		alphabet = byteutil.Alphabet(L)
	}
	C := make(map[byte]int, len(alphabet))
	count := 0
	for _, c := range L {
		if _, ok := C[c]; !ok {
			C[c] = count
		}
		count++
	}
	return C
}

// ComputeOccurrence returns occurrence information.
// Occ(c, k) is the number of occurrences of character c in the prefix L[1..k]
func ComputeOccurrence(bwt []byte, letters []byte) map[byte]*[]int {
	if letters == nil {
		letters = byteutil.Alphabet(bwt)
	}
	occ := make(map[byte]*[]int, len(letters)-1)
	for _, letter := range letters {
		t := make([]int, 1, len(bwt))
		t[0] = 0
		occ[letter] = &t
	}
	t := make([]int, 1, len(bwt))
	t[0] = 1
	occ[bwt[0]] = &t
	var letter, k byte
	var v *[]int
	for _, letter = range bwt[1:] {
		for k, v = range occ {
			if k == letter {
				*v = append(*v, (*v)[len(*v)-1]+1)
			} else {
				*v = append(*v, (*v)[len(*v)-1])
			}
		}
	}
	return occ
}
