package compare

import (
	"fmt"
)

type Diff interface {
	String() string
	Result() (interface{}, interface{})
}

type StringDiff struct {
	Ori    string
	Replay string
}

func (diff *StringDiff) String() string {
	return fmt.Sprintf("%v, %v", diff.Ori, diff.Replay)
}

func (diff *StringDiff) Result() (interface{}, interface{}) {
	return diff.Ori, diff.Replay
}

type IntDiff struct {
	Ori    int
	Replay int
}

func (diff *IntDiff) String() string {
	return fmt.Sprintf("%v, %v", diff.Ori, diff.Replay)
}

func (diff *IntDiff) Result() (interface{}, interface{}) {
	return diff.Ori, diff.Replay
}

type BytesDiff struct {
	Ori    []byte
	Replay []byte
}

func (diff *BytesDiff) String() string {
	return fmt.Sprintf("%v, %v", diff.Ori, diff.Replay)
}

func (diff *BytesDiff) Result() (interface{}, interface{}) {
	return diff.Ori, diff.Replay
}
