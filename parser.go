package frienvironment

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse(name, code string) Program {
	re := regexp.MustCompile(`"[^"]*"|\S+`)
	parsedCode := re.FindAllString(preParse(code), -1)
	for _, p := range parsedCode {
		fmt.Println(p)
	}
	return LoadProgram(name, parsedCode)
}

func preParse(code string) string {
	out := ""
	for _, line := range strings.Split(code, "\n") {

		out += strings.Split(line, "//")[0] + "\n"
	}
	return out
}
