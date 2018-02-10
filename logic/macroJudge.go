package logic

import (
	"strconv"
	"strings"

	"github.com/chentaihan/NginxParse/util"
)

type JudgeFun func(first, second int64) bool

type MacroJudge struct {
	funMap   map[string]JudgeFun
	priority map[string]uint8
}

var macroJudge *MacroJudge

func GetMacroJudge() *MacroJudge {
	if macroJudge == nil {
		macroJudge = &MacroJudge{}
		macroJudge.initMap()
	}
	return macroJudge
}

func (macro *MacroJudge) initMap() {
	macro.funMap = make(map[string]JudgeFun, 18)
	macro.priority = make(map[string]uint8, 18)
	priorityValue := uint8(0)
	macro.funMap["!"] = func(first, second int64) bool { //!val
		return second == 0
	}
	macro.funMap["++"] = func(first, second int64) bool { //val++
		first++
		return first > 0
	}
	macro.funMap["--"] = func(first, second int64) bool { //val--
		first--
		return first > 0
	}
	macro.priority["!"] = priorityValue
	macro.priority["++"] = priorityValue
	macro.priority["--"] = priorityValue
	priorityValue++

	macro.funMap["+"] = func(first, second int64) bool {
		return first+second > 0
	}
	macro.funMap["-"] = func(first, second int64) bool {
		return first-second > 0
	}
	macro.priority["+"] = priorityValue
	macro.priority["-"] = priorityValue
	priorityValue++

	macro.funMap["<<"] = func(first, second int64) bool {
		return first<<uint64(second) > 0
	}
	macro.funMap[">>"] = func(first, second int64) bool {
		return first>>uint64(second) > 0
	}
	macro.priority["<<"] = priorityValue
	macro.priority[">>"] = priorityValue
	priorityValue++

	macro.funMap[">="] = func(first, second int64) bool {
		return first >= second
	}
	macro.funMap[">"] = func(first, second int64) bool {
		return first > second
	}
	macro.funMap["<="] = func(first, second int64) bool {
		return first <= second
	}
	macro.funMap["<"] = func(first, second int64) bool {
		return first >= second
	}
	macro.priority[">="] = priorityValue
	macro.priority[">"] = priorityValue
	macro.priority["<="] = priorityValue
	macro.priority["<"] = priorityValue
	priorityValue++

	macro.funMap["!="] = func(first, second int64) bool {
		return first != second
	}
	macro.funMap["=="] = func(first, second int64) bool {
		return first == second
	}
	macro.priority["!="] = priorityValue
	macro.priority["=="] = priorityValue
	priorityValue++

	macro.funMap["&"] = func(first, second int64) bool {
		return (first & second) > 0
	}
	macro.priority["&"] = priorityValue
	priorityValue++

	macro.funMap["^"] = func(first, second int64) bool {
		return first^second > 0
	}
	macro.priority["^"] = priorityValue
	priorityValue++

	macro.funMap["|"] = func(first, second int64) bool {
		return (first | second) > 0
	}
	macro.priority["|"] = priorityValue
	priorityValue++

	macro.funMap["&&"] = func(first, second int64) bool {
		return first > 0 && second > 0
	}
	macro.priority["&&"] = priorityValue
	priorityValue++

	macro.funMap["||"] = func(first, second int64) bool {
		return first > 0 || second > 0
	}
	macro.priority["||"] = priorityValue
	priorityValue++
}

func (macro *MacroJudge) Parse(line string) string {
	outBuffer := util.NewBufferWriter(len(line))
	inBuffer := util.NewBufferWriter(len(line))
	inBuffer.WriteString(line)
	inBuffer.ReplaceByte(NEWLINE_REPLACE_KEY, '\n')
	macro.parseContent(inBuffer, outBuffer)
	return outBuffer.ToString()
}

func (macro *MacroJudge) parseContent(buffer, outBuf *util.BufferWriter) {
	const (
		LOCATION_OUT    = 0
		LOCATION_INIF   = 1
		LOCATION_INELSE = 2
	)
	location := LOCATION_OUT
	ifTrue := false
	depth := 0
	isJudge := false
	lineBuffer := util.NewBufferWriter(buffer.Size())
	for buffer.MoveNext() {
		str := buffer.Current()
		str = str[util.NotEmptyIndex(str):]
		if depth <= 0 {
			location = LOCATION_OUT
		}
		if strings.HasPrefix(str, "#if") {
			depth++
			if depth <= 1 {
				location = LOCATION_INIF
			}
			if !isJudge {
				isJudge = true
				tmpLine := strings.TrimRight(str[3:], "\n")
				ifTrue = macro.judge(tmpLine)
				continue
			}

		}
		if strings.HasPrefix(str, "#else") {
			if depth <= 1 {
				location = LOCATION_INELSE
			}
			if depth <= 1 {
				if ifTrue {
					break
				} else {
					continue
				}
			}

		}

		if strings.HasPrefix(str, "#endif") {
			depth--
			if depth <= 0 {
				continue
			}
		}

		switch location {
		case LOCATION_INIF:
			if ifTrue {
				lineBuffer.WriteString(str)
			}
		case LOCATION_INELSE:
			if !ifTrue && depth >= 1 {
				lineBuffer.WriteString(str)
			}
		default:
			if depth <= 0 {
				outBuf.WriteString(str)
			}
		}

	}
	if lineBuffer.Size() > 0 {
		macro.parseContent(lineBuffer, outBuf)
	}
}

