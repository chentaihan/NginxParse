package main

import (
	"strings"
)

type MacroInfo struct {
	Name  string
	Value string
}

type Macro struct {
	MacroList map[string]*MacroInfo
}

var macro *Macro = nil

func GetMacro() *Macro {
	if macro == nil {
		macro = &Macro{
			MacroList: make(map[string]*MacroInfo, 1024),
		}
	}
	return macro
}

//判断是不是有效宏
func (macro *Macro) IsStartStruct(line string) bool {
	return strings.HasPrefix(line, NGX_DEFINE)
}

func (macro *Macro) AddMacroInfo(key string, mf *MacroInfo){
	macro.MacroList[key] = mf
}

//解析宏
func (macro *Macro) ParseStruct(filePath string, writer *BufferWriter) bool {
	writer = macro.FormatStruct(writer)
	//fmt.Println(writer.ToString())
	line := writer.ToString()
	index := strings.Index(line, NGX_DEFINE)
	if index >= 0 {
		line = line[index+1+NGX_DEFINE_LEN:]
		key := getLegalString(line)
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

		macro.AddMacroInfo(key, &MacroInfo{name, value})
	}
	return true
}

//格式化宏
func (macro *Macro) FormatStruct(bufWriter *BufferWriter) *BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := NewBufferWriter(bufWriter.Size)
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

//是否存在指定的宏
func (macro *Macro) Exist(macroName string) bool {
	if _, ok := macro.MacroList[macroName]; ok {
		return true
	}
	return false
}

func (macro *Macro) GetMacroValue(macroName string) string {
	index := strings.Index(macroName, "(")
	key := macroName
	if index > 0 {
		key = macroName[0:index]
	}

	macroInfo, ok := macro.MacroList[key]
	if !ok {
		return ""
	}

	value := macroInfo.Value

	if index > 0 {
		actualName := macroName[index:]
		formalParams := getLegalStrings(macroInfo.Name)
		actualParams := getLegalStrings(actualName)
		return macro.replaceParams(value, formalParams, actualParams)
	}

	return value
}

//宏替换，实参代替形参
func (macro *Macro) replaceParams(value string, formalParams, actualParams []string) string {
	minLen := len(formalParams)
	if minLen > len(actualParams) {
		minLen = len(actualParams)
	}
	for i := 0; i < minLen; i++ {
		value = strings.Replace(value, formalParams[i], actualParams[i], -1)
	}
	return value
}
