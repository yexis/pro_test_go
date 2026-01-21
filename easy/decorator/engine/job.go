package engine

type JobBase struct {
	index int
	name  string
	P     []interface{}
}

func (j *JobBase) Indexed(i int) JobI {
	j.index = i
	return j
}

func (j *JobBase) Named(n string) JobI {
	j.name = n
	return j
}

func (j *JobBase) Paramed(params ...interface{}) JobI {
	j.P = append(j.P, params...)
	return j
}

func (j *JobBase) Params() []interface{} {
	return j.P
}

// JobOption ... 用于给 Job 添加属性
type JobOption func(ActionI)

func WithName(name string) JobOption {
	return func(a ActionI) {
		a.Named(name)
	}
}

func WithIndex(index int) JobOption {
	return func(a ActionI) {
		a.Indexed(index)
	}
}

func WithParams(p ...interface{}) JobOption {
	return func(a ActionI) {
		a.Paramed(p...)
	}
}
