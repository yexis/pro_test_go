package engine

import (
	"context"
	"sync"

	"github.com/yexis/pro_test_go/easy/easylogger"
)

type Context struct {
	easylogger.LoggerWriter
	logId string

	context context.Context
	cancel  context.CancelFunc

	globalWaiter sync.WaitGroup
	errOnce      sync.Once
	err          error
	mu           sync.RWMutex
}

type ContextOption func(*Context)

func WithLogger(l easylogger.LoggerWriter) ContextOption {
	return func(ctx *Context) {
		ctx.LoggerWriter = l
	}
}

// WithLogID 设置context的logid
func WithLogID(logid string) ContextOption {
	return func(c *Context) {
		c.logId = logid
	}
}

func NewContext(oriCtx context.Context) *Context {
	ctx, cancel := context.WithCancel(oriCtx)
	return &Context{
		context: ctx,
		cancel:  cancel,
	}
}

func (c *Context) Context() context.Context {
	return c.context
}

func (c *Context) initCallerContext(parent context.Context) {
	if parent == nil {
		parent = context.Background()
	}
	c.context, c.cancel = context.WithCancel(parent)
}

func (c *Context) SetErr(err error) {
	if err == nil {
		return
	}
	c.errOnce.Do(func() {
		c.mu.Lock()
		c.err = err
		c.mu.Unlock()
		if c.cancel != nil {
			c.cancel()
		}
	})
}

func (c *Context) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

func (c *Context) AddGlobalWaiter(delta int) {
	c.globalWaiter.Add(delta)
}

func (c *Context) GlobalWaiterDone() {
	c.globalWaiter.Done()
}

func (c *Context) WaitGlobalWaiter() {
	c.globalWaiter.Wait()
}
