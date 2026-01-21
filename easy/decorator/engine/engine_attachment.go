package engine

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func (en *Engine) attachment(ac ActionI) {
	en.attachmentWalk(ac, &en.index)
}

func (en *Engine) attachmentWalk(ac ActionI, idx *int) {
	if ac == nil {
		return
	}

	// 先处理本节点的 name/index（仅当未设置时）
	switch v := ac.(type) {
	case *BlockJob:
		// BlockJob 自身不持有 JobBase，name/index 归属其内层 ActionI
		// 直接下钻
		en.attachmentWalk(v.ActionI, idx)
		return
	default:
		curIndex, curName := en.getIndexName(ac)
		if curIndex == 0 {
			*idx++
			ac.Indexed(*idx)
		}
		if curName == "" {
			if name := en.defaultName(ac, *idx); name != "" {
				ac.Named(name)
			} else {
				ac.Named(fmt.Sprintf("%s_%d", "unknown", *idx))
			}
		}
	}

	// 再递归处理子节点
	switch v := ac.(type) {
	case *SerialJob:
		for _, child := range v.actions {
			en.attachmentWalk(child, idx)
		}
	case *ConcurrentJob:
		for _, child := range v.actions {
			en.attachmentWalk(child, idx)
		}
	case *SwitchJob:
		en.attachmentWalk(v.judge, idx)
		for _, child := range v.cases {
			en.attachmentWalk(child, idx)
		}
	case *CatchJob:
		en.attachmentWalk(v.c, idx)
		en.attachmentWalk(v.e, idx)
	case *RecoverJob:
		en.attachmentWalk(v.c, idx)
		if v.e != nil {
			en.attachmentWalk(v.e, idx)
		}
	case *TimerJob:
		en.attachmentWalk(v.a, idx)
	}
}

func (en *Engine) getIndexName(ac ActionI) (int, string) {
	switch v := ac.(type) {
	case *SingleJob:
		return v.index, v.name
	case *SerialJob:
		return v.index, v.name
	case *ConcurrentJob:
		return v.index, v.name
	case *SwitchJob:
		return v.index, v.name
	case *CatchJob:
		return v.index, v.name
	case *RecoverJob:
		return v.index, v.name
	case *TimerJob:
		return v.index, v.name
	default:
		return 0, ""
	}
}

func (en *Engine) defaultName(ac ActionI, index int) string {
	switch v := ac.(type) {
	case *SingleJob:
		// 基础类
		return en.ctrlName(v.c)
	default:
		// 包装类：用类型名
		t := reflect.TypeOf(ac)
		if t == nil {
			return ""
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return fmt.Sprintf("%s__%d", t.Name(), index)
	}
}

func (en *Engine) ctrlName(c Ctrl) string {
	if c == nil {
		return ""
	}
	pc := reflect.ValueOf(c).Pointer()
	if pc == 0 {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	name := fn.Name()
	// 只保留最后一段，避免过长（pkg/path.Func）
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	return name
}

type Attachment struct {
	Index int
	Name  string
}

// VisitJobs ...
func (en *Engine) VisitJobs() []Attachment {
	root := en.g
	var out []Attachment

	var visit func(a ActionI)
	visit = func(a ActionI) {
		if a == nil {
			return
		}

		// BlockJob 自己没 JobBase，name/index 属于内部 ActionI
		if bj, ok := a.(*BlockJob); ok {
			visit(bj.ActionI)
			return
		}

		// 读取当前节点的 index/name（同包可读私有字段）
		switch v := a.(type) {
		case *SingleJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
		case *SerialJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			for _, c := range v.actions {
				visit(c)
			}
		case *ConcurrentJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			for _, c := range v.actions {
				visit(c)
			}
		case *SwitchJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			visit(v.judge)
			for _, c := range v.cases {
				visit(c)
			}
		case *CatchJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			visit(v.c)
			visit(v.e)
		case *RecoverJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			visit(v.c)
			if v.e != nil {
				visit(v.e)
			}
		case *TimerJob:
			out = append(out, Attachment{Index: v.index, Name: v.name})
			visit(v.a)
		default:
			// 未知实现：只收集不到 index/name，但尽量不崩
		}
	}

	visit(root)
	return out
}
