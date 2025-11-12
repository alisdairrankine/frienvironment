package frienvironment

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type VM struct {
	stack   Stack[int]
	words   map[string]func()
	program Program
	buf     []byte
	memory  [4096]int

	docs string
}

func (vm *VM) Run() error {
	vm.program.state = "running"
	for vm.program.ptr < len(vm.program.code) && vm.program.state != "halted" {
		if vm.program.state == "waiting" {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		step := vm.program.code[vm.program.ptr]
		vm.program.ptr++

		if word, exists := vm.words[step]; exists {
			word()
			continue
		}

		if funcPtr, exists := vm.program.functions[step]; exists {
			vm.program.rp.Push(vm.program.ptr)
			vm.program.ptr = funcPtr
			continue
		}

		if n, err := strconv.Atoi(step); err == nil {
			vm.stack.Push(n)
			continue
		}

		if strings.HasSuffix(step, ":") {
			continue
		}

		if strings.HasPrefix(step, ":") {
			if jmp, exists := vm.program.labels[strings.TrimPrefix(step, ":")]; exists {
				vm.program.ptr = jmp
			}
			continue
		}

		if strings.HasPrefix(step, `"`) && strings.HasSuffix(step, `"`) {
			vm.stack.PushMany(stringToIntSlice(strings.Trim(step, `"`))...)
			continue
		}

		if addr, exists := vm.program.vars[step]; exists {
			vm.stack.Push(addr)
		}

		return fmt.Errorf("unrecognised word: %s", step)

	}
	return nil

}
func (vm *VM) AddWord(word, docs string, wordFn func()) {
	vm.words[word] = wordFn
	vm.docs += fmt.Sprintf("- %q: %s\n", word, docs)
}

func (vm *VM) init() {

	vm.AddWord("+", "(a b -- a+b) Adds 2 numbers", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(b + a)
	})

	vm.AddWord("-", "(a b -- b-a) Subtracts a number", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(b - a)
	})
	vm.AddWord("/", "(a b -- b/a) Divides a number", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(b / a)
	})

	vm.AddWord("*", "(a b -- a*b) Multiplies 2 numbers", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(b * a)
	})

	vm.AddWord("=", "(a b -- a=b?1:0) Checks equality", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		if b == a {
			vm.stack.Push(1)
		} else {
			vm.stack.Push(0)
		}
	})

	vm.AddWord("<", "(a b -- b<a?1:0) Checks b is less than a", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		if b < a {
			vm.stack.Push(1)
		} else {
			vm.stack.Push(0)
		}
	})

	vm.AddWord(">", "(a b -- b>a?1:0) Checks b is greater than a", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		if b > a {
			vm.stack.Push(1)
		} else {
			vm.stack.Push(0)
		}
	})

	vm.AddWord(".", "(a -- ) Outputs top stack value", func() {
		fmt.Println(vm.stack.Pop())
	})

	vm.AddWord("DUP", "(a -- a a) Duplicates top stack value", func() {
		a := vm.stack.Pop()
		vm.stack.Push(a)
		vm.stack.Push(a)
	})

	vm.AddWord("DROP", "(a -- ) Removes top stack value", func() {
		vm.stack.Pop()
	})

	vm.AddWord("SWAP", "(a b -- b a) Swaps top 2 stack values", func() {
		a := vm.stack.Pop()
		b := vm.stack.Pop()
		vm.stack.Push(a)
		vm.stack.Push(b)
	})

	vm.AddWord("STORE", "(a b -- ) Stores a in memory address b", func() {
		addr := vm.stack.Pop()
		data := vm.stack.Pop()
		vm.memory[addr] = data
	})

	vm.AddWord("LOAD", "(a -- b) Loads memory address a", func() {
		addr := vm.stack.Pop()
		vm.stack.Push(vm.memory[addr])
	})

	vm.AddWord("IF", "If top stack value is 1, continue to next step, otherwise skip to ELSE", func() {
		condition := vm.stack.Pop()
		if condition == 0 {
			// Skip to ELSE or THEN
			depth := 1
			for depth > 0 {
				if vm.program.ptr >= len(vm.program.code) {
					panic("Unmatched IF: program pointer out of bounds")
				}
				step := vm.program.code[vm.program.ptr]
				vm.program.ptr++
				switch step {
				case "IF":
					depth++
				case "ELSE":
					if depth == 1 {
						depth--
					}
				case "THEN":
					depth--
				}
			}
		}
	})

	vm.AddWord("ELSE", "", func() {
		// Skip to THEN
		depth := 1
		for depth > 0 {
			if vm.program.ptr >= len(vm.program.code) {
				panic("Unmatched ELSE: program pointer out of bounds")
			}
			step := vm.program.code[vm.program.ptr]
			vm.program.ptr++
			switch step {
			case "IF":
				depth++
			case "THEN":
				depth--
			}
		}
	})

	vm.AddWord("THEN", "End of IF", func() {
		// No-op, handled by control flow
	})

	vm.AddWord("HALT", "Stops the VM", func() {
		vm.program.state = "halted"
	})

	vm.AddWord("GOTO", "(a --) Jump to top stack step", func() {
		vm.program.ptr = vm.stack.Pop()
	})

	vm.AddWord("RET", "Return to previous return pointer", func() {
		vm.program.ptr = vm.program.rp.Pop()
	})

	vm.AddWord("DEF", "Define a function", func() {

	})

	vm.AddWord("BUF", "(a -- ) Store byte in output buffer", func() {
		b := vm.stack.Pop()
		vm.buf = append(vm.buf, byte(b))
	})

	// vm.AddWord("SYSCALL.OUT", func() {
	// 	fmt.Println(string(vm.buf))
	// 	vm.buf = []byte{}
	// })

	vm.AddWord("YIELD", "Set the program into a waiting state", func() {
		vm.program.state = "waiting"
	})

	vm.AddWord("SLEEP", "(a -- ) Sleep for top stack value milliseconds", func() {
		s := vm.stack.Pop()

		time.Sleep(time.Millisecond * time.Duration(s))
	})

	vm.AddWord("VAR", "allocate a variable", func() {
		vm.program.ptr++
	})

	vm.AddWord("SYSCALL.INTERRUPT.REGISTER", "DO NOT USE: Experimental interrupt system", func() {
		functionPtr := vm.stack.Pop()
		interruptType := vm.stack.Pop()
		vm.program.interruptHandlers[InterruptType(interruptType)] = functionPtr
	})

}

