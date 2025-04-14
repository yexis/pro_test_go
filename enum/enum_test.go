package enum

import (
	"fmt"
	"testing"
)

// TriggerBrainAgentState ... 访问brain-agent状态
type TriggerBrainAgentState int

func ToTriggerBrainAgentStatePtr(ts TriggerBrainAgentState) *TriggerBrainAgentState {
	return &ts
}

const (
	// NoTriggerByDefault ... 默认值不访问
	NoTriggerByDefault TriggerBrainAgentState = iota
	// NoTriggerByEmptySatisfaction ... 空满意度
	NoTriggerByEmptySatisfaction
	// NoTriggerByUSDecide ... us满意度决定
	NoTriggerByUSDecide
	// NoTriggerByDuplexReject ... 全双工拒识
	NoTriggerByDuplexReject

	// TriggerByUsDecide ... us满意度决定
	TriggerByUsDecide = 100
	// TriggerByLLMResponse ... 大模型结果，bot不支持
	TriggerByLLMResponse
)

func TestEnum(t *testing.T) {
	s := TriggerBrainAgentState(TriggerByUsDecide)
	fmt.Println(s)

	s = TriggerBrainAgentState(TriggerByLLMResponse)
	fmt.Println(s)
}
