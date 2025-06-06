package orderlist

type HookKey int

const (
	errorKey HookKey = iota
	costKey
)

type Recorder interface {
	GetRecordKey() HookKey
}

type record struct {
	hk    HookKey
	index int
}

func (r *record) GetRecordKey() HookKey {
	return r.hk
}

type errorRecord struct {
	record
	err error
}

type costRecord struct {
	record
	ms int
}

type listEventType int

const (
	DataType listEventType = iota
	ErrorType
	EndType
)

type NodeStatistic struct {
	ms  int   // cost time
	err error // err
}
