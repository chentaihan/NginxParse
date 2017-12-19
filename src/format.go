package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	COMMAND_HEADER = "<tr><th colspan=\"6\">{{.ModuleName}}</th></tr><tr><td colspan=\"3\">文件名：{{.FileName}}</td><td colspan=\"3\">struct：{{.StructName}}</td></tr>"
	COMMAND_TH     = "<tr><td width=\"16%\">name</td><td width=\"24%\">type</td><td width=\"17%\">set</td><td width=\"17%\">conf</td><td width=\"13%\">offset</td><td width=\"13%\">post</td></tr>"
	COMMAND_TD     = "<tr><td>{{.Name}}</td><td>{{.Type}}</td><td>{{.Set}}</td><td>{{.Conf}}</td><td>{{.Offset}}</td><td>{{.Post}}</td></tr>"
	COMMAND_HTML   = "command.html"
)

const (
	MODULE_HEADER = "<tr><th colspan=\"6\">{{.ModuleName}}</th></tr><tr><td colspan=\"3\">文件名：{{.FileName}}</td><td colspan=\"3\">struct：{{.StructName}}</td></tr>"
	MODULE_TH     = "<tr><td width=\"33%\" colspan=\"2\">Context</td><td width=\"34%\" colspan=\"2\">Command</td><td width=\"33%\" colspan=\"2\">Type</td></tr>"
	MODULE_TD     = "<tr><td colspan=\"2\">{{.Context}}</td><td colspan=\"2\">{{.Command}}</td><td colspan=\"2\">{{.Type}}</td></tr>"
	MODULE_HTML   = "module.html"
)

const (
	VARIABLE_HEADER = "<tr><th colspan=\"%d\">%s</th></tr><tr><td colspan=\"%d\">文件名：%s</td><td colspan=\"%d\">struct：%s</td></tr>"
	VARIABLE_HTML   = "variable.html"
)

type formatStruct func(varList []*TableInfo) []byte

var nginxSourcePath string = "" //nginx源码路径

func removeSourcePath(filePath string) string {
	if nginxSourcePath == "" {
		nginxSourcePath = os.Args[0]
		if !strings.HasSuffix(nginxSourcePath, "/") {
			nginxSourcePath += "/"
		}
	}
	return filePath[len(nginxSourcePath)+1:]
}

func formatVariableList(varList []*TableInfo) []byte {
	buf := make([]byte, 0, 1024)

	for _, variable := range varList {
		variable.FileName = removeSourcePath(variable.FileName)
		buf = append(buf, formatTable(variable)...)
	}
	return buf
}

func formatHeader(variable *TableInfo) string {
	colspan1 := len(variable.Title)
	colspan3 := colspan1 / 2
	colspan2 := colspan1 - colspan3
	return fmt.Sprintf(VARIABLE_HEADER, colspan1, variable.ModuleName, colspan2, variable.FileName, colspan3,variable.StructName)
}

func formatTH(title []string) string {
	writer := NewBufferWriter(0)
	writer.Write([]byte("<tr>"))
	for i, _ := range title {
		td := fmt.Sprintf("<td >%s</td>", title[i])
		writer.WriteString(td)
	}
	writer.Write([]byte("</tr>"))
	return writer.ToString()
}

func formatTD(content [][]string) string {
	writer := NewBufferWriter(0)
	for _, tr := range content {
		writer.Write([]byte("<tr>"))
		for _, td := range tr {
			tdStr := fmt.Sprintf("<td >%s</td>", td)
			writer.WriteString(tdStr)
		}
		writer.Write([]byte("</tr>"))
	}

	return writer.ToString()
}

func formatTable(variable *TableInfo) []byte {
	writer := NewBufferWriter(0)
	writer.Write([]byte("<table>"))

	writer.WriteString(formatHeader(variable))
	writer.WriteString(formatTH(variable.Title))
	writer.WriteString(formatTD(variable.Content))

	writer.Write([]byte("</table><br/><br/>"))

	return writer.GetBuffer()
}

func OutPut(structInfo *Assignment) {
	fileName := structInfo.StructInfo.StructName + ".html"
	outputFile(FILE_CONFIG_FORMAT, fileName, formatVariableList, structInfo.Tables)
}

func outputFile(formatFile, outputFile string, formatFunc formatStruct, list []*TableInfo) {
	configFormat := getConfigFile(formatFile)
	content, err := ioutil.ReadFile(configFormat)
	if err != nil {
		return
	}
	var index int = -1
	for i := 0; i < len(content); i++ {
		if content[i] == '$' && content[i+1] == '$' {
			index = i
		}
	}

	outputFile = getOutPutFile(outputFile)
	ioutil.WriteFile(outputFile, content[:index], 0777)

	buf := formatFunc(list)
	WriteFileAppend(outputFile, buf, 0666)

	index += 2
	WriteFileAppend(outputFile, content[index:], 0666)
}
