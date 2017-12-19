package main

import (
	"io/ioutil"
	"path"
	"strings"
)

//获取指定目录下所有.C后缀的文件
func getFileList(dir string) ([]string, []string) {
	queue := NewQueue()
	queue.Enqueue(dir)
	cList := make([]string, 0, 100)
	hList := make([]string, 0, 100)

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
				ext := path.Ext(name)
				ext = strings.ToLower(ext)
				if ext == ".c" {
					cList = append(cList, dir+name)
				} else if ext == ".h" {
					hList = append(hList, dir+name)
				}
			}
		}
	}

	return cList, hList
}
