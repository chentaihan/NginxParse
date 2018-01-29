package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type FileParse struct {
	Structs []IParse
}

func (fileParse *FileParse) Register(stParse IParse) {
	fileParse.Structs = append(fileParse.Structs, stParse)
}

func (fileParse *FileParse) isStartStruct(line string) IParse {
	for _, stt := range fileParse.Structs {
		if stt.IsStartStruct(line) {
			return stt
		}
	}
	return nil
}

//解析主流程
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
	depth := 0
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
		//过滤注释
		line = filterNote(line, &inNote)
		//合并空格
		line = mergeSpace(line)

		if curStruct == nil {
			if curStruct = fileParse.isStartStruct(line); curStruct != nil {
				buffer = &BufferWriter{}
				isOK = true
			}
		}

		if curStruct != nil {
			line = strings.Trim(line, " ")
			if isEmptyLine(line) {
				continue
			}

			buffer.WriteString(line)
			depth += getDepth(line)
			if depth == 0 && curStruct.IsEndStruct(line) {
				curStruct.ParseStruct(fullPath, buffer)
				curStruct = nil
				//str := buffer.ToString()
				//fmt.Println(str)
				buffer = &BufferWriter{}
			}
		}
	}

	return isOK
}

/******************************过滤注释 start***************************************/

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

/******************************过滤注释 end***************************************/

func getDepth(line string) int {
	count := 0

	keys := map[int]string{
		1:  "{",	//存在一个{，加1
		-1: "}",	//存在一个}，减1
	}
	for cnt, key := range keys {
		tmpLine := line
		for {
			index := strings.Index(tmpLine, key)
			if index >= 0 {
				count += cnt
				if index == len(tmpLine) -1 {
					break
				}
				tmpLine = tmpLine[index+1:]
			} else {
				break
			}
		}
	}

	return count
}
