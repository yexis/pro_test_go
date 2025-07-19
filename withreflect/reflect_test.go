package withreflect

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string
}

func TestReflect(t *testing.T) {
	p := &Person{Name: "Tom"}
	val := reflect.ValueOf(p) // withreflect.Value 表示一个 *Person
	typ := val.Type()         // withreflect.Type 表示类型: *main.Person
	kind := typ.Kind()        // withreflect.Kind 是 Ptr（指针）

	typElem := typ.Elem()
	typElemKind := typElem.Kind()

	valElem := val.Elem() // 获取 p 指向的值，即 Person
	valElemKind := valElem.Kind()

	elemType := valElem.Type()  // 类型: main.Person
	elemKind := elemType.Kind() // withreflect.Kind 是 Struct（结构体）

	/*
		p         = &{Tom}
		val       = &{Tom}
		typ       = *withreflect.Person
		kind      = ptr
		elemVal   = {Tom}
		elemType  = withreflect.Person
		elemKind  = struct
	*/
	fmt.Println("p         =", p)
	fmt.Println("val       =", val)
	fmt.Println("typ       =", typ)
	fmt.Println("kind      =", kind)

	fmt.Println("typElem     =", typElem)
	fmt.Println("typElemKind =", typElemKind)
	fmt.Println("valElem     =", valElem)
	fmt.Println("valElemKind =", valElemKind)

	fmt.Println("elemType  =", elemType)
	fmt.Println("elemKind  =", elemKind)
}
