package fmi

import (
	"testing"

	"github.com/shenwei356/bwt"
)

type Case struct {
	s, q string
	m    int
	r    []int
}

var cases = []Case{
	{"", "abc", 0, []int{}},
	{"mississippi", "", 0, []int{}},
	{"mississippi", "iss", 0, []int{1, 4}},
	{"abcabcabc", "abc", 0, []int{0, 3, 6}},
	{"abcabcabc", "gef", 0, []int{}},
	{"abcabcabc", "gef", 0, []int{}},
	{"abcabcabc", "xef", 0, []int{}},
	{"abcabcabc", "xabcb", 1, []int{}},
	{"abcabcabc", "xabcb", 2, []int{2}},
	{"abcabd", "abc", 1, []int{0, 3}},

	{"acctatac", "ac", 0, []int{0, 6}},
	{"acctatac", "tac", 0, []int{5}},
	{"acctatac", "tac", 1, []int{3, 5}},
	{"acctatac", "taz", 1, []int{3, 5}},
	{"ccctatac", "tzc", 1, []int{5}},
	{"acctatac", "atac", 0, []int{4}},
	{"acctatac", "acctatac", 0, []int{0}},
	{"acctatac", "acctatac", 1, []int{0}},
	{"acctatac", "cctatac", 1, []int{1}},

	{"acctatac", "caa", 2, []int{1, 2, 3, 4, 5}},
	{"acctatac", "caa", 3, []int{0, 1, 2, 3, 4, 5}},
}

func TestLocate(t *testing.T) {
	var err error
	var match bool
	var fmi *FMIndex

	for i, c := range cases {
		fmi = NewFMIndex()
		_, err = fmi.Transform([]byte(c.s))
		if err != nil {
			if c.s == "" && err == bwt.ErrEmptySeq {
				continue
			} else {
				t.Errorf("case #%d: Transform: %s", i+1, err)
				return
			}
		}

		match, err = fmi.Match([]byte(c.q), c.m)
		if err != nil {
			t.Errorf("case #%d: Locate: %s", i, err)
			return
		}

		if match != (len(c.r) > 0) {
			t.Errorf("case #%d: Match '%s' in '%s' (allow %d mismatch), result: %v. right answer: %v", i+1, c.q, c.s, c.m, match, len(c.r) > 0)
			return
		}

	}
}

func TestMatch(t *testing.T) {
	var err error
	var loc []int
	var fmi *FMIndex

	for i, c := range cases {
		fmi = NewFMIndex()
		_, err = fmi.Transform([]byte(c.s))
		if err != nil {
			if c.s == "" && err == bwt.ErrEmptySeq {
				continue
			} else {
				t.Errorf("case #%d: Transform: %s", i+1, err)
				return
			}
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
