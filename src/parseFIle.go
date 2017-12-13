package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//struct以;结尾
func isEndOfStruct(line string) bool {
	isEnd := strings.Index(line, ";")
	if isEnd < 0 {
		return false
	}
	index := strings.Index(line, "//")
	if index >= 0 && index < isEnd {
		return false
	}
	return true
}

func isEmptyLine(line string) bool {
	for i := len(line) - 1; i >= 0; i++ {
		if line[i] != ' ' {
			return false
		}
	}
	return true
}

type StructParse struct {
	Parse      IParse //解析接口
	structType int    //解析结构类型
}

type FileParse struct {
	Structs   []*StructParse
	curStruct *StructParse
}

type StructContent struct {
	buffer     *BufferWriter
	filePath   string
	structType int //解析结构类型
}

func (fileParse *FileParse) Register(stParse *StructParse) {
	fileParse.Structs = append(fileParse.Structs, stParse)
}

func (fileParse *FileParse) Check(line string) bool {
	for _, stt := range fileParse.Structs {
		if stt.Parse.Check(line) {
			fileParse.curStruct = stt
			return true
		}
	}
	return false
}

func (fileParse *FileParse) Parse(fullPath string) bool {
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	inStruct := false
	isOK := false
	inNote := false
	var buffer *BufferWriter
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		line = strings.TrimRight(line, "\n")
		line = strings.Trim(line, " ")
		if isEmptyLine(line) {
			continue
		}

		if !inStruct {
			if fileParse.Check(line) {
				buffer = &BufferWriter{}
				buffer.WriteString(line)
				inStruct = true
				isOK = true
				continue
			}
		}

		if inStruct {
			line = filterNote(line, &inNote)
			buffer.WriteString(line)

			if isEndOfStruct(line) {
				inStruct = false
				buffer = formatSturct(buffer)
				fmt.Println(buffer.ToString())
				fileParse.curStruct.Parse.ParseStruct(fullPath, buffer)
			}
		}
	}

	return isOK
}

func filterNote(line string, inNote *bool) string {
	if *inNote {
		return filterInNote(line, inNote)
	}
	return filterNotInNote(line, inNote)
}

func filterInNote(line string, inNote *bool) string {
	start := strings.Index(line, "*/")
	if start >= 0 { // 结束*/
		*inNote = false
		return filterNotInNote(line[start+2:], inNote)
	}
	return ""
}

func filterNotInNote(line string, inNote *bool) string {
	if start := strings.Index(line, "//"); start >= 0 {
		line = line[0:start]
	}

	for {
		start := strings.Index(line, "/*")
		if start >= 0 { // 进入/*
			*inNote = true
			if end := strings.Index(line, "*/"); end >= 0 { //结束*/
				line = line[0:start] + line[end+2:]
				*inNote = true
			} else { //没有*/
				line = line[0:start]
				break
			}
		} else {
			break
		}
	}

	return line
}

//将buffer中的struct格式化成容易解析的样子
func formatSturct(bufWriter *BufferWriter) *BufferWriter {
	inBuf := bufWriter.GetBuffer()
	outBuf := NewBufferWriter(bufWriter.Size + 64)
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
		if newLine&BEFORE_LINE > 0 {
			outBuf.WriteChar('\n')
		}
		if addChar {
			outBuf.WriteChar(val)
		}
		if newLine&AFTER_LINE > 0 {
			outBuf.WriteChar('\n')
		}

	}
	return outBuf
}
