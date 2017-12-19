package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"fmt"
)

func isEmptyLine(line string) bool {
	for i := len(line) - 1; i >= 0; i++ {
		if line[i] != ' ' {
			return false
		}
	}
	return true
}

type FileParse struct {
	Structs   []IParse
}

func (fileParse *FileParse) Register(stParse IParse) {
	fileParse.Structs = append(fileParse.Structs, stParse)
}

func (fileParse *FileParse) Check(line string) IParse {
	for _, stt := range fileParse.Structs {
		if stt.Check(line) {
			return stt
		}
	}
	return nil
}

func (fileParse *FileParse) Parse(fullPath string) bool {
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	isOK := false
	inNote := false
	var buffer *BufferWriter
	var curStruct IParse = nil
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

		if curStruct == nil {
			if curStruct = fileParse.Check(line); curStruct != nil {
				buffer = &BufferWriter{}
				isOK = true
			}
		}
		if curStruct != nil {
			line = filterNote(line, &inNote)
			line = strings.Trim(line, " ")
			buffer.WriteString(line)

			if curStruct.IsEndStruct(line) {
				buffer = curStruct.FormatStruct(buffer)
				fmt.Println(buffer.ToString())
				curStruct.ParseStruct(fullPath, buffer)
				curStruct = nil
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
			if end := strings.Index(line, "*/"); end >= 0 { //结束*/
				line = line[0:start] + line[end+2:]
				*inNote = false
			} else { //没有*/
				line = line[0:start]
				*inNote = true
				break
			}
		}
		break
	}

	return line
}
