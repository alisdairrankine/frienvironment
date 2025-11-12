package frienvironment

type InterruptType int

const (
	InterruptTypeNULL            InterruptType = 0
	InterruptTypeKeyPressed      InterruptType = 1
	InterruptTypeMouseDown       InterruptType = 2
	InterruptTypeWindowMouseDown InterruptType = 3
)
