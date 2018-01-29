package util

import "testing"

func Test_GetLegalString(t *testing.T) {
	ret := util.GetLegalString(" abcd efgeg (fdfdfdf)")
	t.Log(ret)

	ret = util.GetLegalString(" abcdefg")
	t.Log(ret)

	ret = util.GetLegalString("abcdefg")
	t.Log(ret)
}


func Test_GetLegalStrings(t *testing.T) {
	ret := util.GetLegalStrings(" abcd efgeg (fdfdfdf)")
	for _, str := range ret{
		t.Log(str)
	}

	ret = util.GetLegalStrings(" abcd efgeg (fdfdfdf")
	for _, str := range ret{
		t.Log(str)
	}

	ret = util.GetLegalStrings("abcd efgeg (fdfdfdf ")
	for _, str := range ret{
		t.Log(str)
	}
}