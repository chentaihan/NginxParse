package main

import (
	"html/template"
)

type CommandWriter struct {
	//Write(p []byte) (n int, err error)
	buf []byte
}

func NewCommandWriter() *CommandWriter {
	writer := &CommandWriter{}
	writer.buf = make([]byte, 0, 1024*16)
	return writer
}

func (writer *CommandWriter) Write(p []byte) (n int, err error) {
	writer.buf = append(writer.buf, p...)
	n = len(p)
	err = nil
	return n, err
}

const (
	MODULENAME       = "<tr><th colspan=\"6\">{{.ModuleName}}</th></tr>"
	FILE_STRUCT_NAME = "<tr><td colspan=\"3\">文件名：{{.FileName}}</td><td colspan=\"3\">struct：{{.StructName}}</td></tr>"
	COMMANDLINE_TH   = "<tr><td width=\"16%\">name</td><td width=\"24%\">type</td><td width=\"17%\">set</td><td width=\"17%\">conf</td><td width=\"13%\">offset</td><td width=\"13%\">post</td></tr>"
	COMMANDLINE_TD   = "<tr><td>{{.Name}}</td><td>{{.TypeVal}}</td><td>{{.Set}}</td><td>{{.Conf}}</td><td>{{.Offset}}</td><td>{{.Post}}</td></tr>"
)

func commandInfoFormat(cmdInfo *CommandInfo) []byte {
	writer := NewCommandWriter()

	writer.Write([]byte("<table>"))
	tmpl := template.New("tmpl1")
	tmpl.Parse(MODULENAME + FILE_STRUCT_NAME)
	tmpl.Execute(writer, cmdInfo)
	writer.Write([]byte(COMMANDLINE_TH))

	for _, cmd := range cmdInfo.CmdList {
		tmpl = template.New("tmpl1")
		tmpl.Parse(COMMANDLINE_TD)
		tmpl.Execute(writer, cmd)
	}
	writer.Write([]byte("</table><br/><br/>"))

	return writer.buf
}
