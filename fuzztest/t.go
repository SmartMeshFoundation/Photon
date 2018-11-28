package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	fmt.Println(b == 0x04)
}
