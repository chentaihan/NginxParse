package logic

import "fmt"

/**
所有结构体的定义
 */

type StructInfo struct {
	FileName     string
	ModuleName   string
	StructName   string
	Rename       string
	Fields       []string
	StructString string
}

type TableInfo struct {
	StructInfo
	VarName string
	Content [][]string
}

var defines *Defines

type Defines struct {
	Map map[string]*StructInfo
}

func GetDefines() *Defines {
	if defines == nil {
		defines = &Defines{
			Map : make(map[string]*StructInfo, 128),
		}
		return defines
	}
	return defines
}

func (defines *Defines) Add(structName string, sttInfo *StructInfo) bool {
	if sttInfo != nil {
		defines.Map[sttInfo.StructName] = sttInfo
		return true
	}
	return false
}

func (defines *Defines) Get(structName string) *StructInfo {
	if sttInfo, ok := defines.Map[structName]; ok {
		return sttInfo
	}
	return nil
}

func (defines *Defines) Size() int {
	return len(defines.Map)
}

func (defines *Defines) StructNameList() []string {
	structNames := make([]string, 0, defines.Size())
	for name, _ := range defines.Map {
		structNames = append(structNames, name)
	}
	return structNames
}

func (defines *Defines) Structlist() []*StructInfo {
	structList := make([]*StructInfo, 0, defines.Size())
	for _, structInfo := range defines.Map {
		structList = append(structList, structInfo)
	}
	return structList
}

func (defines *Defines) Print(){
	fmt.Println()
	fmt.Println("-----------------Defines-------------------------")
	for _,stt := range defines.Map{
		if stt != nil {
			fmt.Println("struct ", stt.StructName, "------------------")
			for _, field := range stt.Fields{
				fmt.Println(field)
			}
			fmt.Println()
		}
	}
}
