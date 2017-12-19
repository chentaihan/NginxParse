package main

import "testing"

func Test_getLegalString(t *testing.T) {
	ret := getLegalString(" abcd efgeg (fdfdfdf)")
	t.Log(ret)

	ret = getLegalString(" abcdefg")
	t.Log(ret)

	ret = getLegalString("abcdefg")
	t.Log(ret)
}


func Test_getLegalStrings(t *testing.T) {
	ret := getLegalStrings(" abcd efgeg (fdfdfdf)")
	for _, str := range ret{
		t.Log(str)
	}

	ret = getLegalStrings(" abcd efgeg (fdfdfdf")
	for _, str := range ret{
		t.Log(str)
	}

	ret = getLegalStrings("abcd efgeg (fdfdfdf ")
	for _, str := range ret{
		t.Log(str)
	}
}