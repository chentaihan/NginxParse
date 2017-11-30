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

func Test_getFileList(t *testing.T) {
	dir := "/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/"
	fileList := make([]string, 0)
	getFileList(dir, fileList)
}
