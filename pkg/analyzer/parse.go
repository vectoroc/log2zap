package analyzer

import (
	"regexp"
	"strings"
)

var reWithNames = regexp.MustCompile(`((?:[\w_-]+=)?\\?['"]?%[^ '"]\\?['"]?)`)

func parseFormat(format string) []formatVar {
	var res []formatVar
	for _, match := range reWithNames.FindAllString(format, -1) {
		item := formatVar{raw: match}
		item.format = match
		parts := strings.Split(match, "=")
		if len(parts) != 1 {
			item.format = parts[1]
			item.name = parts[0]
		}

		item.format = strings.Trim(item.format, ` '"\`)
		res = append(res, item)
	}

	return res
}

type formatVar struct {
	raw    string
	name   string
	format string
}

func cleanUpFormatString(format string, vars []formatVar) string {
	for _, v := range vars {
		format = strings.Replace(format, v.raw, "", 1)
	}
	return strings.Trim(strings.ReplaceAll(format, "  ", " "), ".,: \n\t")
}
