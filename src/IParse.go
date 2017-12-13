package main

type IParse interface {
	Check(line string) bool		//判断是不是有效结构体
	ParseStruct(filePath string, writer *BufferWriter) bool	//解析出结构体内容
}
