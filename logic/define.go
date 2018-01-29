package logic

import (
	"fmt"
	"strings"

	"github.com/chentaihan/NginxParse/util"
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
	if def.Struct.StructName == "ngx_http_upstream_server_t" {
		i := 0
		i++
	}
	writer = def.FormatStruct(writer)
	structStr := writer.ToString()
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
	fmt.Println(writer.ToString())
	outBuf := util.NewBufferWriter(writer.Size)
	inMacro := false
	macroBuf := util.NewBufferWriter(64)
	for index := 0; index < len(inBuf); index++ {
		val := inBuf[index]
		//将宏处理成一行，去掉宏中的分号
		if !inMacro {
			if val == '#' {
				ifStr := string(inBuf[index+1 : index+3])
				if ifStr == "if" {
					inMacro = true
				}
			}
		}
		if inMacro {
			if val != ';' {
				macroBuf.WriteChar(val)
			}

			if val == '#' {
				endif := string(inBuf[index+1 : index+6])
				if endif == "endif" {
					inMacro = false
					outBuf.Write(macroBuf.GetBuffer())
					outBuf.WriteString(endif)
					outBuf.WriteChar('\n')
					index += 5
				}
			}
			continue
		}

		if val != ';' {
			outBuf.WriteChar(val)
		}
		if val == ';' || val == '{' {
			outBuf.WriteChar('\n')
		}

	}
	outBuf = def.formatUnion(outBuf)
	return def.formatMacro(outBuf)
}

func (def *Define) macroReplace(writer *util.BufferWriter) *util.BufferWriter{
	return writer
}


//格式化struct中的union
func (def *Define) formatUnion(writer *util.BufferWriter) *util.BufferWriter {
	unionLineIndex := -1
	writer.Reset()
	def.unionParse.Reset()
	outBuf := util.NewBufferWriter(writer.Size)
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
			def.unionParse.AddLine(line)
			//union每个字段后面加;
			if unionLineIndex > 1 && !isTail {
				def.unionParse.AddLine(";")
			}
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

//struct中的宏处理
func (def *Define) formatMacro(writer *util.BufferWriter) *util.BufferWriter {
	writer.Reset()
	def.unionParse.Reset()
	outBuf := util.NewBufferWriter(writer.Size)
	for writer.MoveNext() {
		exist := true
		line := writer.Current()
		if strings.Index(line, "#if") == 0 {
			macroName := def.getMacroName(line)
			//宏不存在就去掉这个字段
			if macroName != "" {
				if !GetMacro().Exist(macroName) {
					exist = false
				}else{
					outBuf.WriteString(def.getMacroField(line))
				}
			}
		}
		if exist {
			outBuf.WriteString(line)
		}

	}
	return outBuf
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
	if index > 0 && end > 0 && end > index{
		line = line[index+1:end]
		return strings.Trim(line, " ")
	}
	return ""
}