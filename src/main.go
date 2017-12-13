package main

import (
	"fmt"
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

	fileParse := FileParse{}

	cmdMgr := NewCommandManager()
	cmdStruct := StructParse{
		structType: STRUCT_TYPE_COMMAND,
		Parse:      cmdMgr,
	}
	fileParse.Register(&cmdStruct)

	moduleMgr := NewModuleManager()
	moduleStruct := StructParse{
		structType: STRUCT_TYPE_MODULE,
		Parse:      moduleMgr,
	}
	fileParse.Register(&moduleStruct)

	fileList := getFileList(dir)

	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}

	outPutCommand(cmdMgr.CmdInfo)
	outPutModule(moduleMgr.moduleInfo)
	fmt.Println("ok")
}
