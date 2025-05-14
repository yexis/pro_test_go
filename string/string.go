package string

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ToByteSlice ... convert string to bytes with no copy
func ToByteSlice(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

/*
CopyString
golang中的string是只读轻量的，字符串底层实际是一个结构体（感觉这一点和指针比较类似）
形如：

	type string struct {
	    Data *byte // 指向底层只读字节数组
	    Len  int
	}

传参时：仅复制 Data 和 Len（通常 16 字节），不会复制底层内容。
不可变性：Go 字符串是只读的，设计上防止修改，从而可以放心共享底层内存。
*/
func CopyString() {
	s := "12345"
	func(t string) {
		p := (*reflect.StringHeader)(unsafe.Pointer(&t))
		fmt.Println("__string_copy__ end", t, p.Data, p.Len)
	}(s)
	q := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Println("__string_copy__ start", s, q.Data, q.Len)
}