func (macro *MacroJudge) judge(line string) bool {
	//JUDGE_REPLACE: #if IS_OK && (IS_OK) || (IS_OK == 1) || (NGX_HAVE_FILE_AIO || NGX_COMPAT)
	//JUDGE_BRACE: #if 1 && (1) || (1 == 1) || (1 || 1)
	//JUDGE_CALC #if 1 && 1 || 1 || 1
	line = macro.judgeReplace(line)
	line = macro.judgeBrace(line)
	return macro.judgeCalc(line)
}

//第一步 宏替换
func (macro *MacroJudge) judgeReplace(line string) string {
	macroSlice := util.GetLegalStrings(line)
	for _, macroName := range macroSlice {
		macroValue := "0"
		if util.IsIntValue(macroName) {
			//数字不用替换
			macroValue = macroName
		} else {
			//其他的字符串一律当做宏来处理，不存在就给0
			macroValue = GetMacro().GetMacroValue(macroName)
			if macroValue == "" {
				macroValue = "0"
			}
		}

		line = strings.Replace(line, macroName, macroValue, -1)
	}
	return line
}

//第二步 计算括号中的值
func (macro *MacroJudge) judgeBrace(line string) string {
	line = util.RemoveBlank(line)
	start := -1
	i := 0
	for ; i < len(line); i++ {
		if line[i] == '(' {
			start = i
		}
		if line[i] == ')' {
			tmpLine := line[start+1: i]

			calcVal := "0"
			if macro.judgeCalc(tmpLine) {
				calcVal = "1"
			}
			line = strings.Replace(line, "("+tmpLine+")", calcVal, -1)
			//每替换完一对(),在重新查找替换
			i = 0
		}
	}
	return line
}

//将1&&0拆分成：slice{1,&&,0}
func (macro *MacroJudge) split(line string) []string {
	curVal := ""
	slice := make([]string, 0, 3)
	size := len(line)
	prevIsInt := true
	for i := 0; i < size; i++ {
		var curIsInt bool
		if line[i] >= '0' && line[i] <= '9' {
			curIsInt = true
		} else {
			curIsInt = false
		}

		if prevIsInt != curIsInt {
			prevIsInt = curIsInt
			slice = append(slice, curVal)
			curVal = ""
		}
		curVal += string(line[i])
	}
	if curVal != "" {
		slice = append(slice, curVal)
	}
	return slice
}

//第三步 计算 1||2&&3或1>=2||2<3
func (macro *MacroJudge) judgeCalc(line string) bool {
	operatorSlice := macro.split(line)

	for size := len(operatorSlice); size > 1;{
		operatorMap := make(map[int]uint8, size/2)
		for i := 1; i < size; i += 2 {
			operatorMap[i] = macro.priority[operatorSlice[i]]
		}

		maxIndex := macro.getMinPriorityIndex(operatorMap)
		operator := operatorSlice[maxIndex]
		index := maxIndex - 1
		first, _ := strconv.Atoi(operatorSlice[index])
		second, _ := strconv.Atoi(operatorSlice[maxIndex+1])
		fun := macro.funMap[operator]
		result := "0"
		if fun != nil && fun(int64(first), int64(second)) {
			result = "1"
		}
		operatorSlice[index] = result
		tmpSlice := operatorSlice[:maxIndex]
		if maxIndex+2 < size {
			tmpSlice = append(tmpSlice, operatorSlice[maxIndex+2:]...)
		}
		operatorSlice = tmpSlice
		size = len(operatorSlice)
		delete(operatorMap, maxIndex)
	}

	if result, _ := strconv.Atoi(operatorSlice[0]); result > 0 {
		return true
	}
	return false
}

//优先级值越小，优先级越高
func (macro *MacroJudge) getMinPriorityIndex(priorityMap map[int]uint8) int {
	var minKey int = -1
	for key, value := range priorityMap {
		if minKey == -1 {
			minKey = key
		}
		if value < priorityMap[minKey] {
			minKey = key
		}
	}
	return minKey
}
