package main

import (
	"io/ioutil"
	"path"
	"strings"
)

//过滤出.c文件
func isCFile(fileName string) bool {
	ext := path.Ext(fileName)
	ext = strings.ToLower(ext)
	if ext == ".c" {
		return true
	}
	return false
}

func getFileList(dir string, list map[string]bool){
	dirList, _ := ioutil.ReadDir(dir)
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	for _, item := range dirList {
		name := item.Name()
		if item.IsDir() {
			if name != "." && name != ".." {
				getFileList(dir+item.Name(), list)
			}
		} else {
			if isCFile(name) {
				list[dir + name] = true
			}
		}
	}
}

