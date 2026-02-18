package vm

import (
	"encoding/binary"
	"fmt"
)

type VM struct {
	MMIO *MMIO

	pc uint16

	interruptChan chan uint16

	running bool

	Debug bool
}

const (
	FlagFault                = 0b00000001
	FlagWaiting              = 0b00000010
	FlagStackOverflow        = 0b00000100
	FlagStackUnderflow       = 0b00001000
	FlagReturnStackOverflow  = 0b00010000
	FlagReturnStackUnderflow = 0b00100000
	FlagDivideByZero         = 0b01000000
)

func (vm *VM) RegisterDevice(num int, device Device) {
	vm.MMIO.devices[num] = device
}

const (
	AddrEntrypoint uint16 = 0x0000

	AddrStatus uint16 = 0x0004

	AddrStackPointer uint16 = 0x0002
	AddrStackStart   uint16 = 0x0100
	AddrStackEnd     uint16 = 0x01FF

	AddrReturnStackPointer uint16 = 0x0003
	AddrReturnStackStart   uint16 = 0x0200
	AddrReturnStackEnd     uint16 = 0x02FF
)

func New() *VM {
	return &VM{
		MMIO: &MMIO{
			data: [65536]byte{},
		},
		interruptChan: make(chan uint16),
	}
}

func (vm *VM) LoadProgram(data []byte) {
	for i, b := range data {
		vm.MMIO.WriteByte(uint16(i+0x400), b)
	}
	vm.MMIO.WriteByte(0, 0x04)
}

func (vm *VM) PushStack(b byte) {
	sp := uint16(vm.MMIO.ReadByte(AddrStackPointer))
	start := AddrStackStart
	if sp > 0xFE {
		vm.setFault(FlagStackOverflow)
		return
	}
	vm.MMIO.WriteByte(start+sp, b)
	vm.MMIO.WriteByte(AddrStackPointer, vm.MMIO.ReadByte(AddrStackPointer)+1)
}

func (vm *VM) PopStack() byte {
	sp := uint16(vm.MMIO.ReadByte(AddrStackPointer))
	if sp == 0x00 {
		vm.setFault(FlagStackUnderflow)
		return 0
	}
	sp1 := sp - 1

	b := vm.MMIO.ReadByte(AddrStackStart + sp1)
	vm.MMIO.WriteByte(AddrStackPointer, vm.MMIO.ReadByte(AddrStackPointer)-1)
	return b
}
func (vm *VM) PushStack16(b uint16) {
	sp := uint16(vm.MMIO.ReadByte(AddrStackPointer))
	start := AddrStackStart
	if sp >= 0xFE {
		vm.setFault(FlagStackOverflow)
		return
	}
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, b)

	sp1 := sp + 1
	vm.MMIO.WriteByte(start+sp, bs[0])
	vm.MMIO.WriteByte(start+sp1, bs[1])
	vm.MMIO.WriteByte(AddrStackPointer, vm.MMIO.ReadByte(AddrStackPointer)+2)

}

func (vm *VM) PopStack16() uint16 {
	sp := uint16(vm.MMIO.ReadByte(AddrStackPointer))
	start := AddrStackStart
	if sp < 0x02 {
		vm.setFault(FlagStackUnderflow)
		return 0
	}
	sp1 := sp - 1
	sp2 := sp - 2
	a := vm.MMIO.ReadByte(start + sp2)
	b := vm.MMIO.ReadByte(start + sp1)
	vm.MMIO.WriteByte(AddrStackPointer, vm.MMIO.ReadByte(AddrStackPointer)-2)
	return (uint16(a) << 8) | uint16(b)
}

func (vm *VM) PushReturnStack(b uint16) {
	sp := uint16(vm.MMIO.ReadByte(AddrReturnStackPointer))
	start := AddrReturnStackStart
	if sp >= 0xFE {
		vm.setFault(FlagReturnStackOverflow)
		return
	}
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, b)

	sp1 := sp + 1

	vm.MMIO.WriteByte(start+sp, bs[0])
	vm.MMIO.WriteByte(start+sp1, bs[1])
	vm.MMIO.WriteByte(AddrReturnStackPointer, vm.MMIO.ReadByte(AddrReturnStackPointer)+2)
}

