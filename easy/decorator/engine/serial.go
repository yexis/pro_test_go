package engine

type SerialJob struct {
	JobBase
	actions []ActionI
}

func newSerialJob(a []ActionI) *SerialJob {
	return &SerialJob{actions: a}
}

func (so *SerialJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	var e error
	var i interface{} = input
	if task == nil {
		return nil, New(EM1303MissingParams)
	}
	if task.Context != nil {
		if err := task.Context.Err(); err != nil {
			return nil, err
		}
	}
	for _, action := range so.actions {
		i, e = action.Do(task, i, stage, action.Params()...)
		// action 通过返回 *SerialStop 来主动中断后续 action（非 error）
		if stop, ok := i.(*SerialStop); ok {
			return stop.Result, nil
		}
		if e == nil {
			continue
		}
		return nil, e
	}
	return i, e
}
