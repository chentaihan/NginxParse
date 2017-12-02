package main

import (
	"io"
	"os"
)

func isDir(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
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

func main() {

	//if len(os.Args) != 2 {
	//	fmt.Println("need one param as nginx source path")
	//	return
	//}
	//dir := os.Args[1]
	//if !isDir(dir) {
	//	fmt.Println(dir, " is not a path")
	//	return
	//}
	dir := "/Users/didi/OpenSource/nginx-1.12.2/src"
	initParse()

	fileList := make(map[string]bool, 0)
	getFileList(dir, fileList)
	cmdList := make([]*CommandInfo, 0, 50)
	for fileName, _ := range fileList {
		commandInfo := parseFile(fileName)
		if commandInfo != nil {
			cmdList = append(cmdList, commandInfo)
		}
	}
	outputFile := "./nginxConf.txt"
	os.Remove(outputFile)
	for _, cmdInfo := range cmdList {
		buf := commandInfoFormat(cmdInfo)
		WriteFile(outputFile, buf, 0666)
	}

}