func (vm *VM) PopReturnStack() uint16 {
	sp := uint16(vm.MMIO.ReadByte(AddrReturnStackPointer))

	start := AddrReturnStackStart
	if sp < 0x02 {
		vm.setFault(FlagReturnStackUnderflow)
		return 0
	}
	sp1 := sp - 1
	sp2 := sp - 2
	a := vm.MMIO.ReadByte(start + sp2)
	b := vm.MMIO.ReadByte(start + sp1)
	vm.MMIO.WriteByte(AddrReturnStackPointer, vm.MMIO.ReadByte(AddrReturnStackPointer)-2)
	return (uint16(a) << 8) | uint16(b)
}

func (vm *VM) setFault(flags byte) {
	vm.MMIO.WriteByte(AddrStatus, flags|FlagFault)
	vm.running = false
}

func (vm *VM) SetFlag(flag byte) {
	vm.MMIO.WriteByte(AddrStatus, vm.MMIO.ReadByte(AddrStatus)|flag)
}

func (vm *VM) CheckFlag(flag byte) bool {
	return vm.MMIO.ReadByte(AddrStatus)&flag == flag
}

func (vm *VM) UnsetFlag(flag byte) {
	vm.MMIO.WriteByte(AddrStatus, vm.MMIO.ReadByte(AddrStatus) & ^flag)
}

func (vm *VM) advancePC(n uint16) {
	vm.pc += n
}

func (vm *VM) Run() {
	vm.pc = binary.BigEndian.Uint16(
		[]byte{
			vm.MMIO.ReadByte(0),
			vm.MMIO.ReadByte(1),
		},
	)
	vm.running = true

	go func() {
		for vm.running {
			vm.execute(vm.MMIO.ReadByte(vm.pc))
			if vm.CheckFlag(FlagWaiting) {
				callbackPtr := <-vm.interruptChan

				addr := make([]byte, 2)
				addr[0] = vm.MMIO.ReadByte(callbackPtr)
				addr[1] = vm.MMIO.ReadByte(callbackPtr + 1)
				vm.pc = binary.BigEndian.Uint16(addr)
				vm.UnsetFlag(FlagWaiting)
				fmt.Println("continue")
			}
		}
	}()
}

func (vm *VM) Stop() {
	vm.running = false
}

