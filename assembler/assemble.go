package assembler

import (
	"strconv"
	"strings"

	"github.com/alisdairrankine/frienvironment/vm"
)

func Assemble(program string) []byte {
	out := []byte{}
	for _, line := range strings.Split(program, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			continue
		}
		out = append(out, ParseLine(line)...)
	}
	return out
}

func ParseLine(line string) []byte {
	var expect2Bytes bool
	parts := strings.Split(line, " ")
	out := []byte{}
	switch strings.ToLower(parts[0]) {
	case "yield":
		out = append(out, vm.YieldInstruction)
	case "halt":
		out = append(out, vm.HaltInstruction)
	case "dup":
		out = append(out, vm.DupInstruction)
	case "drop":
		out = append(out, vm.DropInstruction)
	case "swap":
		out = append(out, vm.SwapInstruction)
	case "rot":
		out = append(out, vm.RotInstruction)
	case "over":
		out = append(out, vm.OverInstruction)
	case "nip":
		out = append(out, vm.NipInstruction)
	case "tuck":
		out = append(out, vm.TuckInstruction)
	case "tor":
		out = append(out, vm.ToRInstruction)
	case "fromr":
		out = append(out, vm.FromRInstruction)
	case "fetchr":
		out = append(out, vm.FetchRInstruction)
	case "add":
		out = append(out, vm.AddInstruction)
	case "add16":
		expect2Bytes = true
		out = append(out, vm.Add16Instruction)
	case "sub":
		out = append(out, vm.SubInstruction)
	case "sub16":
		expect2Bytes = true
		out = append(out, vm.Sub16Instruction)
	case "mul":
		out = append(out, vm.MulInstruction)
	case "mul16":
		expect2Bytes = true
		out = append(out, vm.Mul16Instruction)
	case "div":
		out = append(out, vm.DivInstruction)
	case "div16":
		expect2Bytes = true
		out = append(out, vm.Div16Instruction)
	case "mod":
		out = append(out, vm.ModInstruction)
	case "mod16":
		expect2Bytes = true
		out = append(out, vm.Mod16Instruction)
	case "and":
		out = append(out, vm.AndInstruction)
	case "and16":
		expect2Bytes = true
		out = append(out, vm.And16Instruction)
	case "or":
		out = append(out, vm.OrInstruction)
	case "or16":
		expect2Bytes = true
		out = append(out, vm.Or16Instruction)
	case "xor":
		out = append(out, vm.XorInstruction)
	case "xor16":
		expect2Bytes = true
		out = append(out, vm.Xor16Instruction)
	case "not":
		out = append(out, vm.NotInstruction)
	case "not16":
		expect2Bytes = true
		out = append(out, vm.Not16Instruction)
	case "inc":
		out = append(out, vm.IncInstruction)
	case "inc16":
		expect2Bytes = true
		out = append(out, vm.Inc16Instruction)
	case "dec":
		out = append(out, vm.DecInstruction)
	case "dec16":
		expect2Bytes = true
		out = append(out, vm.Dec16Instruction)
	case "shl":
		out = append(out, vm.ShlInstruction)
	case "shl16":
		expect2Bytes = true
		out = append(out, vm.Shl16Instruction)
	case "shr":
		out = append(out, vm.ShrInstruction)
	case "shr16":
		expect2Bytes = true
		out = append(out, vm.Shr16Instruction)
	case "jz":
		out = append(out, vm.JzInstruction)
	case "jnz":
		out = append(out, vm.JnzInstruction)
	case "call":
		out = append(out, vm.CallInstruction)
	case "ret":
		out = append(out, vm.RetInstruction)
	case "eq":
		out = append(out, vm.EqInstruction)
	case "eq16":
		expect2Bytes = true
		out = append(out, vm.Eq16Instruction)
	case "nq":
		out = append(out, vm.NqInstruction)
	case "nq16":
		expect2Bytes = true
		out = append(out, vm.Nq16Instruction)
	case "gt":
		out = append(out, vm.GtInstruction)
	case "gt16":
		expect2Bytes = true
		out = append(out, vm.Gt16Instruction)
	case "lt":
		out = append(out, vm.LtInstruction)
	case "lt16":
		expect2Bytes = true
		out = append(out, vm.Lt16Instruction)
	case "push":
		out = append(out, vm.PushInstruction)
	case "push16":
		expect2Bytes = true
		out = append(out, vm.Push16Instruction)
	case "store":
		out = append(out, vm.StoreInstruction)
	case "store16":
		out = append(out, vm.Store16Instruction)
	case "load":
		out = append(out, vm.LoadInstruction)
	case "load16":
		expect2Bytes = true
		out = append(out, vm.Load16Instruction)
	}

	if len(parts) == 2 {
		data := parts[1]

		if strings.HasPrefix(data, "'") {
			out = append(out, []byte(strings.Trim(data, "'"))...)

		} else {
			base := 10
			if strings.HasPrefix(data, "0x") {
				data = strings.TrimPrefix(data, "0x")
				base = 16
			}
			d, _ := strconv.ParseUint(data, base, 16)
			if expect2Bytes {
				high := byte((d & 0xFF00) >> 8)
				out = append(out, high)
			}
			low := 0xFF & byte(d)
			out = append(out, low)

		}
	}

	return out
}
