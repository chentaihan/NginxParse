package logic

/**
宏处理
*/

import (
	"strings"

	"github.com/chentaihan/NginxParse/util"
)

type Typedef struct {
	StructName   string
	StructRename string
}

func NewTypedef() *Typedef {
	return &Typedef{}
}

//判断是不是struct重定义
//typedef struct ngx_stream_session_s  ngx_stream_session_t;
//typedef ngx_stream_session_s  ngx_stream_session_t;
func (typedef *Typedef) IsStartStruct(line string) bool {
	isStartOK := strings.HasPrefix(line, util.NGX_TYPEDEF)
	isEndOK := strings.HasSuffix(line, ";")
	index := -1
	if isStartOK && isEndOK {
		strArr := strings.Split(line, " ")
		if len(strArr) == 4 && strArr[1] == util.STRUCT {
			index = 2
		}
		if len(strArr) == 3 {
			index = 1
		}
		if index >= 0 {
			typedef.StructName = strArr[index]
			index++
			typedef.StructRename = strings.TrimRight(strArr[index], ";")
		}
	}
	return index >= 0
}

//解析宏
func (typedef *Typedef) ParseStruct(filePath string, writer *util.BufferWriter) bool {
	GetTypedefs().Add(typedef.StructName, typedef.StructRename)
	return true
}

//宏不是已\\结尾
func (typedef *Typedef) IsEndStruct(line string) bool {
	return true
}
