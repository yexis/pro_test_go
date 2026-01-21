package engine

type SwitchJob struct {
	JobBase
	judge ActionI
	cases []ActionI
}

func newSwitchJob(j ActionI, c []ActionI) *SwitchJob {
	return &SwitchJob{
		judge: j,
		cases: c,
	}
}

func (sj *SwitchJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	if task == nil {
		return nil, New(EM1303MissingParams)
	}
	if task.Context != nil {
		if err := task.Context.Err(); err != nil {
			return nil, err
		}
	}

	l := len(sj.cases)
	if l <= 0 {
		return nil, New(EM1303MissingParams)
	}
	ji, je := sj.judge.Do(task, input, stage, sj.judge.Params()...)
	if je != nil {
		return nil, je
	}

	s, ok := ji.(*Selection)
	if !ok {
		return nil, FromCode(SwitchValNotSelection)
	}
	if s.Index < 0 || s.Index >= l {
		return nil, FromCode(SwitchBeyond)
	}

	selected := sj.cases[s.Index]
	return selected.Do(task, input, stage, selected.Params()...)
}
