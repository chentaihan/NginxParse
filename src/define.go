package main

import "strings"

type Define struct {
	StructInfo *StructPraseInfo
	checker    Checker
}

func NewParseDefine(sutInfo *StructPraseInfo) *Define {
	varMgr := &Define{
		StructInfo: sutInfo,
	}
	sturctName := sutInfo.Include[0]
	sturctName = sturctName[0:len(sturctName)-1] + "s"
	varMgr.checker = Checker{
		Include:   []string{sturctName, "struct", "{"},
		IncludeNo: []string{"(", ")"},
	}
	return varMgr
}

//判断是不是有效结构体
func (varMgr *Define) Check(line string) bool {
	return varMgr.checker.Check(line)
}

//解析出结构体内容
func (mgr *Define) ParseStruct(filePath string, writer *BufferWriter) bool {
	lines := make([]string, 0, 4)
	for writer.MoveNext() {
		lines = append(lines, writer.Current())
	}
	mgr.parseFieldName(lines)
	//如果没有指定字段，就显示全部字段
	if len(mgr.StructInfo.Fields) == 0 {
		for i := 0; i < len(mgr.StructInfo.FieldNames); i++ {
			mgr.StructInfo.Fields = append(mgr.StructInfo.Fields, i)
		}
	}
	return true
}

func (mgr *Define) parseFieldName(lines []string) {
	for i := 1; i < len(lines)-1; i++ {
		line := mgr.getFieldName(lines[i])
		mgr.StructInfo.FieldNames = append(mgr.StructInfo.FieldNames, line)
	}
}

func (mgr *Define) getFieldName(line string) string {
	index := strings.Index(line, "(")
	if index < 0 {
		index = strings.LastIndex(line, " ")
	}
	line = line[index:]
	return getLegalString(line)
}

//将buffer中的struct赋值格式化成容易解析的样子
func (mgr *Define) FormatStruct(bufWriter *BufferWriter) *BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := NewBufferWriter(bufWriter.Size)
	for index, val := range inBuf {
		//重复的空格合并
		if val == ' ' && index < len(inBuf)-1 && inBuf[index+1] == ' ' {
			continue
		}
		if val != ';' {
			outBuf.WriteChar(val)
		}
		if val == ';' || val == '{' {
			outBuf.WriteChar('\n')
		}
	}
	return outBuf
}

//struct 定义以}结尾
func (mgr *Define) IsEndStruct(line string) bool {
	isEnd := strings.Index(line, "}")
	if isEnd < 0 {
		return false
	}
	return true
}