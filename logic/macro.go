package logic

/**
宏处理
 */

import (
	"strings"

	"github.com/chentaihan/NginxParse/util"
)

type Macro struct {
	Name  string
	Value string
}

func NewMacro() *Macro {
	return &Macro{}
}

//其中的宏直接返回实参
var returnActualVal = []string{
	"ngx_string",
}

//判断是不是有效宏
func (macro *Macro) IsStartStruct(line string) bool {
	return strings.HasPrefix(line, util.NGX_DEFINE)
}

//解析宏
func (macro *Macro) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	writer = macro.FormatStruct(writer)
	line := writer.ToString()
	index := strings.Index(line, util.NGX_DEFINE)
	if index >= 0 {
		line = line[index+1+util.NGX_DEFINE_LEN:]
		key := util.GetLegalString(line)
		index = strings.Index(line, key) + len(key)
		line = line[index:]
		line = strings.Trim(line, " ")
		var name string = ""
		var value string = ""
		if len(line) > 0 {
			if line[0] == '(' { //宏有参数
				end := strings.Index(line, ")") + 1
				name = line[0:end]
				line = line[end:]
				index = end
				name = strings.Trim(name, " ")

			}
			value = line
			value = strings.Trim(value, " ")
		}
		value = strings.Trim(value, "\n")
		GetMacros().Add(key, &Macro{name, value})
	}
	return true
}

//格式化宏
func (macro *Macro) FormatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := util.NewBufferWriter(bufWriter.Size())
	for _, val := range inBuf {
		if val != '\\' {
			outBuf.WriteChar(val)
		}
	}
	return outBuf
}

//宏不是已\\结尾
func (macro *Macro) IsEndStruct(line string) bool {
	if strings.HasSuffix(line, "\\") {
		return false
	}
	return true
}
