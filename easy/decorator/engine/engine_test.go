package engine

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func switchJobCaseC(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
	fmt.Println("I am switchJobCaseC")
	return nil, nil
}
func TestEngineByPeople(t *testing.T) {
	serialFuncA := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialFuncA")
		return nil, nil
	}
	serialFuncB := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialFuncB")
		return nil, nil
	}
	serialGroupFuncA := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialGroupFuncA")
		return nil, nil
	}
	serialGroupFuncAException := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialGroupFuncAException")
		return nil, nil
	}
	serialGroupFuncBSkip := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialGroupFuncBSkip")
		return StopSerial(nil), nil
	}
	serialGroupFuncCSkipped := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am serialGroupFuncCSkipped")
		return nil, nil
	}
	concurrentJobA := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am concurrentJobA")
		return nil, nil
	}
	concurrentJobB := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		time.Sleep(3 * time.Second)
		fmt.Println("I am concurrentJobB")
		return nil, nil
	}
	switchJobJudge := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am switchJobJudge")
		panic(errors.New("mock panic"))
		return nil, nil
	}
	switchJobCaseA := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am switchJobCaseA")
		return nil, nil
	}
	switchJobCaseB := func(task *Task, input interface{}, stage *Stage, params ...interface{}) (interface{}, error) {
		fmt.Println("I am switchJobCaseB")
		return nil, nil
	}

	enCtx := NewContext(context.Background())
	en := NewEngine(enCtx)
	process := en.Build(
		en.Recover(
			en.Serial(
				en.ToActions(
					en.Call(serialFuncA),
					en.Call(serialFuncB),
					en.Serial(
						en.ToActions(
							en.Catch(
								en.Call(serialGroupFuncA, WithName("serialGroupFuncA")),
								en.Call(serialGroupFuncAException),
							),
							en.Call(serialGroupFuncBSkip),
							en.Call(serialGroupFuncCSkipped)),
					),
					en.Concurrent(
						en.ToActions(
							en.Block(en.Call(concurrentJobA), true),
							en.Block(en.Call(concurrentJobB), false),
						),
					),
					en.Recover(
						en.Switch(
							en.ToActions(
								en.Call(switchJobJudge),
								en.Call(switchJobCaseA),
								en.Call(switchJobCaseB),
								en.Call(switchJobCaseC),
							),
							WithName("switch"),
						),
					),
				),
			),
		),
	)

	attach := process.VisitJobs()
	for _, att := range attach {
		s := fmt.Sprintf("(%d, %s)", att.Index, att.Name)
		fmt.Println(s)
	}

	_, _ = process.Run()
}
