package withdefer

import (
	"errors"
	"fmt"
	"testing"
)

// defer 返回值捕获
func TestDeferModifyVal(t *testing.T) {
	handler := func() error {
		var err error
		defer func() {
			err = errors.New("the first defer error")
		}()
		defer func() {
			err = errors.New("the second defer error")
		}()
		err = errors.New("normal error")
		return err
	}
	err := handler()
	if err != nil {
		fmt.Println("final error:", err.Error())
	}
}

func TestNamedDeferModifyVal(t *testing.T) {
	handler := func() (err error) {
		defer func() {
			err = errors.New("the first defer error")
		}()
		defer func() {
			err = errors.New("the second defer error")
		}()
		err = errors.New("normal error")
		return err
	}
	err := handler()
	if err != nil {
		fmt.Println("final error:", err.Error())
	}
}

// defer 变量捕获
// 1. 变量捕获
// 写法一
// x := 10
//
//	defer func(int x) {
//	   fmt.Println(x)
//	}()
//
// 写法二
// x := 10
// defer fmt.Println(x)
// 两种写法的作用是一致的
//
// 2. 引用捕获（闭包）
// x := 10
//
//	defer func() {
//	   fmt.Println(x)
//	}
func TestNamedDeferModifyVal2(t *testing.T) {
	func() {
		x := 10
		defer fmt.Println("defer-val", x)
		x = 20
		fmt.Println("update-val", x)
	}()

	func() {
		x := 10
		defer func() {
			fmt.Println("defer-val", x)
		}()
		x = 20
		fmt.Println("update-val", x)
	}()
}
