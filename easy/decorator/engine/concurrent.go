package engine

import (
	"sync"
)

type ConcurrentJob struct {
	JobBase
	actions    []ActionI
	currWaiter sync.WaitGroup
}

func newConcurrentJob(a []ActionI) *ConcurrentJob {
	cj := &ConcurrentJob{
		actions: a,
	}
	return cj
}

func (cj *ConcurrentJob) Do(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	if task == nil {
		return nil, New(EM1303MissingParams)
	}
	baseCtx := task.Context
	if baseCtx != nil {
		if err := baseCtx.Err(); err != nil {
			return nil, err
		}
	}
	l := len(cj.actions)
	res := make([]interface{}, l)
	for idx, action := range cj.actions {
		var needWait BlockType
		v, ok := action.(BlockActionI)
		if ok {
			needWait = v.NeedWait()
		}
		if needWait == BlockTypeBlock {
			cj.currWaiter.Add(1)
		} else if needWait == BlockTypeFinalBlock {
			if baseCtx != nil {
				baseCtx.globalWaiter.Add(1)
			}
		} else if needWait == BlockTypeNoBlock {
			// do nothing
		}
		localIdx := idx
		localAction := action
		// 并发分支用独立 Task，避免 data race；但共享同一个 Context（全局 Context）
		childTask := &Task{
			Content: task.Content,
			Context: baseCtx,
		}
		go func(needWait BlockType, tsk *Task) {
			defer func() {
				if needWait == BlockTypeBlock {
					cj.currWaiter.Done()
				} else if needWait == BlockTypeFinalBlock {
					if baseCtx != nil {
						baseCtx.GlobalWaiterDone()
					}
				}
			}()
			i, e := localAction.Do(tsk, input, stage, localAction.Params()...)
			if e != nil {
				if needWait == BlockTypeBlock {
					res[localIdx] = e
				} else if needWait == BlockTypeFinalBlock {
					task.Context.SetErr(e)
				}
			} else {
				if needWait == BlockTypeBlock {
					res[localIdx] = i
				} else if needWait == BlockTypeFinalBlock {
					// do nothing
				}
			}

		}(needWait, childTask)
	}

	cj.currWaiter.Wait()
	if baseCtx != nil {
		if err := baseCtx.Err(); err != nil {
			return nil, err
		}
	}
	return res, nil
}
