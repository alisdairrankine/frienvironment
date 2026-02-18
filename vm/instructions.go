package vm

const (
	YieldInstruction   byte = 0x00
	HaltInstruction    byte = 0x01
	DupInstruction     byte = 0x02
	DropInstruction    byte = 0x03
	SwapInstruction    byte = 0x04
	RotInstruction     byte = 0x05
	OverInstruction    byte = 0x06
	NipInstruction     byte = 0x07
	TuckInstruction    byte = 0x08
	ToRInstruction     byte = 0x09
	FromRInstruction   byte = 0x0A
	FetchRInstruction  byte = 0x0B
	AddInstruction     byte = 0x10
	Add16Instruction   byte = 0x11
	SubInstruction     byte = 0x12
	Sub16Instruction   byte = 0x13
	MulInstruction     byte = 0x14
	Mul16Instruction   byte = 0x15
	DivInstruction     byte = 0x16
	Div16Instruction   byte = 0x17
	ModInstruction     byte = 0x18
	Mod16Instruction   byte = 0x19
	AndInstruction     byte = 0x20
	And16Instruction   byte = 0x21
	OrInstruction      byte = 0x22
	Or16Instruction    byte = 0x23
	XorInstruction     byte = 0x24
	Xor16Instruction   byte = 0x25
	NotInstruction     byte = 0x26
	Not16Instruction   byte = 0x27
	IncInstruction     byte = 0x28
	Inc16Instruction   byte = 0x29
	DecInstruction     byte = 0x2A
	Dec16Instruction   byte = 0x2B
	ShlInstruction     byte = 0x2C
	Shl16Instruction   byte = 0x2D
	ShrInstruction     byte = 0x2E
	Shr16Instruction   byte = 0x2F
	JzInstruction      byte = 0x40
	JnzInstruction     byte = 0x41
	CallInstruction    byte = 0x42
	RetInstruction     byte = 0x43
	EqInstruction      byte = 0x44
	Eq16Instruction    byte = 0x45
	NqInstruction      byte = 0x46
	Nq16Instruction    byte = 0x47
	GtInstruction      byte = 0x48
	Gt16Instruction    byte = 0x49
	LtInstruction      byte = 0x4A
	Lt16Instruction    byte = 0x4B
	PushInstruction    byte = 0x50
	Push16Instruction  byte = 0x51
	StoreInstruction   byte = 0x52
	Store16Instruction byte = 0x53
	LoadInstruction    byte = 0x54
	Load16Instruction  byte = 0x55
)

var InstrName = map[byte]string{
	YieldInstruction:   "Yield",
	HaltInstruction:    "Halt",
	DupInstruction:     "Dup",
	DropInstruction:    "Drop",
	SwapInstruction:    "Swap",
	RotInstruction:     "Rot",
	OverInstruction:    "Over",
	NipInstruction:     "Nip",
	TuckInstruction:    "Tuck",
	ToRInstruction:     "ToR",
	FromRInstruction:   "FromR",
	FetchRInstruction:  "FetchR",
	AddInstruction:     "Add",
	Add16Instruction:   "Add16",
	SubInstruction:     "Sub",
	Sub16Instruction:   "Sub16",
	MulInstruction:     "Mul",
	Mul16Instruction:   "Mul16",
	DivInstruction:     "Div",
	Div16Instruction:   "Div16",
	ModInstruction:     "Mod",
	Mod16Instruction:   "Mod16",
	AndInstruction:     "And",
	And16Instruction:   "And16",
	OrInstruction:      "Or",
	Or16Instruction:    "Or16",
	XorInstruction:     "Xor",
	Xor16Instruction:   "Xor16",
	NotInstruction:     "Not",
	Not16Instruction:   "Not16",
	IncInstruction:     "Inc",
	Inc16Instruction:   "Inc16",
	DecInstruction:     "Dec",
	Dec16Instruction:   "Dec16",
	ShlInstruction:     "Shl",
	Shl16Instruction:   "Shl16",
	ShrInstruction:     "Shr",
	Shr16Instruction:   "Shr16",
	JzInstruction:      "Jz",
	JnzInstruction:     "Jnz",
	CallInstruction:    "Call",
	RetInstruction:     "Ret",
	EqInstruction:      "Eq",
	Eq16Instruction:    "Eq16",
	NqInstruction:      "Nq",
	Nq16Instruction:    "Nq16",
	GtInstruction:      "Gt",
	Gt16Instruction:    "Gt16",
	LtInstruction:      "Lt",
	Lt16Instruction:    "Lt16",
	PushInstruction:    "Push",
	Push16Instruction:  "Push16",
	StoreInstruction:   "Store",
	Store16Instruction: "Store16",
	LoadInstruction:    "Load",
	Load16Instruction:  "Load16",
}
