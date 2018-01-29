package logic

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_macro(t *testing.T) {
	macro := GetMacro()
	macro.AddMacroInfo("NGX_MODULE_SIGNATURE_5", &MacroInfo{
		Name:  "NGX_MODULE_SIGNATURE_5",
		Value: "NGX_MODULE_SIGNATURE_5",
	})

	macro.AddMacroInfo("NGX_MODULE_SIGNATURE_12", &MacroInfo{
		Name:  "NGX_MODULE_SIGNATURE_12",
		Value: "NGX_MODULE_SIGNATURE_12",
	})

	list := macro.MacroList
	for key, val := range list {
		if val.Name != "" {
			fmt.Printf("%s%s=%s", key, val.Name, val.Value)
			fmt.Println()
		} else {
			fmt.Printf("-------------------------%s%s=%s", key, val.Name, val.Value)
			fmt.Println()
		}
	}

	t.Log(macro.Exist("NGX_MODULE_SIGNATURE_5"))
	t.Log(macro.Exist("NGX_MODULE_SIGNATURE_511"))
	t.Log(macro.GetMacroValue("NGX_MODULE_SIGNATURE_12"))
	t.Log(macro.GetMacroValue("NGX_MODULE_SIGNATURE_1211"))
	t.Log(macro.GetMacroValue("ngx_conf_merge_uint_value(conf__, prev__, default__)"))
	t.Log(macro.GetMacroValue("ngx_conf_merge_uint_value(conf__, prev__)"))
	t.Log(macro.GetMacroValue("ngx_conf_merge_uint_value()"))

}

func Test_Map(t *testing.T) {
	mapList := make(map[int]MacroInfo, 4)
	for i := 0; i < 10; i++ {
		mapList[i] = MacroInfo{
			Name:  strconv.Itoa(i),
			Value: strconv.Itoa(i),
		}
	}

	t.Log(mapList[0].Value)
	t.Log(mapList[4].Value)
	t.Log(mapList[9].Value)
}
