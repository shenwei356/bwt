package bwt

import "bytes"

// SliceOfByteSlice is [][]byte
type SliceOfByteSlice [][]byte

func (s SliceOfByteSlice) Len() int { return len(s) }
func (s SliceOfByteSlice) Less(i, j int) bool {
	return bytes.Compare(s[i], s[j]) < 0
}
func (s SliceOfByteSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