func (vm *VM) execute(instr byte) {
	if vm.Debug {
		fmt.Println(InstrName[instr])
	}
	switch instr {
	case YieldInstruction:
		vm.SetFlag(FlagWaiting)
	case HaltInstruction:
		vm.UnsetFlag(FlagFault)
		vm.running = false
		vm.advancePC(1)
	case DupInstruction:
		a := vm.PopStack()
		vm.PushStack(a)
		vm.PushStack(a)
		vm.advancePC(1)
	case DropInstruction:
		vm.PopStack()
		vm.advancePC(1)
	case SwapInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(b)
		vm.PushStack(a)
		vm.advancePC(1)
	case RotInstruction:
		c := vm.PopStack()
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(b)
		vm.PushStack(c)
		vm.PushStack(a)
		vm.advancePC(1)
	case OverInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a)
		vm.PushStack(b)
		vm.PushStack(a)
		vm.advancePC(1)
	case NipInstruction:
		b := vm.PopStack()
		/*a :=*/ vm.PopStack()
		vm.PushStack(b)
		vm.advancePC(1)
	case TuckInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(b)
		vm.PushStack(a)
		vm.PushStack(b)
		vm.advancePC(1)
	case ToRInstruction:
		addr := vm.PopStack16()
		vm.PushReturnStack(addr)
		vm.advancePC(1)
	case FromRInstruction:
		addr := vm.PopReturnStack()
		vm.PushStack16(addr)
		vm.advancePC(1)
	case FetchRInstruction:
		addr := vm.PopReturnStack()
		vm.PushStack16(addr)
		vm.PushReturnStack(addr)
		vm.advancePC(1)
	case AddInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a + b)
		vm.advancePC(1)
	case Add16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(a + b)
		vm.advancePC(1)
	case SubInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(b - a)
		vm.advancePC(1)
	case Sub16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(b - a)
		vm.advancePC(1)
	case MulInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a * b)
		vm.advancePC(1)
	case Mul16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(a * b)
		vm.advancePC(1)
	case DivInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		if a == 0 {
			vm.setFault(FlagDivideByZero)
			return
		}
		vm.PushStack(b / a)
		vm.advancePC(1)
	case Div16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		if a == 0 {
			vm.setFault(FlagDivideByZero)
			return
		}
		vm.PushStack16(b / a)
		vm.advancePC(1)
	case ModInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		if a == 0 {
			vm.setFault(FlagDivideByZero)
			return
		}
		vm.PushStack(b % a)
		vm.advancePC(1)
	case Mod16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		if a == 0 {
			vm.setFault(FlagDivideByZero)
			return
		}
		vm.PushStack16(b % a)
		vm.advancePC(1)
	case AndInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a & b)
		vm.advancePC(1)
	case And16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(a & b)
		vm.advancePC(1)
	case OrInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a | b)
		vm.advancePC(1)
	case Or16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(a | b)
		vm.advancePC(1)
	case XorInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a ^ b)
		vm.advancePC(1)
	case Xor16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		vm.PushStack16(a ^ b)
		vm.advancePC(1)
	case NotInstruction:
		a := vm.PopStack()
		vm.PushStack(^a)
		vm.advancePC(1)
	case Not16Instruction:
		a := vm.PopStack16()
		vm.PushStack16(^a)
		vm.advancePC(1)
	case IncInstruction:
		a := vm.PopStack()
		vm.PushStack(a + 1)
		vm.advancePC(1)
	case Inc16Instruction:
		a := vm.PopStack16()
		vm.PushStack16(a + 1)
		vm.advancePC(1)
	case DecInstruction:
		a := vm.PopStack()
		vm.PushStack(a - 1)
		vm.advancePC(1)
	case Dec16Instruction:
		a := vm.PopStack16()
		vm.PushStack16(a - 1)
		vm.advancePC(1)
	case ShlInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a << b)
		vm.advancePC(1)
	case Shl16Instruction:
		b := vm.PopStack()
		a := vm.PopStack16()
		vm.PushStack16(a << b)
		vm.advancePC(1)
	case ShrInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		vm.PushStack(a >> b)
		vm.advancePC(1)
	case Shr16Instruction:
		b := vm.PopStack()
		a := vm.PopStack16()
		vm.PushStack16(a >> b)
		vm.advancePC(1)
	case JzInstruction:
		c := vm.PopStack()
		addr := vm.PopStack16()
		if c == 0 {
			vm.pc = addr
		} else {
			vm.advancePC(1)
		}
	case JnzInstruction:
		c := vm.PopStack()
		addr := vm.PopStack16()
		if c > 0 {
			vm.pc = addr
		} else {
			vm.advancePC(1)
		}
	case CallInstruction:
		addr := vm.PopStack16()
		vm.PushReturnStack(vm.pc + 1)
		vm.pc = addr
	case RetInstruction:
		vm.pc = vm.PopReturnStack()
	case EqInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		c := 0
		if a == b {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case Eq16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		c := 0
		if a == b {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case NqInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		c := 0
		if a != b {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case Nq16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		c := 0
		if a != b {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case GtInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		c := 0
		if b > a {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case Gt16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		c := 0
		if b > a {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case LtInstruction:
		b := vm.PopStack()
		a := vm.PopStack()
		c := 0
		if b < a {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case Lt16Instruction:
		b := vm.PopStack16()
		a := vm.PopStack16()
		c := 0
		if b < a {
			c = 1
		}
		vm.PushStack(byte(c))
		vm.advancePC(1)
	case PushInstruction:
		a := vm.MMIO.ReadByte(vm.pc + 1)
		vm.PushStack(a)
		vm.advancePC(2)
	case Push16Instruction:
		a, b := vm.MMIO.ReadByte(vm.pc+1), vm.MMIO.ReadByte(vm.pc+2)
		vm.PushStack16(binary.BigEndian.Uint16([]byte{a, b}))
		vm.advancePC(3)
	case StoreInstruction:
		c := vm.PopStack()
		addr := vm.PopStack16()

		vm.MMIO.WriteByte(addr, c)
		vm.advancePC(1)
	case Store16Instruction:
		b := vm.PopStack()
		a := vm.PopStack()
		addr := vm.PopStack16()
		vm.MMIO.WriteByte(addr, a)
		vm.MMIO.WriteByte(addr+1, b)
		vm.advancePC(1)
	case LoadInstruction:
		addr := vm.PopStack16()
		vm.PushStack(vm.MMIO.ReadByte(addr))
		vm.advancePC(1)
	case Load16Instruction:
		addr := vm.PopStack16()
		b := vm.MMIO.ReadByte(addr)
		a := vm.MMIO.ReadByte(addr + 1)
		vm.PushStack(b)
		vm.PushStack(a)
		vm.advancePC(1)
	}
}

func (vm *VM) Interrupt(callbackPtr uint16) error {
	vm.interruptChan <- callbackPtr
	return nil
}
