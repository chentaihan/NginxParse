package main

import (
	"fmt"
	"encoding/json"
)

func main() {

	initParse()

	dir := "/Users/didi/OpenSource/nginx-1.12.2/src/http/modules/"
	fileList := make([]string, 0)
	fileList = getFileList(dir, fileList)
	cmdList := make([]*CommandInfo, 0, 50)

	for _,fileName := range fileList {
		commandInfo := parseFile(fileName)
		if commandInfo != nil {
			cmdList = append(cmdList, commandInfo)
		}else{
			fmt.Println(fileName)
		}
	}
	buf,_ := json.Marshal(cmdList)
	fmt.Println(string(buf))
}
