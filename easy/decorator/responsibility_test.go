package decorator

import (
	"errors"
	"fmt"
	"testing"
)

func TestResponsibilityContextNormal(t *testing.T) {
	step1 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, nil
		}},
		ErrTriggerStop: true,
	}

	step2 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, nil
		}},
		ErrTriggerStop: true,
	}

	step3 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, nil
		}},
		ErrTriggerStop: true,
	}

	en := NewResponsibilityEngine()
	en.Use(step1, step2, step3)
	input := []byte(`123`)
	i, e := en.Start(nil, input, nil)
	if e != nil {
		fmt.Printf("final err:%s\n", e.Error())
		t.Fatal(e)
	}
	fmt.Printf("final res:%s\n", i.([]byte))
}

func TestResponsibilityContextError(t *testing.T) {
	step1 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, errors.New("test error 1")
		}},
		ErrTriggerStop: false,
	}

	step2 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, errors.New("test error 2")
		}},
		ErrTriggerStop: true,
	}

	step3 := &ResponsibilityAction{
		Action: Action{C: func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
			return input, errors.New("test error 3")
		}},
		ErrTriggerStop: true,
	}

	en := NewResponsibilityEngine()
	en.Use(step1, step2, step3)
	input := []byte(`123`)
	i, e := en.Start(nil, input, nil)
	if e != nil {
		fmt.Printf("final err:%s\n", e.Error())
		t.Fatal(e)
		return
	}
	fmt.Printf("final res:%s\n", i.([]byte))
}
