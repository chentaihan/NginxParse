package logic

import (
	"testing"
)

func Test_isStructHeader(t *testing.T) {

	list := []string{
		"typedef struct {",
		"typedef struct  {",
		"typedef struct name {",
		"typedef struct name  {",
		"typedef struct name name {",
		"typedef struct   name name {",
		"typedef struct n&ame {",
		"typede f struct name name {",
		" typede f stur ct name name {",

		"struct  {",
		"struct name {",
		"struct  name  {",
		"struct na$me {",
		"struct name name {",
		"struct name  name  name {",
		"tydsds pedef struct   name name {",

		"typ ede f struct name name {",
		" typede f stur ct name name {",
		"fdj fdjkfjd fdjfk fdkj ",
		"fdfdf",
		"fdfj fdfj ",
		"fdf fdfjd fdhjf",
	}

	var define Define
	for _, str := range list {
		isHeader := define.IsStartStruct(str)
		t.Logf("%s = %t", str, isHeader)
	}

}
