package engine

type CatchJob struct {
	JobBase
	c ActionI
	e ActionI
}

func newCatchJob(c ActionI, e ActionI) *CatchJob {
	return &CatchJob{c: c, e: e}
}

func (cj *CatchJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	var e error
	var i interface{} = input
	i, e = cj.c.Do(task, input, stage, cj.c.Params()...)
	if e != nil {
		i, e = cj.e.Do(task, e, stage, cj.e.Params()...)
		if e != nil {
			return nil, e
		}
		return i, nil
	}
	return i, e
}
