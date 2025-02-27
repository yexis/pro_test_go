package assertion

import (
	"fmt"
	"testing"
)

// TestAssert ... nil 不能被断言成 *struct
func TestAssert(t *testing.T) {
	type A struct {
		v int
	}
	var i interface{} = nil

	if _, ok := i.(*A); ok {
		fmt.Println("assert succeed")
	} else {
		fmt.Println("assert failed")
	}
}
