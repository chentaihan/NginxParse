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

	cFileList, hFileList := getFileList(dir)

	defineList := make([]IParse, 0, len(ConfigInfo.ParseList))
	assignmentList := make([]IParse, 0, len(ConfigInfo.ParseList))
	macroList := make([]IParse, 1, 1)
	macroList[0] = NewParseMacro()
	for i, _ := range ConfigInfo.ParseList {
		defineList = append(defineList, NewParseDefine(&ConfigInfo.ParseList[i]))
		assignmentList = append(assignmentList, NewAssignment(&ConfigInfo.ParseList[i]))
	}

	//解析宏
	parseMacro(macroList, hFileList)
	//解析结构体定义
	parseDefine(defineList, hFileList)
	//解析结构体变量
	parseAssignment(assignmentList, cFileList)

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
	//fmt.Println("------------------------------------------")
	//fmt.Println("------------------------------------------")
	//for _, sut := range structList {
	//	if mgr, ok := sut.(*Macro); ok {
	//		list := mgr.MacroList
	//		for key, val := range list {
	//			if val.Name != "" {
	//				fmt.Printf("%s%s=%s", key, val.Name, val.Value)
	//				fmt.Println()
	//			} else {
	//				fmt.Printf("-------------------------%s%s=%s", key, val.Name, val.Value)
	//				fmt.Println()
	//			}
	//		}
	//	}
	//}
	//
	//macro := GetMacro()
	//fmt.Println(macro.Exist("NGX_MODULE_SIGNATURE_5"))
	//fmt.Println(macro.Exist("NGX_MODULE_SIGNATURE_511"))
	//fmt.Println(macro.GetMacroValue("NGX_MODULE_SIGNATURE_12"))
	//fmt.Println(macro.GetMacroValue("NGX_MODULE_SIGNATURE_1211"))
	//fmt.Println(macro.GetMacroValue("ngx_conf_merge_uint_value(conf__, prev__, default__)"))
	//fmt.Println(macro.GetMacroValue("ngx_conf_merge_uint_value(conf__, prev__)"))
	//fmt.Println(macro.GetMacroValue("ngx_conf_merge_uint_value()"))

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
