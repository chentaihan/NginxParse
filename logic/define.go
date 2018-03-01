package logic

/**
struct定义解析
*/

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"

	"github.com/chentaihan/NginxParse/util"
)

const (
	NEWLINE_REPLACE_KEY = '$'
)

type Define struct {
	Struct      StructInfo
	structParse *StructParse
	unionParse  *Union
}

func NewParseDefine() *Define {
	return &Define{
		structParse: NewStructParsee(util.STRUCT),
		unionParse:  NewUnion(),
	}
}

//判断是不是结构体头部
func (def *Define) IsStartStruct(line string) bool {
	isHead := def.structParse.IsHead(line)
	if isHead {
		def.Struct.StructName = def.structParse.StructName
	}
	return isHead
}

//struct结尾
func (def *Define) IsEndStruct(line string) bool {
	isTail := def.structParse.IsTail(line)
	if isTail {
		def.Struct.Rename = def.structParse.Rename
		if def.Struct.StructName == "" {
			def.Struct.StructName = def.Struct.Rename
		}
	}
	return isTail
}

//解析出结构体内容
func (def *Define) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	def.Struct.FileName = filePath
	def.Struct.ModuleName = util.ParseModuleName(filePath)
	structStr := writer.ToString()
	writer = def.FormatStruct(writer)
	structStr = writer.ToString()
	def.Struct.StructString = structStr
	if def.Struct.StructName == "ngx_stream_module_t" {
		fmt.Println(structStr)
	}
	lines := make([]string, 0, 4)
	writer.Reset()
	for writer.MoveNext() {
		lines = append(lines, writer.Current())
	}
	util.Println("-----------------------------------", def.Struct.StructName)
	for i := 1; i < len(lines)-1; i++ {
		fieldName := def.getFieldName(lines[i])
		if fieldName != "" {
			def.Struct.Fields = append(def.Struct.Fields, fieldName)
		} else {
			util.Println("错误struct：" + structStr)
			panic(filePath)
		}
	}
	structInfo := def.Struct
	GetDefines().Add(structInfo.StructName, &structInfo)
	def.Struct = StructInfo{}
	def.structParse.Reset()
	return true
}

/**
获取字段名称
字段结构，如：int *val
*/
func (def *Define) getFieldName(line string) string {
	if index := strings.Index(line, " "); index >= 0 {
		line = line[index:]
		return util.GetLegalString(line)
	}
	return ""
}

//将buffer中的struct赋值格式化成容易解析的样子
func (def *Define) FormatStruct(writer *util.BufferWriter) *util.BufferWriter {
	outBuf := def.precompileHandler(writer) //1.预编译处理成一行
	outBuf = def.macroHandler(outBuf)       //2.宏处理
	outBuf = def.unionHandler(outBuf)       //3.联合体处理
	outBuf = def.multFieldHandler(outBuf)   //4.一行多个字段转换成一行一字段
	outBuf = def.partFieldHandler(outBuf)   //5.多行一个字段转换成一行一字段
	return outBuf
}

//1.预编译处理成一行
func (def *Define) precompileHandler(writer *util.BufferWriter) *util.BufferWriter {
	inBuf := writer.GetBuffer()
	outBuf := util.NewBufferWriter(writer.Size())
	macroDepth := 0
	macroBuf := util.NewBufferWriter(64)
	for index := 0; index < len(inBuf); index++ {
		val := inBuf[index]
		if val == '#' {
			ifStr := string(inBuf[index+1 : index+3])
			if ifStr == "if" {
				macroDepth++
			}
		}

		if macroDepth > 0 {
			//\n用$替换
			if val != '\n' {
				macroBuf.WriteChar(val)
			} else {
				macroBuf.WriteChar(NEWLINE_REPLACE_KEY)
			}
			if val == '#' {
				endif := string(inBuf[index+1 : index+6])
				if endif == "endif" {
					macroDepth--
					if macroDepth == 0 {
						outBuf.Write(macroBuf.GetBuffer())
						macroBuf.Clear()
						outBuf.WriteString(endif)
						index += 5
					}
				}
			}
			continue
		}
		outBuf.WriteChar(val)
	}
	return outBuf
}

//2.宏处理
func (def *Define) macroHandler(writer *util.BufferWriter) *util.BufferWriter {
	writer.Reset()
	def.unionParse.Reset()
	outBuf := util.NewBufferWriter(writer.Size())
	for writer.MoveNext() {
		line := writer.Current()
		if strings.Index(line, "#if") == 0 {
			line = GetPreCompile().Parse(line[0 : len(line)-1])
		} else {
			line = def.replaceMacro(line)
		}
		outBuf.WriteString(line)
	}
	return outBuf
}

