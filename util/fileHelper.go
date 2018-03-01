package util

import (
	"io"
	"os"
	"path"
)

//获取配置文件完整路径
func GetConfigFile(fileName string) string {
	curDir := os.Args[0]
	return path.Dir(path.Dir(curDir)) + "/" + PATH_CONFIG + fileName
}

//获取输出文件完整路径
func GetOutPutFile(fileName string) string {
	curDir := os.Args[0]
	path := path.Dir(path.Dir(curDir)) + "/" + PATH_OUTPUT
	os.MkdirAll(path, 0755)
	return path + fileName
}

func MkDir(path string) bool {
	return os.MkdirAll(path, 0755) != nil
}

//获源码完整路径
func GetSourceCodePath() string {
	curDir := os.Args[0]
	return path.Dir(path.Dir(curDir)) + "/" + PATH_SOURCECODE
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, perm)
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

func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func PathExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

//目录下文件和文件夹数量
func FileCount(dirname string) int {
	f, err := os.Open(dirname)
	if err != nil {
		return 0
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return 0
	}
	return len(list)
}