package main

type StructInfo struct {
	FileName     string
	ModuleName   string
	StructName   string
	StructString string
}

type TableInfo struct {
	StructInfo
	Title   []string
	Content [][]string
}
