package engine

type RecoverJob struct {
	JobBase
	c ActionI
	e ActionI
}

func newRecoverJob(c ActionI) *RecoverJob {
	return &RecoverJob{c: c}
}

func (rj *RecoverJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (i interface{}, e error) {
	i = input
	defer func() (interface{}, error) {
		if re := recover(); re != nil {
			//fmt.Printf("panic:%v stack:%v\n", re, string(debug.Stack()))
			var ok bool
			if e, ok = re.(error); ok && rj.e != nil {
				i, e = rj.e.Do(task, e, stage, rj.e.Params()...)
			}
		}
		return i, e
	}()
	i, e = rj.c.Do(task, input, stage, rj.c.Params()...)
	return i, e
}
