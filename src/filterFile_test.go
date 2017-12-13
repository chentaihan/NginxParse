package main

import (
	"testing"
)

func Test_isCFile(t *testing.T) {
	isCFile1 := isCFile("fdfdfjdkf.c")
	if !isCFile1 {
		t.Error("isCFile error")
	}

	isCFile1 = isCFile("fdjfdfkjdkfdkf.h")
	if isCFile1 {
		t.Error("isCFile error")
	}
}
