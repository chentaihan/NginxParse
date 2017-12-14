package main

type VariableManager struct {
	VarInfo []*VariableInfo
	checker Checker
}

func NewVariableManager() *VariableManager {
	varMgr := &VariableManager{
		checker: Checker{
			needStr:    []string{NGX_HTTP_VARIABLE_T, "[]"},
			notNeedStr: []string{"(", ")"},
		},
	}
	return varMgr
}

//判断是不是有效结构体
func (varMgr *VariableManager) Check(line string) bool {
	return varMgr.checker.Check(line)
}

//解析出ModuleInfo
func (mgr *VariableManager) parseModuleInfo(filePath string, writer *BufferWriter) *VariableInfo {
	writer.MoveNext()
	line := writer.Current()
	VarInfo := &VariableInfo{}
	VarInfo.FileName = filePath
	VarInfo.ModuleName = parseModuleName(filePath)
	VarInfo.StructName = parseStructName(line, NGX_HTTP_VARIABLE_T)
	VarInfo.StructString = writer.ToString()
	VarInfo.VarList = make([]*Variable, 0, 4)
	return VarInfo
}

//解析出结构体内容
func (mgr *VariableManager) ParseStruct(filePath string, writer *BufferWriter) bool {
	writer.Reset()
	VarInfo := mgr.parseModuleInfo(filePath, writer)
	fieldCount := 8
	lines := make([]string, fieldCount, fieldCount)
	for !writer.IsEnd() {
		i := 0
		for ; i < fieldCount && writer.MoveNext(); i++ {
			lines[i] = writer.Current()
		}

		variable := mgr.parseCommand(lines)
		if variable != nil {
			VarInfo.VarList = append(VarInfo.VarList, variable)
		}
	}

	mgr.VarInfo = append(mgr.VarInfo, VarInfo)
	return true
}

func (mgr *VariableManager) parseCommand(lines []string) *Variable {
	name := parseCommandName(lines[1])
	if name == "" {
		return nil
	}
	cmd := &Variable{}
	cmd.Name = name
	cmd.SetHandler = lines[2]
	cmd.GetHandler = lines[3]
	cmd.Data = lines[4]
	cmd.Flags = lines[5]
	cmd.Index = lines[6]
	return cmd
}
