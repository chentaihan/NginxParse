package logic

import (
	"fmt"
	"strings"
	"unsafe"
	"bytes"

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
		structParse: NewStructParsee(STRUCT),
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
	if def.Struct.StructName == "ngx_http_geo_ctx_t" || def.Struct.StructName == "ngx_http_geo_ctx_s"{
		i := 0
		i++
	}
	structStr := writer.ToString()
	writer = def.FormatStruct(writer)
	structStr = writer.ToString()
	def.Struct.StructString = structStr

	lines := make([]string, 0, 4)
	writer.Reset()
	for writer.MoveNext() {
		lines = append(lines, writer.Current())
	}

	for i := 1; i < len(lines)-1; i++ {
		line := def.getFieldName(lines[i])
		if line != "" {
			def.Struct.Fields = append(def.Struct.Fields, line)
		} else {
			fmt.Println("错误struct：" + structStr)
			//panic(filePath)
		}
	}
	structInfo := def.Struct
	GetStructManager().Add(&structInfo)
	def.Struct = StructInfo{}
	return true
}

func (def *Define) getFieldName(line string) string {
	//字段结构，如：int *val
	if index := strings.LastIndex(line, " "); index >= 0 {
		line = line[index:]
		return util.GetLegalString(line)
	}
	return ""
}

//将buffer中的struct赋值格式化成容易解析的样子
func (def *Define) FormatStruct(writer *util.BufferWriter) *util.BufferWriter {
	inBuf := writer.GetBuffer()
	outBuf := util.NewBufferWriter(writer.Size())
	macroDepth := 0
	macroBuf := util.NewBufferWriter(64)
	for index := 0; index < len(inBuf); index++ {
		val := inBuf[index]
		//将宏处理成一行，去掉宏中的分号
		if val == '#' {
			ifStr := string(inBuf[index+1: index+3])
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
				endif := string(inBuf[index+1: index+6])
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

	outBuf = def.formatMacro(outBuf)

	outBuf = def.formatUnion(outBuf)
	outBuf = util.MergeSequenceChar(outBuf.ToString(), '\n')
	structStr := outBuf.ToString()
	fmt.Println(structStr)
	return outBuf
}

//struct中的宏处理
func (def *Define) formatMacro(writer *util.BufferWriter) *util.BufferWriter {
	writer.Reset()
	def.unionParse.Reset()
	outBuf := util.NewBufferWriter(writer.Size())
	for writer.MoveNext() {
		line := writer.Current()
		if strings.Index(line, "#if") == 0 {
			line = GetMacroJudge().Parse(line[0: len(line)-1])
		} else {
			line = def.replaceMacro(line)
		}
		outBuf.WriteString(line)
	}
	return outBuf
}

//格式化struct中的union
func (def *Define) formatUnion(writer *util.BufferWriter) *util.BufferWriter {
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

	return outBuf
}

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
			tmpStr = tmpStr[0: len(tmpStr)-1]
			str = str[0: len(str)-1]
		}
		if util.IsLegalMacro(tmpStr) {
			hasReplace = true
			macroValue := GetMacro().GetMacroValue(str)
			if strings.HasSuffix(macroValue, ";") {
				macroValue = macroValue[0: len(macroValue)-1]
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
		line = line[index+1: end]
		return strings.Trim(line, " ")
	}
	return ""
}
