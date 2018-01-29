package main

import (
	"io/ioutil"
	"path"
	"strings"
)

//获取指定目录下所有.c和.h后缀的文件
func getFileList(dir string) []string {
	var queue Queue
	queue.Enqueue(dir)
	list := make([]string, 0, 100)

	for queue.Size() > 0 {
		dir = queue.Dequeue().(string)
		dirList, _ := ioutil.ReadDir(dir)
		if !strings.HasSuffix(dir, "/") {
			dir += "/"
		}

		for _, item := range dirList {
			name := item.Name()
			if item.IsDir() {
				if name != "." && name != ".." {
					queue.Enqueue(dir + item.Name())
				}
			} else {
				ext := strings.ToLower(path.Ext(name))
				if ext == ".c" || ext == ".h"{
					list = append(list, dir+name)
				}
			}
		}
	}

	return list
}
