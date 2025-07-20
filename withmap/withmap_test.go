package withmap

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	mp := make(map[int]int)
	mp[0] = 2
	SetValToMap(mp)
	for k, v := range mp {
		fmt.Println(k, v)
	}
}
