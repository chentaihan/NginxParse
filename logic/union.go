package logic

import (
	"fmt"
	"strings"

	"github.com/chentaihan/NginxParse/util"
)

type Union struct {
	buffer util.BufferWriter
}

func NewUnion() *Union {
	return &Union{}
}

func (union *Union) IsHead(line string) bool {
	keys := []string{util.UNION, "{"}
	for _,key := range keys  {
		if index := strings.Index(line, key); index < 0{
			return false
		}
	}
	return true
}

func (union *Union) IsTail(line string) bool {
	keys := []string{"}"}
	for _,key := range keys  {
		if index := strings.Index(line, key); index < 0{
			return false
		}
	}
	return true
}

func (union *Union) AddLine(line string) bool {
	union.buffer.WriteString(line)
	return true
}

func (union *Union) Reset() {
	union.buffer.Clear()
}

//将union格式化成一行
func (union *Union) GetUnionField() string {
	union.buffer.RemoveCutChar()
	return fmt.Sprintf("%s %s", util.UNION_TAG, union.buffer.ToString())
}
