package main

import "strings"

type Checker struct{
	needStr []string
	notNeedStr []string
}


func (ckr *Checker) Check(line string) bool{
	return ckr.checkNeedCmd(line) && !ckr.checkNotNeedCmd(line)
}

func (ckr *Checker) checkNeedCmd(line string) bool {
	for _, cmd := range ckr.needStr {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

func (ckr *Checker) checkNotNeedCmd(line string) bool {
	for _, cmd := range ckr.notNeedStr {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