func NewVM(app Program) *VM {
	vm := &VM{
		stack:   Stack[int]{},
		words:   map[string]func(){},
		program: app,
		memory:  [4096]int{},
	}
	vm.init()

	return vm
}

func (vm *VM) WriteToBuffer(b []byte) {
	vm.buf = append(vm.buf, b...)
}

func (vm *VM) ClearBuffer() {
	vm.buf = []byte{}
}

func (vm *VM) ReadFromBuffer() []byte {
	return vm.buf
}

func (vm *VM) State() string {
	return vm.program.state
}

func (vm *VM) Interrupt(typ InterruptType) {
	if handler, exists := vm.program.interruptHandlers[typ]; exists {
		vm.program.state = "waiting"
		vm.program.rp.Push(vm.program.ptr)
		vm.program.ptr = handler
	}
}

func (vm *VM) Docs() string {
	builtinDocs := `
Notes:
- Labels are suffixed with ":"
- Jump to a label by prefixing the label name with ":"
- Functions are defined with "DEF"
- Pointers to functions are prefixed with "&"
- Strings are quoted with "'"

Example Program:
` + "```" + `
'Loading... (1)'
PRINT-STRING
1000
SLEEP

'Loading... (2)'
PRINT-STRING
1000
SLEEP

'Loading... (3)'
PRINT-STRING
1000
SLEEP

'Hello World!'
PRINT-STRING
YIELD

DEF PRINT-STRING
    out: BUF
    DUP 0 = IF
        DROP
        SYSCALL.OUT
        RET
    ELSE
        :out
    THEN
RET
` + "```\n"
	return "Defined Words:\n" + vm.docs + builtinDocs
}

func stringToIntSlice(s string) []int {
	i := []int{}

	for _, b := range []byte(s) {
		i = append(i, int(b))
	}
	i = append(i, 0)
	slices.Reverse(i)
	return i
}
