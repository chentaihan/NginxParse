package main

import "strings"

//从路径中取出模块名称
func parseModuleName(fullPath string) string {
	start := strings.LastIndex(fullPath, "/")
	fullPath = fullPath[start+1:]
	end := strings.Index(fullPath, ".")
	return fullPath[:end]
}

//是否是合法字符
func isLegalChar(c uint8) bool {
	if (c >= '0' && c <= '9') || c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

func parseCommandName(line string) string {
	index := strings.Index(line, NGX_STRING)
	if index < 0 {
		return ""
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
	start := -1
	for ; index < len(line); index++ {
		c := line[index]
		if isLegalChar(c) {
			if start < 0 {
				start = index
			}
			continue
		}
		if start > 0 {
			return line[start:index]
		}
	}
	return ""
}
