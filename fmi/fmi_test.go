package fmi

import "testing"

func TestBWIndex(t *testing.T) {
	s := "abracadabra"
	fmi := NewFMIndex()

	_, err := fmi.Transform([]byte(s))
	if err != nil {
		t.Errorf("Test failed: TestBWIndex: %s", err)
	}
	// fmt.Println(fmi)
	if fmi.Count([]byte("bra")) != 2 || fmi.Count([]byte("a")) != 5 || fmi.Count([]byte("b")) != 2 {
		t.Errorf("Test failed: TestBWIndex: %s", "Count failed")
	}

	_, err = fmi.TransformForLocate([]byte(s))
	if err != nil {
		t.Errorf("Test failed: TestBWIndex: %s", err)
	}
	fmi.Locate([]byte("ab"), 0)
}
