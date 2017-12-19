package main

import "strings"

type Checker struct{
	Include   []string
	IncludeNo []string
}


func (ckr *Checker) Check(line string) bool{
	return ckr.checkNeedCmd(line) && !ckr.checkNotNeedCmd(line)
}

func (ckr *Checker) checkNeedCmd(line string) bool {
	for _, cmd := range ckr.Include {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

func (ckr *Checker) checkNotNeedCmd(line string) bool {
	for _, cmd := range ckr.IncludeNo {
		if !strings.Contains(line, cmd) {
			return false
		}
	}
	return true
}

