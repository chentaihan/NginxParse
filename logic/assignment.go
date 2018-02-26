package logic

/*
结构体实例
 */

import (
	"strings"

	"github.com/chentaihan/NginxParse/util"
	"fmt"
)

type Assignment struct {
	checker Checker
}

func NewAssignment() *Assignment {
	varMgr := &Assignment{
		checker: Checker{
			Include:   []string{"=", "{"},
			IncludeNo: []string{"(", ")"},
		},
	}
	return varMgr
}

//判断是不是有效结构体
func (varMgr *Assignment) IsStartStruct(line string) bool {
	return varMgr.checker.Check(line)
}

//struct 赋值以;结尾
func (mgr *Assignment) IsEndStruct(line string) bool {
	isEnd := strings.Index(line, ";")
	if isEnd < 0 {
		return false
	}
	return true
}

//解析出ModuleInfo
func (mgr *Assignment) parseTableInfo(filePath string, writer *util.BufferWriter) *TableInfo {
	table := &TableInfo{}
	writer.MoveNext()
	array := util.GetLegalStrings(writer.Current())
	if len(array) >= 2 {
		table.VarName = array[len(array)-1]
		if "ngx_event_core_module_ctx" == table.VarName {
			fmt.Println(writer.ToString())
		}
		structName := array[len(array)-2]
		sttInfo := GetDefines().Get(structName)
		if sttInfo != nil {
			table.StructInfo = *sttInfo
		} else {
			util.Println("错误：struct:", structName, " not exist")
			return nil
		}
	}

	table.Content = make([][]string, 0, 1)
	return table
}

//解析出结构体内容
func (mgr *Assignment) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	writer = mgr.FormatStruct(writer)
	table := mgr.parseTableInfo(filePath, writer)
	if table == nil {
		return false
	}
	util.Println(writer.ToString())
	writer.Reset()
	writer.MoveNext()
	structHeader := writer.Current()
	//是否是结构体数组
	isArray := strings.Index(structHeader, "[") > 0
	fieldsSize := len(table.Fields) + 2
	lines := make([]string, 0, fieldsSize)
	if !isArray {
		lines = append(lines, "{")
	}

	for !writer.IsEnd() {
		for len(lines) < fieldsSize && writer.MoveNext() {
			line := strings.TrimRight(writer.Current(), "\n")
			lines = append(lines, line)
		}
		if len(lines) < fieldsSize {
			break
		}
		fields := mgr.parseFields(lines)
		if len(fields) > 0 {
			table.Content = append(table.Content, fields)
		} else {
			fmt.Println("错误 ", writer.ToString())
			//panic(writer.ToString())
			break
		}
		lines = lines[0:0]
	}

	GetAssignments().Add(table)
	return true
}

func (mgr *Assignment) parseFields(lines []string) []string {
	fieldsLen := len(lines)
	if fieldsLen > 2 {
		lines = lines[1:fieldsLen-1]
		fields := make([]string, 0, fieldsLen-2)
		for _, val := range lines {
			fields = append(fields, val)
		}
		return fields
	}
	return nil
}

//将buffer中的struct赋值格式化成容易解析的样子
func (mgr *Assignment) formatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := util.NewBufferWriter(bufWriter.Size() + 64)
	bracketCount := 0
	inBracketCount := strings.Count(bufWriter.ToString(), "[") + 1 //大括号深度
	inLittleBracket := false                                       //在小括号内部

	const BEFORE_LINE = 1
	const AFTER_LINE = 2
	const BEFORE_AFTER_LINE = 3

	for index, val := range inBuf {
		//重复的空格合并
		if val == ' ' && index < len(inBuf)-1 && inBuf[index+1] == ' ' {
			continue
		}
		//结构体第一行空格保留
		if val == ' ' && bracketCount > 0 {
			continue
		}
		newLine := 0
		addChar := true
		if val == '{' {
			newLine = AFTER_LINE
			bracketCount++
		} else if val == '}' {
			newLine = BEFORE_AFTER_LINE
			bracketCount--
		}

		if val == '(' {
			inLittleBracket = true
		} else if val == ')' {
			inLittleBracket = false
		}

		if val == ',' && !inLittleBracket {
			if bracketCount == inBracketCount {
				newLine = AFTER_LINE
			}
			//}后面的,去掉
			//if inBuf[index-1] == '}' {
			addChar = false
			//}
		}
		isAddN := false
		if newLine&BEFORE_LINE > 0 && index > 0 && inBuf[index-1] != '\n' {
			outBuf.WriteChar('\n')
			isAddN = true
		}
		if addChar {
			outBuf.WriteChar(val)
		}
		if !isAddN && newLine&AFTER_LINE > 0 && inBuf[index+1] != '\n' {
			outBuf.WriteChar('\n')
		}
	}
	return outBuf
}

//将buffer中的struct赋值格式化成容易解析的样子
func (mgr *Assignment) FormatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	outBuf := mgr.formatStruct(bufWriter)
	outBuf = mgr.replaceMacro(outBuf)
	outBuf = mgr.formatStruct(outBuf)
	return outBuf
}

func (mgr *Assignment) replaceMacro(bufWriter *util.BufferWriter) *util.BufferWriter {
	bufWriter.Reset()
	outBuf := util.NewBufferWriter(bufWriter.Size() + 64)
	for bufWriter.MoveNext() {
		line := bufWriter.Current()
		macroName := strings.Trim(line, "\n")
		value := GetMacros().GetMacroValue(macroName)
		if value != "" {
			outBuf.WriteString(value)
			outBuf.WriteChar('\n')
		} else {
			outBuf.WriteString(line)
		}
	}
	return outBuf
}
