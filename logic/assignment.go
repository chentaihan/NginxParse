package logic

/*
结构体实例
*/

import (
	"fmt"
	"strings"

	"github.com/chentaihan/NginxParse/util"
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
func (asst *Assignment) IsStartStruct(line string) bool {
	return asst.checker.Check(line)
}

//struct 赋值以;结尾
func (asst *Assignment) IsEndStruct(line string) bool {
	isEnd := strings.Index(line, ";")
	if isEnd < 0 {
		return false
	}
	return true
}

//解析出ModuleInfo
func (asst *Assignment) parseTableInfo(filePath string, writer *util.BufferWriter) *TableInfo {
	table := &TableInfo{}
	writer.MoveNext()
	array := util.GetLegalStrings(writer.Current())
	if len(array) >= 2 {
		table.VarName = array[len(array)-1]
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
func (asst *Assignment) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	writer = asst.FormatStruct(writer)
	table := asst.parseTableInfo(filePath, writer)
	if table == nil {
		return false
	}
	fmt.Println(writer.ToString())
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
		fields := asst.parseFields(lines)
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

//1.宏处理
func (asst *Assignment) macroHandler(writer *util.BufferWriter) *util.BufferWriter {
	writer = GetPreCompile().InOneLine(writer) //1.预编译处理成一行
	writer.Reset()
	outBuf := util.NewBufferWriter(writer.Size())
	for writer.MoveNext() {
		line := writer.Current()
		if strings.Index(line, "#if") == 0 {
			line = GetPreCompile().Parse(line[0 : len(line)-1])
		}
		outBuf.WriteString(line)
	}
	writer.Recycle()
	return outBuf
}

func (asst *Assignment) parseFields(lines []string) []string {
	fieldsLen := len(lines)
	if fieldsLen > 2 {
		lines = lines[1: fieldsLen-1]
		fields := make([]string, 0, fieldsLen-2)
		for _, val := range lines {
			fields = append(fields, val)
		}
		return fields
	}
	return nil
}

//将buffer中的struct赋值格式化成容易解析的样子
func (asst *Assignment) formatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := util.NewBufferWriter(bufWriter.Size())
	bracketCount := 0
	inBracketCount := strings.Count(bufWriter.ToString(), "[") + 1 //大括号深度
	inLittleBracket := false                                       //在小括号内部
	inQuote := false                                               //是否在双引号里面
	const BEFORE_LINE = 1
	const AFTER_LINE = 2
	const BEFORE_AFTER_LINE = 3

	for index, val := range inBuf {
		//重复的空格合并
		if val == ' ' && index < len(inBuf)-1 && inBuf[index+1] == ' ' {
			continue
		}
		//去掉结构体字段中的空格
		if val == ' ' && bracketCount > 0 {
			continue
		}
		if val == '"' {
			inQuote = !inQuote
		}
		newLine := 0
		addChar := true
		if val == '{' {
			bracketCount++
			if bracketCount == inBracketCount {
				newLine = AFTER_LINE
			}
		} else if val == '}' {
			if bracketCount == inBracketCount {
				newLine = BEFORE_AFTER_LINE
			}
			bracketCount--
		}

		if val == '(' {
			inLittleBracket = true
		} else if val == ')' {
			inLittleBracket = false
		}

		if !inQuote {
			if val == ',' && !inLittleBracket {
				if bracketCount == inBracketCount {
					newLine = AFTER_LINE
				}
				//字段后面的,去掉
				if bracketCount <= inBracketCount {
					addChar = false
				}
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
		} else {
			outBuf.WriteChar(val)
		}
	}
	return outBuf
}

//将buffer中的struct赋值格式化成容易解析的样子
func (asst *Assignment) FormatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	outBuf := asst.macroHandler(bufWriter)
	outBuf = asst.formatStruct(outBuf)
	outBuf = asst.replaceMacro(outBuf)
	outBuf = asst.formatStruct(outBuf)
	return outBuf
}

func (asst *Assignment) replaceMacro(bufWriter *util.BufferWriter) *util.BufferWriter {
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
	bufWriter.Recycle()
	return outBuf
}
