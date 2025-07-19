package easypool

import "sync"

var (
	// useObjectPool 使用对象池
	useObjectPool = true
)

// DisableObjectPool ...
// it's for debug
func DisableObjectPool() {
	useObjectPool = false
}

// ObjectPool 对象池
type ObjectPool interface {
	Get() any
	Put(any)
}

type objectPool struct {
	raw     *sync.Pool
	creator func() any
}

// NewObjectPool ...
// to new an object pool
func NewObjectPool(creator func() any) ObjectPool {
	if creator == nil {
		panic("ObjectPool creator is nil")
	}
	return &objectPool{
		raw: &sync.Pool{
			New: creator,
		},
		creator: creator,
	}
}

func (op *objectPool) Get() any {
	var obj any
	if useObjectPool {
		obj = op.raw.Get()
	}

	if obj == nil {
		obj = op.creator()
	}

	return obj
}

func (op *objectPool) Put(obj any) {
	if useObjectPool {
		op.raw.Put(obj)
	}
}
