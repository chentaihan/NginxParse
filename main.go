package main

import (
	"fmt"
	"os"

	"github.com/chentaihan/NginxParse/logic"
	"github.com/chentaihan/NginxParse/util"
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

	util.LoadConfig()
	dir := "/Users/didi/OpenSource/nginx-1.12.2/src"

	fileList := util.GetFileList(dir)

	defineList := make([]logic.IParse, 0, len(util.ConfigInfo.ParseList))
	assignmentList := make([]logic.IParse, 0, len(util.ConfigInfo.ParseList))
	macroList := make([]logic.IParse, 1, 1)
	macroList[0] = logic.GetMacro()
	for i, _ := range util.ConfigInfo.ParseList {
		defineList = append(defineList, logic.NewParseDefine())
		assignmentList = append(assignmentList, logic.NewAssignment(&util.ConfigInfo.ParseList[i]))
	}

	//解析宏
	parseMacro(macroList, fileList)

	//fileList = []string{
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/core/ngx_buf.h",
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/ngx_http_geo_module.c",
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream.h",
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/stream/ngx_stream_geoip_module.c",
	//	"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream_round_robin.h",
	//	"/usr/local/Cellar/go/1.9.2/src/github.com/chentaihan/NginxParse/output/test.c",
	//}

	//解析结构体定义
	parseDefine(defineList, fileList)
	//解析结构体变量
	//parseAssignment(assignmentList, fileList)

	fmt.Println("---------------------result output------------------")
	//sstMgr := *logic.GetStructManager()
	//for _, sstInfo := range sstMgr{
	//	//fmt.Println(sttName)
	//	fmt.Println(sstInfo.StructString)
	//}

	fmt.Println("ok")
}

func parseMacro(structList []logic.IParse, fileList []string) {
	fileParse := logic.FileParse{}
	for _, sut := range structList {
		fileParse.Register(sut)
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}
}

func macroOutput() {

}

func parseAssignment(structList []logic.IParse, fileList []string) {
	fileParse := logic.FileParse{}
	for _, sut := range structList {
		fileParse.Register(sut)
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}
	for _, sut := range structList {
		if mgr, ok := sut.(*logic.Assignment); ok {
			logic.OutPut(mgr)
		}
	}
}

func parseDefine(structList []logic.IParse, fileList []string) {
	fileParse := logic.FileParse{}
	for i, _ := range structList {
		fileParse.Register(structList[i])
	}
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}
}
