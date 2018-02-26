package logic

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_macro(t *testing.T) {
	macro := GetMacros()
	macro.Add("NGX_MODULE_SIGNATURE_5", &MacroInfo{
		Name:  "NGX_MODULE_SIGNATURE_5",
		Value: "NGX_MODULE_SIGNATURE_5",
	})

	macro.Add("NGX_MODULE_SIGNATURE_12", &MacroInfo{
		Name:  "NGX_MODULE_SIGNATURE_12",
		Value: "NGX_MODULE_SIGNATURE_12",
	})
	macro.Add("NGX_MODULE", &MacroInfo{
		Name:  "NGX_MODULE",
		Value: "",
	})

	macro.Add("ngx_conf_merge_uint_value", &MacroInfo{
		Name:  "(conf, prev, default)",
		Value: "if(conf == NGX_CONF_UNSET_UINT){conf = (prev == NGX_CONF_UNSET_UINT) ? default : prev;}",
	})

	list := macro.Map
	for key, val := range list {
		if val.Name != "" {
			fmt.Printf("%s%s=%s", key, val.Name, val.Value)
			util.Println()
		} else {
			fmt.Printf("-------------------------%s%s=%s", key, val.Name, val.Value)
			util.Println()
		}
	}

	t.Log(macro.Exist("NGX_MODULE_SIGNATURE_5"))
	t.Log(macro.Exist("NGX_MODULE_SIGNATURE_511"))
	t.Log(macro.GetMacroValue("NGX_MODULE_SIGNATURE_12"))
	t.Log(macro.GetMacroValue("NGX_MODULE"))
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
