package splitter

import (
	"fmt"
)

type BytesSplitter = Splitter[byte]
type Splitter[T byte | int8 | int16] struct {
	data     []T
	offset   int
	overlap  int
	needDrop bool
}

func NewSplitter[T byte | int8 | int16]() *Splitter[T] {
	return &Splitter[T]{
		offset:  0,
		overlap: 0,
		data:    make([]T, 0),
	}
}

func (sp *Splitter[T]) SetDrop(b bool) *Splitter[T] {
	sp.needDrop = b
	return sp
}

func (sp *Splitter[T]) SetOverlap(d int) *Splitter[T] {
	if d >= 0 {
		sp.overlap = d
	}
	return sp
}

func (sp *Splitter[T]) Read(l int) (bool, []T, error) {
	if l <= 0 {
		return false, nil, fmt.Errorf("invalid read length")
	}
	start := sp.offset - sp.overlap
	if start < 0 || start+l > len(sp.data) {
		return false, nil, fmt.Errorf("read length too big")
	}
	out := sp.data[sp.offset : sp.offset+l]
	sp.offset = start + l

	if sp.needDrop && sp.offset > 0 {
		// dropped the read
		sp.data = sp.data[sp.offset:]
		sp.offset = 0
	}

	return true, out, nil
}

func (sp *Splitter[T]) Bytes() (bool, []T, error) {
	return true, sp.data, nil
}

func (sp *Splitter[T]) Write(data []T) error {
	sp.data = append(sp.data, data...)
	return nil
}
