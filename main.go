package main

import (
	"os"

	"github.com/chentaihan/NginxParse/logic"
	"github.com/chentaihan/NginxParse/util"
	"fmt"
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
	//	util.Println("need one param as nginx source path")
	//	return
	//}

	//dir := os.Args[1]
	//if !isDir(dir) {
	//	util.Println(dir, " is not a path")
	//	return
	//}

	util.LoadConfig()
	dir := "/Users/didi/OpenSource/nginx-1.12.2/src"

	fileList := util.GetFileList(dir)

	parseList := []logic.IParse{
		logic.NewMacro(),
		logic.NewParseDefine(),
		logic.NewAssignment(),
	}
	fileParse := logic.FileParse{}
	for _, parse := range parseList {
		fileParse.Register(parse)
	}

	//解析宏
	parseMacro(fileList)

	//解析结构体定义
	parseDefine(fileList)

	//fileList = []string{
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/core/ngx_buf.h",
	//    //"/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/ngx_http_geo_module.c",
	//    //"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream.h",
	//    //"/Users/didi/OpenSource/nginx-1.12.2/src/stream/ngx_stream_geoip_module.c",
	//    //"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream_round_robin.h",
	//    "/usr/local/Cellar/go/1.9.2/src/github.com/chentaihan/NginxParse/output/test.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/core/nginx.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/stream/ngx_stream.h",
	//}

	//解析结构体变量
	parseAssignment(fileList)

	fmt.Println("ok")
}

func parseMacro(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewMacro())
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			break
		}
	}
	logic.GetMacros().Print()
}

func parseDefine(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewParseDefine())
	for _, fileName := range fileList {
		if fileParse.Parse(fileName) {
			//break
		}
	}
	logic.GetDefines().Print()
}

func parseAssignment(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewAssignment())
	for _, fileName := range fileList {
		fmt.Println("fileName = ", fileName)
		if fileParse.Parse(fileName) {
			//break
		}
	}
	logic.GetAssignments().Print()
}
