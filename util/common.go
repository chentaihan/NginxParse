package util

import (
	"strings"
	"unsafe"
)

//从路径中取出模块名称
func ParseModuleName(fullPath string) string {
	start := strings.LastIndex(fullPath, "/")
	fullPath = fullPath[start+1:]
	end := strings.Index(fullPath, ".")
	return fullPath[:end]
}

//是否是合法字符
func isLegalChar(c byte) bool {
	if (c >= '0' && c <= '9') || c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

//判断是否是int字符串
func IsIntValue(line string) bool{
	for _, c := range line {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

//是否是合法变量名
func IsLegalString(line string) bool {
	if line == "" {
		return false
	}
	for i := len(line) - 1; i >= 0; i-- {
		if !isLegalChar(line[i]) {
			return false
		}
	}
	return true
}

//是否是合法宏名称(0-9_A-Z)
func IsLegalMacro(line string) bool {
	if line == "" {
		return false
	}
	hasUpperCase := false
	for _, c := range line {
		isLegal := (c >= '0' && c <= '9') || c == '_' || (c >= 'A' && c <= 'Z')
		if !isLegal {
			return false
		}
		if c >= 'A' && c <= 'Z' {
			hasUpperCase = true
		}
	}
	return hasUpperCase
}

//从字符串中取出第一个合法变量名
func GetLegalString(line string) string {
	isStart := false
	newLine := make([]byte, 0, len(line))
	for i := 0; i < len(line); i++ {
		if isLegalChar(line[i]) {
			if !isStart {
				isStart = true
			}
			newLine = append(newLine, line[i])
			continue
		}
		if isStart {
			break
		}
	}
	return *(*string)(unsafe.Pointer(&newLine))
}

//从字符串中取出所有合法变量名
func GetLegalStrings(line string) []string {
	isStart := false
	ret := make([]string, 0)
	newLine := make([]byte, 0, len(line))
	for i := 0; i < len(line); i++ {
		if isLegalChar(line[i]) {
			if !isStart {
				isStart = true
			}
			newLine = append(newLine, line[i])
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

//合并指定的连续字符
func MergeSequenceChar(line string, c byte) *BufferWriter {
	inBuf := []byte(line)
	outBuf := NewBufferWriter(len(inBuf))
	for index, val := range inBuf {
		if val == c && index < len(inBuf)-1 && inBuf[index+1] == c {
			continue
		}
		outBuf.WriteChar(val)
	}
	return outBuf
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

func BytesToString(slice [][]byte) string {
	totalLen := 0
	for _, s := range slice {
		totalLen += len(s)
	}
	ret := make([]byte, 0, totalLen)
	for _, s := range slice {
		ret = append(ret, s...)
	}
	return *(*string)(unsafe.Pointer(&ret))
}

func NotEmptyIndex(line string) int {
	for i := 0; i < len(line); i++ {
		if line[i] != ' ' {
			return i
		}
	}
	return 0
}

func GetInt(line string) (int64, bool){
	retVal := int64(0)
	isOk := false
	right := true
	i := 0
	if line[0] == '+' {
		i++
	}else if line[0] == '-' {
		i++
		right = false
	}
	for ; i < len(line); i++ {
		val := line[i]
		if val >= '0' && val <= '9' {
			isOk = true
			retVal = retVal*10 + int64(val- '0')
		}else{
			break
		}
	}
	if !right {
		return -retVal, isOk
	}
	return retVal, isOk
}

func RemoveBlank(line string) string{
	size := len(line)
	buffer := make([]byte,size, size)
	index := 0
	for i := 0; i < size; i++{
		if line[i] != ' '{
			buffer[index] = line[i]
			index++
		}
	}
	buffer = buffer[:index]
	return *(*string)(unsafe.Pointer(&buffer))
}

