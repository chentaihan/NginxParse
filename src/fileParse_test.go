package main

import (
	"testing"
)

func Test_getDepth(t *testing.T) {
	lines := []string{
		"fjdf fkdfk fdkf ",
		"fdjkfd{f}{}{}{Fdf",
		"{dfdf",
		"{{{{{",
		"fdfdffdf}",
		"fdf{fdfdf",
		"fdfdf}Fdfdf}",
		"}}}}}}}}}}}",
		"",
		"{",
		"}",
	}

	for _,line := range lines  {
		depth := getDepth(line)
		t.Logf("%s = %d", line, depth)
	}

}