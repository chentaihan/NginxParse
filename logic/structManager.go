package logic


/**
结构体的定义
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
	Title   []string
	Content [][]string
}

var structManger *StructManager

type StructManager map[string]*StructInfo


func GetStructManager() *StructManager {
	if structManger == nil {
		structManger = &StructManager{}
		return structManger
	}
	return structManger
}

func (sttMgr *StructManager) Add(sttInfo *StructInfo) bool {
	if sttInfo != nil {
		(*sttMgr)[sttInfo.StructName] = sttInfo
		return true
	}
	return false
}

func (sttMgr *StructManager) Get(structName string) *StructInfo {
	if sttInfo, ok := (*sttMgr)[structName]; ok {
		return sttInfo
	}
	return nil
}

func (sttMgr *StructManager) Size() int {
	return len(*sttMgr)
}

func (sttMgr *StructManager) StructNameList() []string {
	structNames := make([]string, 0, sttMgr.Size())
	for name, _ := range *sttMgr {
		structNames = append(structNames, name)
	}
	return structNames
}

func (sttMgr *StructManager) Structlist() []*StructInfo {
	structList := make([]*StructInfo, 0, sttMgr.Size())
	for _, structInfo := range *sttMgr{
		structList = append(structList, structInfo)
	}
	return structList
}
