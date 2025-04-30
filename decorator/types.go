package decorator

import (
	"context"
)

// ErrorMessage ... error message
const (
	EM1301EmptyAction   = "Empty actions"
	EM1302NotAction     = "Type not action"
	EM1303MissingParams = "Missing params"
	EM1304WaitTimeout   = "Waiting timeout after %s"
	EM1305WrongParams   = "Wrong params as %s"
)

// #region type

// Void ... void
type Void struct{}

// #endregion

// #region controller

// Task ... context for single request
type Task struct {
	Context context.Context
	Content interface{}
}

// Stage ... current stage of process
type Stage struct {
	A []*Action
	I int
	D interface{}
}

// Next ... call next function
func (s *Stage) Next(task *Task, input interface{}, stage *Stage) (interface{}, error) {
	l := len(stage.A)
	n := stage.I + 1
	if n >= l {
		return nil, nil
	}
	stage.I = n
	a := stage.A[stage.I]
	return a.C(task, input, stage, a.P...)
}

// Ctrl ... func pointer
type Ctrl func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error)

// Action ... single step or step groups
type Action struct {
	C Ctrl
	P []interface{}
	E Ctrl
}

// Selection ... selection
type Selection struct {
	Index int
	Data  interface{}
}

// #endregion
