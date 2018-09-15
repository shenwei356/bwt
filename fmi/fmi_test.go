package fmi

import (
	"testing"
)

func TestLocate(t *testing.T) {
	var err error
	var loc []int
	var fmi *FMIndex

	type Case struct {
		s, q string
		m    int
		r    []int
	}

	var cases = []Case{
		Case{"mississippi", "iss", 0, []int{1, 4}},
		Case{"abcabcabc", "abc", 0, []int{0, 3, 6}},
		Case{"abcabcabc", "gef", 0, []int{}},
		Case{"abcabd", "abc", 1, []int{0, 3}},

		Case{"acctatac", "ac", 0, []int{0, 6}},
		Case{"acctatac", "tac", 0, []int{5}},
		Case{"acctatac", "tac", 1, []int{3, 5}},
		Case{"acctatac", "atac", 0, []int{4}},
		Case{"acctatac", "acctatac", 0, []int{0}},
		Case{"acctatac", "acctatac", 1, []int{0}},
		Case{"acctatac", "cctatac", 1, []int{1}},

		Case{"acctatac", "caa", 2, []int{1, 2, 3, 4, 5}},
		Case{"acctatac", "caa", 3, []int{0, 1, 2, 3, 4, 5}},
	}

	for i, c := range cases {
		fmi = NewFMIndex()
		_, err = fmi.TransformForLocate([]byte(c.s))
		if err != nil {
			t.Errorf("case #%d: TransformForLocate: %s", i+1, err)
			return
		}

		loc, err = fmi.Locate([]byte(c.q), c.m)
		if err != nil {
			t.Errorf("case #%d: Locate: %s", i, err)
			return
		}

		if len(loc) != len(c.r) {
			t.Errorf("case #%d: Locate '%s' in '%s' (allow %d mismatch), result: %d. right answer: %d", i+1, c.q, c.s, c.m, loc, c.r)
			return
		}

		for j := 0; j < len(loc); j++ {
			if loc[j] != c.r[j] {
				t.Errorf("case #%d: Locate '%s' in '%s' (allow %d mismatch), result: %d. right answer: %d", i+1, c.q, c.s, c.m, loc, c.r)
				return
			}
		}
	}
}
