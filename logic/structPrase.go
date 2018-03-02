package logic

import (
	"strings"
	"github.com/chentaihan/NginxParse/util"
)

/**
解析struct/union
*/

type StructParse struct {
	headFormats [][]string
	tailFormats [][]string
	hasTypedef  int
	StructName  string
	Rename      string
}


func NewStructParsee(structType string) *StructParse {
	sttParse := &StructParse{
		headFormats: [][]string{
			{structType, "", "{"},
			{"typedef", structType, "{"},
			{"typedef", structType, "", "{"},
		},
		tailFormats: [][]string{
			{"}"},
			{"}", ""},
		},
		hasTypedef: 0,
	}

	return sttParse
}

//判断是不是struct/union头部
func (stt *StructParse) IsHead(line string) bool {
	line = util.MergeSequenceChar(line, ' ')
	token := strings.SplitN(line, " ", 4)
	tokenLen := len(token)
	if tokenLen < 3 {
		return false
	}

	//struct头部就三种结构
	//struct name {
	//typedef struct {
	//typedef struct name {

	nameIndex := -1
	for _, format := range stt.headFormats {
		if tokenLen != len(format) {
			continue
		}
		i := 0
		for ; i < len(token); i++ {
			if format[i] == "" {
				if util.IsLegalString(token[i]) {
					nameIndex = i
				} else {
					break
				}
			} else {
				if token[i] != format[i] {
					break
				}
			}

		}
		if i == tokenLen {
			if token[0] == "typedef" {
				stt.hasTypedef = 1
			} else {
				stt.hasTypedef = 0
			}
			if nameIndex >= 0 {
				stt.StructName = token[nameIndex]
			}
			return true
		}
	}

	return false
}

//struct/union结尾
func (stt *StructParse) IsTail(line string) bool {
	if line[len(line)-1] != ';' {
		return false
	}
	line = line[0 : len(line)-1]
	line = util.MergeSequenceChar(line, ' ')
	token := strings.SplitN(line, " ", 2)

	//结构体尾部就两种结构
	//}
	//} name

	formats := stt.tailFormats[stt.hasTypedef]
	if len(token) == len(formats) {
		if token[0] != formats[0] {
			return false
		}
		if len(token) > 1 {
			if util.IsLegalString(token[1]) {
				stt.Rename = token[1]
			}
		}
		return true
	}
	return false
}

func (stt *StructParse) Reset() {
	stt.hasTypedef = 0
	stt.StructName = ""
	stt.Rename = ""
}
