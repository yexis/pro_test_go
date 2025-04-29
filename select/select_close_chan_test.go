package _select

import (
	"testing"
)

// panic
func TestSelectCloseChanBeforeSelect(t *testing.T) {
	for i := 0; i < 10000; i++ {
		SelectCloseChanBeforeSelect()
	}
}

// panic
func TestSelectCloseChanAsyncSelect(t *testing.T) {
	for i := 0; i < 10000; i++ {
		SelectCloseChanSyncSelect()
	}
}

// panic
func TestSelectCloseChanAsyncSelectByAtomic(t *testing.T) {
	for i := 0; i < 10000; i++ {
		SelectCloseChanSyncSelectByAtomic()
	}
}

func TestSelectCloseChanAsyncSelectByMutex(t *testing.T) {
	for i := 0; i < 10000; i++ {
		SelectCloseChanSyncSelectByMutex()
	}
}