//3.联合体处理
func (def *Define) unionHandler(writer *util.BufferWriter) *util.BufferWriter {
	unionLineIndex := -1
	writer.Reset()
	def.unionParse.Reset()
	outBuf := util.NewBufferWriter(writer.Size())
	for writer.MoveNext() {
		line := writer.Current()
		if unionLineIndex == -1 {
			if def.unionParse.IsHead(line) {
				unionLineIndex = 0
			}
		}
		if unionLineIndex > -1 {
			unionLineIndex++
			isTail := def.unionParse.IsTail(line)
			line = strings.Replace(line, "\n", "", -1)
			def.unionParse.AddLine(line)
			if isTail {
				unionLineIndex = -1
				//完整的union结构
				str := def.unionParse.GetUnionField()
				outBuf.WriteString(str)
				outBuf.WriteString("\n")
				def.unionParse.Reset()
			}
		} else {
			outBuf.WriteString(line)
		}
	}
	//删除多余的空行
	return util.MergeSequenceChar(outBuf.ToString(), '\n')
}

/**
宏替换
*/
func (def *Define) replaceMacro(line string) string {
	byteSlice := bytes.SplitN([]byte(line), []byte{')'}, -1)
	hasReplace := false
	for i, slice := range byteSlice {
		str := *(*string)(unsafe.Pointer(&slice))
		index := strings.Index(str, "(")
		var tmpStr string
		if index > 0 {
			str += ")"
			tmpStr = str[0:index]
		} else {
			tmpStr = str
		}
		isEndwithN := false
		if str[len(str)-1] == '\n' {
			isEndwithN = true
			tmpStr = tmpStr[0 : len(tmpStr)-1]
			str = str[0 : len(str)-1]
		}
		if util.IsLegalMacro(tmpStr) {
			hasReplace = true
			macroValue := GetMacros().GetMacroValue(str)
			if strings.HasSuffix(macroValue, ";") {
				macroValue = macroValue[0 : len(macroValue)-1]
			}
			if isEndwithN {
				macroValue += "\n"
			}
			byteSlice[i] = []byte(macroValue)
		}
	}
	if hasReplace {
		return util.BytesToString(byteSlice)
	}
	return line
}

//获取宏名称
func (def *Define) getMacroName(line string) string {
	strLen := len(line)
	if strLen <= 10 {
		return ""
	}
	line = line[3:]
	return util.GetLegalString(line)
}

//获取宏包含的字段定义信息
func (def *Define) getMacroField(line string) string {
	strLen := len(line)
	if strLen <= 10 {
		return ""
	}
	line = line[3:]
	index := strings.Index(line, ")")
	end := strings.Index(line, "#")
	if index > 0 && end > 0 && end > index {
		line = line[index+1 : end]
		return strings.Trim(line, " ")
	}
	return ""
}

//4.一行多个字段转换成一行一字段
func (def *Define) multFieldHandler(writer *util.BufferWriter) *util.BufferWriter {
	writer.Reset()
	outBuf := util.NewBufferWriter(writer.Size())
	for writer.MoveNext() {
		line := writer.Current()
		if strings.Index(line, ",") > 0 {
			outBuf.WriteString(def.oneLineToMultFields(line))
		} else {
			outBuf.WriteString(line)
		}
	}
	return outBuf
}

/*
将一行包含多个字段拆分成多行
如将int a,b,c;转成：
int a;
int b;
int c;
*/
func (def *Define) oneLineToMultFields(line string) string {
	inBrackets := false
	typeStr := ""
	size := len(line)
	hasComma := false //是否存在逗号
	outBuf := util.NewBufferWriter(size)
	start := 0

	writeString := func(str1, str2 string) {
		outBuf.WriteString(str1)
		outBuf.WriteString(" ")
		outBuf.WriteString(str2)
	}

	if index := strings.Index(line, " "); index > 0 {
		typeStr = line[:index]
		index++
		start = index
		for ; index < size; index++ {
			if line[index] == '(' {
				inBrackets = true
			} else if line[index] == ')' {
				inBrackets = false
			}
			if !inBrackets && line[index] == ',' {
				writeString(typeStr, line[start:index] + ";\n")
				start = index + 1
				hasComma = true
			}
		}
	}
	if hasComma {
		writeString(typeStr, line[start:])
		return outBuf.ToString()
	}
	return line
}

/**
5.多行一个字段转换成一行一字段
将多行构成的一个字段转成一行，即一行一个字段
struct每个字段都是以";"结尾，如果\n前面不是";"说明这个\n就是多余的
*/
func (def *Define) partFieldHandler(writer *util.BufferWriter) *util.BufferWriter {
	inBuf := writer.GetBuffer()
	specialChars := ";{}"
	for index := 1; index < len(inBuf); index++ {
		if inBuf[index] == '\n' && !util.ContainsByte(specialChars, inBuf[index-1]) {
			writer.RemoveByte(index)
		}
	}
	return writer
}
