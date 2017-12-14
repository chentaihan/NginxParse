package main

type BaseInfo struct {
	FileName     string
	ModuleName   string
	StructName   string
	StructString string
}

type ModuleInfo struct {
	BaseInfo
	ModuleList []*Module
}

type CommandInfo struct {
	BaseInfo
	CmdList []*Command
}

type VariableInfo struct {
	BaseInfo
	VarList []*Variable
}

//对应nginx的struct ngx_command_s
type Command struct {
	Name   string //配置项名称
	Type   string //配置可以出现在哪些模块，以及参数个数
	Set    string //处理配置的函数
	Conf   string //存放配置的结构体
	Offset string //在结构体中的偏移量
	Post   string //自定义作用
}

//对应nginx的struct ngx_module_t的重要字段
type Module struct {
	Context string //上下文
	Command string //ngx_command_s 名称
	Type    string //模块类型
}

//对应nginx的struct ngx_http_variable_s,nginx变量结构体
type Variable struct {
	Name       string //上下文
	SetHandler string //ngx_command_s 名称
	GetHandler string //模块类型
	Data       string
	Flags      string
	Index      string
}
