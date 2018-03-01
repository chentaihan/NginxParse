package main

import (
	"fmt"

	"github.com/chentaihan/NginxParse/logic"
	"github.com/chentaihan/NginxParse/util"
)

func main() {
	err := util.LoadConfig()
	if err != nil {
		panic("err")
	}

	downloadSourceCode()

	fileList := util.GetFileList(util.GetSourceCodePath())

	//fileList = []string{
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/core/ngx_buf.h",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/ngx_http_geo_module.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream.h",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/stream/ngx_stream_geoip_module.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/http/ngx_http_upstream_round_robin.h",
	//	"/usr/local/Cellar/go/1.9.2/src/github.com/chentaihan/NginxParse/output/test.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/core/nginx.c",
	//	//"/Users/didi/OpenSource/nginx-1.12.2/src/stream/ngx_stream.h",
	//}

	//解析析宏
	parseMacro(fileList)

	//解析结构体定义 和 结构体重定义
	parseDefine(fileList)

	//解析结构体变量
	parseAssignment(fileList)

	fmt.Println("nginx source code parse success")
}

//解析析宏
func parseMacro(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewMacro())
	for _, fileName := range fileList {
		fileParse.Parse(fileName)
	}
	logic.GetMacros().Print()
}

//解析结构体定义 和 结构体重定义
func parseDefine(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewParseDefine())
	fileParse.Register(logic.NewTypedef())
	for _, fileName := range fileList {
		fileParse.Parse(fileName)
	}

	logic.GetTypedefs().Parse(logic.GetDefines())
	logic.GetTypedefs().Print()
	logic.GetDefines().Print()
}

//解析结构体变量
func parseAssignment(fileList []string) {
	fileParse := logic.FileParse{}
	fileParse.Register(logic.NewAssignment())
	for _, fileName := range fileList {
		fileParse.Parse(fileName)
	}
	logic.GetAssignments().Print()
}

//第一次运行的时候，会去下载nginx源码，并解压
func downloadSourceCode() {
	sourceCodePath := util.GetSourceCodePath()
	util.MkDir(sourceCodePath)
	zipFilePath := sourceCodePath + util.NGINX_ZIP
	if !util.FileExists(zipFilePath) {
		if util.Downalod(util.ConfigInfo.SourceCodeUrl, zipFilePath) {
			util.UnZip(zipFilePath, sourceCodePath)
		}
	} else {
		if util.FileCount(sourceCodePath) < 2 {
			util.UnZip(zipFilePath, sourceCodePath)
		}
	}
}
