package frienvironment

import (
	"fmt"
	"regexp"
)

func Parse(name, code string) Program {

	re := regexp.MustCompile(`"[^"]*"|\S+`)

	// 2. Find all matches in the input string.
	parsedCode := re.FindAllString(code, -1)

	for _, p := range parsedCode {
		fmt.Println(p)
	}

	return LoadProgram(name, parsedCode)
}
