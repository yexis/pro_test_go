package engine

import (
	"context"
	"time"
)

type Engine struct {
	index int // Job 编号
	ctx   *Context
	g     ActionI
}

func NewEngine(ctx *Context) *Engine {
	return &Engine{ctx: ctx}
}

func (en *Engine) Build(action ...ActionI) *Engine {
	en.g = en.Serial(en.ToActions(action...))
	en.attachment(en.g)
	return en
}

func (en *Engine) Run(params ...interface{}) (interface{}, error) {
	en.ctx.initCallerContext(context.Background())
	var task = &Task{Context: en.ctx}
	var input interface{} = nil
	i, e := en.g.Do(task, input, nil, params...)
	if e == nil {
		e = en.ctx.Err()
	}
	en.ctx.WaitGlobalWaiter()
	return i, e
}

// Call ... wrap an action
func (en *Engine) Call(c Ctrl, opts ...JobOption) ActionI {
	return en.Apply(newSingle(c), opts...)
}

// Block ... set action's needWait
func (en *Engine) Block(a ActionI, nw bool, opts ...JobOption) ActionI {
	if nw {
		return en.Apply(en.BlockBT(a, BlockTypeBlock), opts...)
	} else {
		return en.Apply(en.BlockBT(a, BlockTypeFinalBlock), opts...)
	}
}

// BlockBT ...
func (en *Engine) BlockBT(a ActionI, bt BlockType, opts ...JobOption) ActionI {
	return en.Apply(newBlockAction(a, bt), opts...)
}

// Serial ... wrap a serial group
func (en *Engine) Serial(actions []ActionI, opts ...JobOption) ActionI {
	return en.Apply(newSerialJob(actions), opts...)
}

// Concurrent ... wrap a concurrent group
func (en *Engine) Concurrent(actions []ActionI, opts ...JobOption) ActionI {
	return en.Apply(newConcurrentJob(actions), opts...)
}

// Switch ... wrap a switch group
func (en *Engine) Switch(actions []ActionI, opts ...JobOption) ActionI {
	l := len(actions)
	if l <= 1 {
		return nil
	}
	return en.Apply(newSwitchJob(actions[0], actions[1:]), opts...)
}

// Catch ... deal with error
func (en *Engine) Catch(c ActionI, e ActionI, opts ...JobOption) ActionI {
	return en.Apply(newCatchJob(c, e), opts...)
}

// Recover ... deal with panic
func (en *Engine) Recover(c ActionI, opts ...JobOption) ActionI {
	return en.Apply(newRecoverJob(c), opts...)
}

// Timer ... wrap a timer group
func (en *Engine) Timer(c ActionI, t time.Duration, opts ...JobOption) ActionI {
	return en.Apply(newTimerJob(c, t), opts...)
}

// Apply ... attach sth to action
func (en *Engine) Apply(a ActionI, opts ...JobOption) ActionI {
	for _, opt := range opts {
		if opt != nil {
			opt(a)
		}
	}
	return a
}

func (en *Engine) ToActions(p ...ActionI) []ActionI {
	return append([]ActionI(nil), p...)
}
