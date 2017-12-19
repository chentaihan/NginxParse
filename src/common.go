package main

import (
	"strings"
)

//从路径中取出模块名称
func parseModuleName(fullPath string) string {
	start := strings.LastIndex(fullPath, "/")
	fullPath = fullPath[start+1:]
	end := strings.Index(fullPath, ".")
	return fullPath[:end]
}

//是否是合法字符
func isLegalChar(c rune) bool {
	if (c >= '0' && c <= '9') || c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

func getLegalString(line string) string {
	isStart := false
	newLine := make([]rune, 0, len(line))
	for _, c := range line {
		if isLegalChar(c) {
			if !isStart {
				isStart = true
			}
			newLine = append(newLine, c)
			continue
		}
		if isStart {
			break
		}
	}
	return string(newLine)
}

func getLegalStrings(line string) []string {
	isStart := false
	ret := make([]string, 0)
	newLine := make([]rune, 0, len(line))
	for _, c := range line {
		if isLegalChar(c) {
			if !isStart {
				isStart = true
			}
			newLine = append(newLine, c)
			continue
		}
		if isStart {
			isStart = false
			ret = append(ret, string(newLine))
			newLine = newLine[0:0]
		}
	}
	if isStart {
		ret = append(ret, string(newLine))
	}
	return ret
}

func parseName(line string) string {
	index := strings.Index(line, NGX_STRING)
	if index < 0 {
		return line
	}
	line = line[index+NGX_STRING_LEN:]
	start := strings.Index(line, "\"")
	if start < 0 {
		return ""
	}
	line = line[start+1:]
	end := strings.Index(line, "\"")
	if end < 0 {
		return ""
	}
	return line[0:end]
}

//解析出模块定义配置信息的struct
func parseStructName(line string, structType string) string {
	index := strings.Index(line, structType) + len(structType)
	line = line[index:]
	return getLegalString(line)
}