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

	LoadConfig()
	dir := "/Users/didi/OpenSource/nginx-1.12.2/src"

	fileList := getFileList(dir)


	defineList := make([]IParse, 0, len(ConfigInfo.ParseList))
	assignmentList := make([]IParse, 0, len(ConfigInfo.ParseList))
	macroList := make([]IParse, 1, 1)
	macroList[0] = GetMacro()
	for i, _ := range ConfigInfo.ParseList {
		defineList = append(defineList, NewParseDefine())
		assignmentList = append(assignmentList, NewAssignment(&ConfigInfo.ParseList[i]))
	}

	//解析宏
	parseMacro(macroList, fileList)

	//fileList = []string{
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/core/ngx_buf.h",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/ngx_http_geo_module.c",
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream.h",
	//}

	//解析结构体定义
	parseDefine(defineList, fileList)
	//解析结构体变量
	//parseAssignment(assignmentList, fileList)

	fmt.Println("---------------------result output------------------")
	sstMgr := *GetStructManager()
	for _, sstInfo := range sstMgr{
		//fmt.Println(sttName)
		fmt.Println(sstInfo.StructString)
	}

	fmt.Println("ok")
}

func parseMacro(structList []IParse, fileList []string) {
	fileParse := FileParse{}
	for _, sut := range structList {
		fileParse.Register(sut)
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}

}

func parseAssignment(structList []IParse, fileList []string) {
	fileParse := FileParse{}
	for _, sut := range structList {
		fileParse.Register(sut)
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}
	for _, sut := range structList {
		if mgr, ok := sut.(*Assignment); ok {
			OutPut(mgr)
		}
	}
}

func parseDefine(structList []IParse, fileList []string) {
	fileParse := FileParse{}
	for i, _ := range structList {
		fileParse.Register(structList[i])
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}

}
