package main

type IParse interface {
	Check(line string) bool                                 //判断是不是有效结构体
	IsEndStruct(line string) bool                           //结构是否结束
	FormatStruct(bufWriter *BufferWriter) *BufferWriter     //格式化struct
	ParseStruct(filePath string, writer *BufferWriter) bool //解析出结构体内容
}
