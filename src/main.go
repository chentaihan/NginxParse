package main

import (
	"fmt"
	"os"
)

func isDir(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
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

	structList := []*StructParse{
		&StructParse{
			structType: STRUCT_TYPE_COMMAND,
			Parse:      NewCommandManager(),
		},
		&StructParse{
			structType: STRUCT_TYPE_MODULE,
			Parse:      NewModuleManager(),
		},
		&StructParse{
			structType: STRUCT_TYPE_VARIABLE,
			Parse:      NewVariableManager(),
		},
	}

	fileParse := FileParse{}
	for _, sut := range structList {
		fileParse.Register(sut)
	}

	fileList := getFileList(dir)

	for _, fileName := range fileList {
		//fileName = "/Users/didi/OpenSource/nginx-1.12.2/src/http/v2/ngx_http_v2_module.c"
		if fileParse.Parse(fileName) {
			//break
		}
	}

	for _, sut := range structList {
		OutPut(sut.Parse)
	}

	fmt.Println("ok")
}
