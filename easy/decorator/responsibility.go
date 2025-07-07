package decorator

type ResponsibilityAction struct {
	Action
	ErrTriggerStop bool
}

type ResponsibilityContext struct {
	index    int
	handlers []*ResponsibilityAction
	Err      error
}

func (c *ResponsibilityContext) Next(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	c.index++
	if c.index >= len(c.handlers) {
		c.Err = nil
		return input, c.Err
	}

	a := c.handlers[c.index]
	i, e := a.Do(task, input, stage, params...)
	if e != nil && a.ErrTriggerStop {
		c.Err = e
		return c.Abort(task, input, stage, params...)
	}

	input = i
	c.Err = nil
	return c.Next(task, input, stage, params...)
}

func (c *ResponsibilityContext) Abort(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	c.index = len(c.handlers)
	return input, c.Err
}

type ResponsibilityEngine struct {
	handlers []*ResponsibilityAction
}

func NewResponsibilityEngine() *ResponsibilityEngine {
	return &ResponsibilityEngine{
		handlers: make([]*ResponsibilityAction, 0, 10),
	}
}

func (re *ResponsibilityEngine) Use(a ...*ResponsibilityAction) *ResponsibilityEngine {
	re.handlers = append(re.handlers, a...)
	return re
}

func (re *ResponsibilityEngine) Start(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	ctx := &ResponsibilityContext{
		index:    -1,
		handlers: re.handlers,
	}
	return ctx.Next(task, input, stage, params...)
}
