package main

import (
	"io"
	"os"
	"path"
)

//获取配置文件完整路径
func getConfigFile(fileName string) string {
	curDir := os.Args[0]
	return path.Dir(path.Dir(curDir)) + "/" + PATH_CONFIG + fileName
}

//获取输出文件完整路径
func getOutPutFile(fileName string) string {
	curDir := os.Args[0]
	path := path.Dir(path.Dir(curDir)) + "/" + PATH_OUTPUT
	os.MkdirAll(path, 0777)
	return path + fileName
}

func WriteFileAppend(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
