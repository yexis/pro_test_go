package engine

type ActionI interface {
	Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error)
	JobI
}

type BlockActionI interface {
	ActionI
	NeedWait() BlockType
}

type JobI interface {
	Named(name string) JobI
	Paramed(params ...interface{}) JobI
	Params() []interface{}
	Indexed(index int) JobI
}

// CostRecorder ...
// in order to state time cost
type CostRecorder interface {
	Start(key string) string
	Finish(key string)
}
