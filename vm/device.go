package vm

type Device interface {
	Write(addr uint16, data byte)
	Read(addr uint16) byte
}
