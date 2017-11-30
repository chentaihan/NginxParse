package main

type CommandInfo struct {
	FileName   string
	ModuleName string
	StructName string
	CmdList    []Command
}

//对应nginx的struct ngx_command_s
type Command struct {
	Name    string //配置项名称
	TypeVal string //配置可以出现在哪些模块，以及参数个数
	Set     string //处理配置的函数
	Conf    string //存放配置的结构体
	Offset  string //在结构体中的偏移量
	Post    string //自定义作用
}
