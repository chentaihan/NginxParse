package logic

/**
所有结构体实例
 */

import "fmt"

type Assignments struct {
	Tables []*TableInfo
}

var assignments *Assignments = nil

func GetAssignments() *Assignments {
	if assignments == nil {
		assignments = &Assignments{
			Tables: make([]*TableInfo, 0, 128),
		}
	}
	return assignments
}

func (asss *Assignments) Add(table *TableInfo) {
	asss.Tables = append(asss.Tables, table)
}

func (asss *Assignments) Print() {
	fmt.Println()
	fmt.Println("-----------------------------------Assignments------------------")
	for _, table := range asss.Tables {
		fmt.Println("assignment", table.StructName, table.VarName, "----------------")
		for _, content := range table.Content {
			fmt.Println("{")
			for _, field := range content {
				if len(table.Content) == 1 && (field == "{" || field == "}") {
					fmt.Println(table.StructString)
				}
				fmt.Println(field)
			}
			fmt.Println("}")
		}
		fmt.Println()
	}
}
