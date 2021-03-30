package fmi

type sMatch struct {
	query      []byte
	start, end int
	mismatches int
}

// Stack struct
type Stack []sMatch

// Empty tell if it is empty
func (s Stack) Empty() bool {
	return len(s) == 0
}

// Peek return the last element
func (s Stack) Peek() sMatch {
	return s[len(s)-1]
}

// Put puts element to stack
func (s *Stack) Put(i sMatch) {
	(*s) = append((*s), i)
}

// Pop pops element from the stack
func (s *Stack) Pop() sMatch {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}
