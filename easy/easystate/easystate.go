package easystate

import (
	"strings"
)

type EasyState[K any] struct {
	idx    int
	buff   strings.Builder
	sep    string
	format func(K) string
}

func NewEasyState[K any](format func(K) string, params ...interface{}) *EasyState[K] {
	es := &EasyState[K]{
		sep:    "_",
		format: format,
	}
	if len(params) > 0 {
		if sz, ok := params[0].(int); ok {
			es.buff.Grow(sz)
		}
	}
	return es
}

func (es *EasyState[K]) SetSep(sep string) *EasyState[K] {
	es.sep = sep
	return es
}

func (es *EasyState[K]) AddState(data K) *EasyState[K] {
	if es.idx > 0 {
		es.buff.WriteString(es.sep)
	}
	es.idx++
	es.buff.WriteString(es.format(data))
	return es
}

func (es *EasyState[K]) GetState() string {
	return es.buff.String()
}
