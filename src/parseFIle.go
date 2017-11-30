package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	NGX_COMMAND_T     = "ngx_command_t"
	NGX_COMMAND_T_LEN = len(NGX_COMMAND_T)
	NGX_STRING        = "ngx_string"
	NGX_STRING_LEN    = len(NGX_STRING)
)

var needCmd []string
var notNeedCmd []string

func initParse() {
	needCmd = []string{NGX_COMMAND_T, "[]"}
	notNeedCmd = []string{"(", ")"}
}

func checkNeedCmd(line string) bool {
	for _, cmd := range needCmd {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

func checkNotNeedCmd(line string) bool {
	for _, cmd := range notNeedCmd {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

//从路径中取出模块名称
func parseModuleName(fullPath string) string {
	start := strings.LastIndex(fullPath, "/")
	fullPath = fullPath[start+1:]
	end := strings.Index(fullPath, ".")
	return fullPath[:end]
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

//是否是合法字符
func isLegalChar(c uint8) bool {
	if (c >= '0' && c <= '9') || c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true
	}
	return false
}

//解析出模块定义配置信息的struct
func parseStructName(line string) string {
	index := strings.Index(line, NGX_COMMAND_T) + NGX_COMMAND_T_LEN
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

func parseCommand(cmd *Command, lines []string) {
	for index, _ := range lines {
		lines[index] = strings.Trim(lines[index], " ")
		lines[index] = strings.Trim(lines[index], ",")
	}
	cmd.TypeVal = lines[0]
	cmd.Set = lines[1]
	cmd.Conf = lines[2]
	cmd.Offset = lines[3]

	start := -1
	line := lines[4]
	for i := 0; i < len(line); i++ {
		c := line[i]
		if isLegalChar(c) {
			if start < 0 {
				start = i
			}
		} else if start >= 0 {
			cmd.Post = line[start:i]
			break
		}
	}
}

func parseFile(fullPath string) *CommandInfo {
	f, err := os.Open(fullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	var commandInfo *CommandInfo = nil
	commandLines := make([]string, 5, 5)
	commandIndex := 0
	inCommand := false
	command := Command{}
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

		if commandInfo == nil {
			hasCmd := checkNeedCmd(line)
			if hasCmd && !checkNotNeedCmd(line) {
				commandInfo = &CommandInfo{}
				commandInfo.FileName = fullPath
				commandInfo.ModuleName = parseModuleName(fullPath)
				commandInfo.StructName = parseStructName(line)
				commandInfo.CmdList = make([]Command, 0, 4)
			}
		} else {
			if isEnd := isEndOfStruct(line); isEnd {
				break
			}
			if !inCommand {
				commandName := parseCommandName(line)
				if commandName != "" {
					inCommand = true
					commandIndex = 0
					command.Name = commandName
				}
			} else {
				if strings.HasSuffix(line, ","){	//一个字段结束
					commandLines[commandIndex] += line
					commandIndex++
					if commandIndex == len(commandLines) {
						inCommand = false
						commandIndex = 0
						parseCommand(&command, commandLines)
						commandInfo.CmdList = append(commandInfo.CmdList, command)
						for key,_ := range commandLines{
							commandLines[key] = ""
						}
					}
				}else{									//处理字段换行的情况
					commandLines[commandIndex] = line
				}

			}
		}
	}
	return commandInfo
}
