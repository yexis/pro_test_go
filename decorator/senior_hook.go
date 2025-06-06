package decorator

// SeniorHook ... u should define K by yourself
type SeniorHook[K comparable] struct {
	mp map[K]*Action
}

func NewSeniorHook[K comparable]() *SeniorHook[K] {
	return &SeniorHook[K]{
		mp: make(map[K]*Action),
	}
}

func (h *SeniorHook[K]) AddHook(k K, action *Action) {
	h.mp[k] = action
}

func (h *SeniorHook[K]) RemoveHook(k K) {
	delete(h.mp, k)
}

func (h *SeniorHook[K]) DoHook(task *Task, input interface{}, stage *Stage, k K) bool {
	if action, ok := h.mp[k]; ok {
		_, e := action.C(task, input, stage, action.P...)
		if e != nil {
			_, _ = action.E(task, input, stage, action.P...)
		}
		return true
	}
	return false
}
