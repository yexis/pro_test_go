package pointer

import (
	"bytes"
	logger "icode.baidu.com/baidu/duer/go-utils/trlogger"
	"testing"
	"time"
)

type A struct {
	id  int
	val string
}

func TestPointer(t *testing.T) {
	pa := &A{
		id:  1,
		val: "lgy",
	}
	var pa2 *A
	logger.Init("/Users/liguoyang/go/src/pro_test_go/pointer")
	lg := logger.NewTraceLogger("123", "agent")
	lg.Debug("__pointer__ pa:%+v", pa)
	lg.Debug("__pointer__ pa2:%+v", pa2)

	s := "liguoyang"
	ps := &s
	lg.Debug("__pointer__ ps:%+v", ps)

	bytes.Buffer{}
	time.Sleep(5 * time.Second)
}
