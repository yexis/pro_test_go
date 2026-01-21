package engine

import (
	"errors"
	"fmt"
)

const (
	SkipSerialGroup       = 666
	ActionIsNotCanWait    = 667
	SwitchBeyond          = 668
	SwitchValNotSelection = 669
)

var codeMsg = map[int]string{
	SkipSerialGroup:       "skip serial group",
	ActionIsNotCanWait:    "action is not CanWait",
	SwitchBeyond:          "switch beyond",
	SwitchValNotSelection: "switch val is not Selection",
}

type InnerError struct {
	code int
	msg  string
}

func (e *InnerError) Code() int {
	return e.code
}

func (e *InnerError) Msg() string {
	return e.msg
}

func (e *InnerError) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.msg)
}

func FromCode(code int) *InnerError {
	var ie InnerError
	if msg, has := codeMsg[code]; has {
		ie.code = code
		ie.msg = msg
	} else {
		ie.code = -1
		ie.msg = "error occurred"
	}
	return &ie
}

func New(msg string) *InnerError {
	return &InnerError{
		code: -1,
		msg:  msg,
	}
}

func IsInnerError(e error) *InnerError {
	var ie *InnerError
	if errors.As(e, &ie) {
		return ie
	}
	return nil
}
