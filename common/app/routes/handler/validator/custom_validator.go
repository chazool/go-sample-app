package validator

import "regexp"

func containsOnly(field string, regExpstr string) bool {
	reg, _ := regexp.Compile(regExpstr)
	return !reg.MatchString(field)
}
