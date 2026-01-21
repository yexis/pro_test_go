package engine

type SingleJob struct {
	JobBase
	c Ctrl
}

func newSingle(c Ctrl) *SingleJob {
	return &SingleJob{c: c}
}

func (sj *SingleJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	if task != nil && task.Context != nil {
		if err := task.Context.Err(); err != nil {
			return nil, err
		}
	}
	i, e := sj.c(task, input, stage, sj.Params()...)
	if e != nil {
		return nil, e
	}
	return i, nil
}
