package main

type ModuleManager struct{
	moduleInfo []*ModuleInfo
	checker    Checker
}


func NewModuleManager() *ModuleManager {
	moduleMgr := &ModuleManager{
		checker: Checker{
			needStr:    []string{NGX_MODULE_T, "{"},
			notNeedStr: []string{"(", ")", ";"},
		},
	}
	return moduleMgr
}

//判断是不是有效结构体
func (cmdMgr *ModuleManager) Check(line string) bool {
	return cmdMgr.checker.Check(line)
}

//解析出ModuleInfo
func (mgr *ModuleManager) parseModuleInfo(filePath string, writer *BufferWriter) *ModuleInfo {
	writer.MoveNext()
	line := writer.Current()
	cmdInfo := &ModuleInfo{}

	cmdInfo.FileName = filePath
	cmdInfo.ModuleName = parseModuleName(filePath)
	cmdInfo.StructName = parseStructName(line, NGX_MODULE_T)
	cmdInfo.StructString = writer.ToString()
	return cmdInfo
}

//解析出结构体内容
func (mgr *ModuleManager) ParseStruct(filePath string, writer *BufferWriter) bool {
	writer.Reset()
	cmdInfo := mgr.parseModuleInfo(filePath, writer)
	fieldCount := 8
	lines := make([]string, fieldCount, fieldCount)
	for !writer.IsEnd() {
		i := 0
		for ; i < fieldCount && writer.MoveNext(); i++ {
			lines[i] = writer.Current()
		}
		if i < fieldCount {
			break
		}
		cmd := mgr.parseModule(lines)
		cmdInfo.ModuleList = append(cmdInfo.ModuleList, cmd)
	}
	mgr.moduleInfo = append(mgr.moduleInfo, cmdInfo)
	return true
}

func (mgr *ModuleManager) parseModule(lines []string) *Module{
	module := &Module{}
	module.Context = lines[0]
	module.Command = lines[1]
	module.Type = lines[2]
	return module
}
