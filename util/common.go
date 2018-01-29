package util

import (
	"strings"
)

//从路径中取出模块名称
func ParseModuleName(fullPath string) string {
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

//是否是合法变量名
func IsLegalString(line string) bool {
	if line == "" {
		return false
	}
	for _, c := range line {
		if !isLegalChar(c) {
			return false
		}
	}
	return true
}

//是否是合法宏名称(0-9_A-Z)
func isLegalMacro(line string) bool {
	if line == "" {
		return false
	}
	for _, c := range line {
		isLegal := (c >= '0' && c <= '9') || c == '_' || (c >= 'A' && c <= 'Z')
		if !isLegal {
			return false
		}
	}
	return true
}

//从字符串中取出第一个合法变量名
func GetLegalString(line string) string {
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

//从字符串中取出所有合法变量名
func GetLegalStrings(line string) []string {
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

func ParseName(line string) string {
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
func ParseStructName(line string, structType string) string {
	index := strings.Index(line, structType) + len(structType)
	line = line[index:]
	return GetLegalString(line)
}

//合并重复空格
func MergeSpace(line string) string {
	if line == "" {
		return line
	}
	inBuf := []byte(line)
	outBuf := NewBufferWriter(len(inBuf))
	for index, val := range inBuf {
		if val == ' ' && index < len(inBuf)-1 && inBuf[index+1] == ' ' {
			continue
		}
		outBuf.WriteChar(val)
	}
	return outBuf.ToString()
}

//判断字符串是否是空行
func IsEmptyLine(line string) bool {
	if line == "" {
		return true
	}
	for i := len(line) - 1; i >= 0; i++ {
		if line[i] != ' ' {
			return false
		}
	}
	return true
}
