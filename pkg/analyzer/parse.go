package analyzer

import (
	"regexp"
	"strings"
)

var rePercentPlaceholders = regexp.MustCompile(`((?:[\w_-]+=)?\\?['"]?%[^ '"]+\\?['"]?)`)
var reStaticVars = regexp.MustCompile(`([\w_-]+=\\?['"]?[\w_-]+\\?['"]?)`)

func parse(format string, re *regexp.Regexp) []formatVar {
	var res []formatVar
	for _, match := range re.FindAllString(format, -1) {
		item := formatVar{raw: match}
		item.value = match
		parts := strings.Split(match, "=")
		if len(parts) != 1 {
			item.value = parts[1]
			item.key = parts[0]
		}

		item.value = strings.Trim(item.value, ` '"\`)
		res = append(res, item)
	}

	return res
}

func ParseFormat(format string) []formatVar {
	return parse(format, rePercentPlaceholders)
}

func ParseStaticVars(msg string) []formatVar {
	return parse(msg, reStaticVars)
}

type formatVar struct {
	raw   string
	key   string
	value string
}

func CleanUpFormatString(format string, vars []formatVar) string {
	pairs := []string{"  ", " ", " .", ".", " ,", ","}
	for _, v := range vars {
		pairs = append(pairs, v.raw, "")
	}
	format = strings.NewReplacer(pairs...).Replace(format)
	format = strings.NewReplacer("  ", " ", " .", ".", " ,", ",").Replace(format)
	return strings.Trim(format, ".,: \n\t")
}
