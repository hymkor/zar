package stringstack

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	//println("start push")
	var stack Stack
	for i := 0; i < 4000; i++ {
		//println(i)
		stack.PushString(fmt.Sprintf("%d", i))
	}
	//println("start pop")
	for i := 4000 - 1; i >= 0; i-- {
		//println(i)
		s, ok := stack.PopString()
		if !ok {
			t.Fatal("stack too short")
			return
		}
		if s != fmt.Sprintf("%d", i) {
			t.Fatalf("data diff %d and %s", i, s)
			return
		}
	}
	_, ok := stack.PopString()
	if ok {
		t.Fatal("could not find stack end")
		return
	}
}
