package main

import (
	"testing"
)

func TestIsChoosedOne(t *testing.T) {
	if isChoosedOne(true, true, false) {
		t.Fatal("isChoosedOne(true,true,false) must be false, but true")
	}
	if !isChoosedOne(true, false, false, false) {
		t.Fatal("isChoosedOne(true,false,false,false) must be true, but false")
	}
	if isChoosedOne(false, false, false) {
		t.Fatal("isChoosedOne(false,false,false) must be false,but true")
	}
}
