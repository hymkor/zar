package stringstack

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func Heap() uint64 {
	var ms runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&ms)

	return ms.Alloc
}

func TestAll(t *testing.T) {
	baseHeap := Heap()

	var stack Stack
	var pureLen uint64 = 0
	for i := 0; i < 4000; i++ {
		//println(i)
		s := fmt.Sprintf("%d", i)
		stack.Push(s)
		pureLen += uint64(len(s) + 2)
	}
	println("Total String Length=", pureLen)
	println("Overhead(stringstack)=", Heap()-baseHeap-pureLen)

	var buffer1 strings.Builder
	for i := 4000 - 1; i >= 0; i-- {
		//println(i)
		ok := stack.PopTo(&buffer1)
		if !ok {
			t.Fatal("stack too short")
			return
		}
		s := buffer1.String()
		if s != fmt.Sprintf("%d", i) {
			t.Fatalf("data diff %d and %s", i, s)
			return
		}
		buffer1.Reset()
	}
	ok := stack.PopTo(&buffer1)
	if ok {
		t.Fatal("could not find stack end")
		return
	}
	baseHeap = Heap()

	buffer := []string{}
	for i := 0; i < 4000; i++ {
		buffer = append(buffer, fmt.Sprintf("%d", i))
	}
	println("Overhead(string slice)=", Heap()-baseHeap-pureLen)
	for i := 4000 - 1; i >= 0; i-- {
		if buffer[i] != fmt.Sprintf("%d", i) {
			t.Fatal("buffer broken?")
		}
	}
}
