package frienvironment

import (
	"fmt"

	"strings"
)

type InterruptType int

const (
	InterruptTypeNULL       InterruptType = 0
	InterruptTypeKeyPressed InterruptType = 1
	InterruptTypeMouseDown  InterruptType = 2
)

type Program struct {
	name  string
	code  []string
	ptr   int
	rp    *Stack[int]
	state string

	functions map[string]int
	labels    map[string]int
	vars      map[string]int

	interruptHandlers map[InterruptType]int
}

func LoadProgram(name string, code []string) Program {
	prog := Program{
		name:              name,
		code:              code,
		functions:         map[string]int{},
		labels:            map[string]int{},
		vars:              map[string]int{},
		interruptHandlers: map[InterruptType]int{},
		rp:                &Stack[int]{},
	}

	nextFreeCell := 1000

	for i, step := range code {
		if strings.HasSuffix(step, ":") {
			prog.labels[strings.TrimSuffix(step, ":")] = i + 1
		}
	}

	for i, step := range code {
		if step == "DEF" {
			prog.functions[code[i+1]] = i + 2
		}
	}

	for i, step := range code {
		if strings.HasPrefix(step, "&") {
			function := strings.TrimPrefix(step, "&")
			if ref, exists := prog.functions[function]; exists {
				code[i] = fmt.Sprintf("%d", ref)
			} else {
				panic(fmt.Errorf("cannot resolve reference for: %s", function))
			}
		}
	}

	for i, step := range code {
		if step == "VAR" {
			prog.vars[code[i+1]] = nextFreeCell
			nextFreeCell++
		}
	}

	return prog
}
