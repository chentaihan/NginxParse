package main

type CommandManager struct {
	CmdInfo []*CommandInfo
	checker Checker
}

func NewCommandManager() *CommandManager {
	cmdMgr := &CommandManager{
		checker: Checker{
			needStr:    []string{NGX_COMMAND_T, "[]"},
			notNeedStr: []string{"(", ")"},
		},
	}
	return cmdMgr
}

//判断是不是有效结构体
func (cmdMgr *CommandManager) Check(line string) bool {
	return cmdMgr.checker.Check(line)
}

//解析出ModuleInfo
func (mgr *CommandManager) parseModuleInfo(filePath string, writer *BufferWriter) *CommandInfo {
	writer.MoveNext()
	line := writer.Current()
	cmdInfo := &CommandInfo{}
	cmdInfo.FileName = filePath
	cmdInfo.ModuleName = parseModuleName(filePath)
	cmdInfo.StructName = parseStructName(line, NGX_COMMAND_T)
	cmdInfo.StructString = writer.ToString()
	cmdInfo.CmdList = make([]*Command, 0, 4)
	return cmdInfo
}

//解析出结构体内容
func (mgr *CommandManager) ParseStruct(filePath string, writer *BufferWriter) bool {
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
		cmd := mgr.parseCommand(lines)
		cmdInfo.CmdList = append(cmdInfo.CmdList, cmd)
	}

	mgr.CmdInfo = append(mgr.CmdInfo, cmdInfo)
	return true
}

func (mgr *CommandManager) parseCommand(lines []string) *Command {
	cmd := &Command{}
	cmd.Name = parseCommandName(lines[1])
	cmd.Type = lines[2]
	cmd.Set = lines[3]
	cmd.Conf = lines[4]
	cmd.Offset = lines[5]
	cmd.Post = lines[6]
	return cmd
}
