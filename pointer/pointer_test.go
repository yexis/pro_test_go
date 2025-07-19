package pointer

import (
	"fmt"
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

	time.Sleep(5 * time.Second)
}

func TestNilPointerButHaveType(t *testing.T) {
	type Str struct {
		v int
	}
	var i interface{} = (*Str)(nil)

	if v, ok := i.(*Str); ok {
		fmt.Println("assert success", v)
	} else {
		fmt.Println("assert failed", v)
	}
}

func TestPrintPointer(t *testing.T) {
	a := 10
	p := &a
	q := &a
	r := &a
	fmt.Printf("p: %p\n", p)
	fmt.Printf("q: %p\n", q)
	fmt.Printf("r: %p\n", r)

}
