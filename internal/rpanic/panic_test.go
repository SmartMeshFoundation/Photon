package rpanic

import (
	"fmt"
	"testing"
)

func panic2() {
	panic("33")
}
func handlePanic() {
	if err := recover(); err != nil {
		fmt.Printf("err=%s", err)
	}
}

func TestPanic(t *testing.T) {
	defer handlePanic()
	panic2()
}
