package fmi

import "testing"

func TestCount(t *testing.T) {
	s := "abracadabra"
	fmi := NewFMIndex()

	_, err := fmi.Transform([]byte(s))
	if err != nil {
		t.Errorf("Test failed: TestBWIndex: %s", err)
	}
	if fmi.Count([]byte("bra")) != 2 || fmi.Count([]byte("a")) != 5 || fmi.Count([]byte("b")) != 2 {
		t.Errorf("Test failed: TestBWIndex: %s", "Count failed")
	}
}

func TestLocate(t *testing.T) {
	var err error
	fmi := NewFMIndex()
	var loc []int

	type Case struct {
		s, q string
		m    int
		r    []int
	}

	var cases = []Case{
		Case{"GATGCGAGAGATG", "GAGA", 0, []int{5, 7}},
		Case{"abracadabra", "ab", 0, []int{0, 7}},

		Case{"abcabd", "abc", 1, []int{0, 3}},
		Case{"abcabd", "abd", 1, []int{0, 3}},
		Case{"abcabd", "bc", 0, []int{1}},
		Case{"abcabd", "bc", 1, []int{1, 4}},
	}

	for i, c := range cases {
		_, err = fmi.TransformForLocate([]byte(c.s))
		if err != nil {
			t.Errorf("case #%d: TransformForLocate: %s", i+1, err)
		}

		loc, err = fmi.Locate([]byte(c.q), c.m)
		if err != nil {
			t.Errorf("case #%d: Locate: %s", i, err)
		}
		if len(loc) != len(c.r) {
			t.Errorf("case #%d: Locate %s at %s, result: %d", i+1, c.q, c.s, loc)
			break
		}

		for j := 0; j < len(loc); j++ {
			if loc[j] != c.r[j] {
				t.Errorf("case #%d: Locate %s at %s, result: %d", i+1, c.q, c.s, loc)
			}
		}
	}

}
