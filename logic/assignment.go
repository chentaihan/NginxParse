package logic

/*
结构体实例化
 */

import (
	"strings"
	"fmt"

	"github.com/chentaihan/NginxParse/util"
)

type Assignment struct {
	Tables     []*TableInfo
	checker    Checker
	StructInfo *util.StructPraseInfo
}

func NewAssignment(sutInfo *util.StructPraseInfo) *Assignment {
	varMgr := &Assignment{
		Tables:     make([]*TableInfo, 0),
		StructInfo: sutInfo,
		checker: Checker{
			Include:   sutInfo.Include,
			IncludeNo: sutInfo.IncludeNo,
		},
	}
	return varMgr
}

//判断是不是有效结构体
func (varMgr *Assignment) IsStartStruct(line string) bool {
	return varMgr.checker.Check(line)
}

//解析出ModuleInfo
func (mgr *Assignment) parseTableInfo(filePath string, writer *util.BufferWriter) *TableInfo {
	writer.MoveNext()
	line := writer.Current()
	table := &TableInfo{}
	table.FileName = filePath
	table.ModuleName = util.ParseModuleName(filePath)
	table.StructName = util.ParseStructName(line, mgr.StructInfo.StructName)
	table.StructString = writer.ToString()

	titleLen := len(mgr.StructInfo.Fields)
	title := make([]string, 0, titleLen)
	for _, index := range mgr.StructInfo.Fields {
		if index >= len(mgr.StructInfo.FieldNames) {
			fmt.Println(index, len(mgr.StructInfo.FieldNames))
		}else{
			title = append(title, mgr.StructInfo.FieldNames[index])
		}

	}
	table.Title = title
	table.Content = make([][]string, 0, 4)
	return table
}

//解析出结构体内容
func (mgr *Assignment) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	writer = mgr.FormatStruct(writer)
	fmt.Println(writer.ToString())
	table := mgr.parseTableInfo(filePath, writer)
	fieldCount := len(mgr.StructInfo.FieldNames) + 2
	lines := make([]string, fieldCount, fieldCount)
	for !writer.IsEnd() {
		i := 0
		for ; i < fieldCount && writer.MoveNext(); i++ {
			lines[i] = writer.Current()
		}
		if i < fieldCount {
			break
		}
		content := mgr.parseContent(lines)
		if content != nil {
			table.Content = append(table.Content, content)
		}
	}

	mgr.Tables = append(mgr.Tables, table)
	return true
}

func (mgr *Assignment) parseContent(lines []string) []string {
	fieldsLen := len(mgr.StructInfo.Fields)
	if fieldsLen > len(lines) {
		return nil
	}
	content := make([]string, 0, fieldsLen)
	lines = lines[1:]
	for _, val := range mgr.StructInfo.Fields {
		line := util.ParseName(lines[val])
		content = append(content, line)
	}
	return content
}

//将buffer中的struct赋值格式化成容易解析的样子
func (mgr *Assignment) formatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := util.NewBufferWriter(bufWriter.Size + 64)
	bracketCount := 0
	inBracketCount := 1      //大括号深度
	inLittleBracket := false //在小括号内部

	const BEFORE_LINE = 1
	const AFTER_LINE = 2
	const BEFORE_AFTER_LINE = 3

	for _, val := range inBuf {
		if val == '[' {
			inBracketCount++
			break
		}
	}

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

		if val == ',' {
			if !inLittleBracket {
				if bracketCount == inBracketCount {
					newLine = AFTER_LINE
				}
				//}后面的,去掉
				//if inBuf[index-1] == '}' {
				addChar = false
				//}
			}
		}
		if newLine&BEFORE_LINE > 0 && inBuf[index+1] != '\n' {
			outBuf.WriteChar('\n')
		}
		if addChar {
			outBuf.WriteChar(val)
		}
		if newLine&AFTER_LINE > 0 && inBuf[index+1] != '\n' {
			outBuf.WriteChar('\n')
		}
	}
	return outBuf
}

//将buffer中的struct赋值格式化成容易解析的样子
func (mgr *Assignment) FormatStruct(bufWriter *util.BufferWriter) *util.BufferWriter {
	outBuf := mgr.formatStruct(bufWriter)
	if mgr.StructInfo.MacroExpand {
		outBuf = mgr.replaceMacro(outBuf)
		return mgr.formatStruct(outBuf)
	}
	return outBuf
}

func (mgr *Assignment) replaceMacro(bufWriter *util.BufferWriter) *util.BufferWriter {
	bufWriter.Reset()
	outBuf := util.NewBufferWriter(bufWriter.Size + 64)
	macro := GetMacro()
	for bufWriter.MoveNext() {
		line := bufWriter.Current()
		macroName := strings.Trim(line, "\n")
		value := macro.GetMacroValue(macroName)
		if value != "" {
			outBuf.WriteString(value)
			outBuf.WriteChar('\n')
		} else {
			outBuf.WriteString(line)
		}
	}
	return outBuf
}

//struct 赋值以;结尾
func (mgr *Assignment) IsEndStruct(line string) bool {
	isEnd := strings.Index(line, ";")
	if isEnd < 0 {
		return false
	}
	return true
}
