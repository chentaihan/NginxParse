package logic

/**
struct所有重定义
 */
import "fmt"

type Typedefs struct {
	Map map[string]string
}

var typedefs *Typedefs = nil

func GetTypedefs() *Typedefs {
	if typedefs == nil {
		typedefs = &Typedefs{
			Map: make(map[string]string, 64),
		}
	}
	return typedefs
}

func (tfs *Typedefs) Add(structName, structRename string) {
	tfs.Map[structName] = structRename
}

func (tfs *Typedefs) Parse(defines *Defines) {
	for key, value := range tfs.Map {
		if _, exist := defines.Map[key]; !exist {
			delete(tfs.Map, key)
		} else {
			if define, exist := defines.Map[value]; !exist {
				defines.Add(value, define)
			}
		}
	}
}

func (tfs *Typedefs) Print() {
	fmt.Println()
	fmt.Println("----------------typedef----------------")
	for key, value := range tfs.Map {
		fmt.Println("typedef struct ", key, value)
	}
}
