package devices

import (
	"fmt"

	"github.com/alisdairrankine/frienvironment/vm"
)

const TerminalDeviceType = 0x05

/**
Terminal Device

Reg  Addr   Name              R/W   Description
---  ----   ----              ---   -----------
0    0x00   device_type       R     0x05 = terminal
1    0x01   data_addr_high    W
2    0x02   data_addr_low     W
3    0x03   data_length       W
4    0x04   write_trigger     W

**/

type Terminal struct {
	vm         *vm.VM
	DataAddr   uint16
	DataLength byte
}

func NewTerminal(vm *vm.VM) *Terminal {
	return &Terminal{
		vm: vm,
	}
}

func (t *Terminal) Write(addr uint16, data byte) {
	if addr == 1 {
		fmt.Print(string(data))
	}

	switch addr & 0x000F {

	case 0x01:
		t.DataAddr = (uint16(data) << 8) | (t.DataAddr & 0x00FF)
	case 0x02:
		t.DataAddr = (uint16(data)) | (t.DataAddr & 0xFF00)
	case 0x03:
		t.DataLength = data
		fmt.Printf("data length : %d\n", data)
	case 0x04:
		data := t.vm.MMIO.ReadData(t.DataAddr, int(t.DataLength))
		fmt.Println("OUT: ", string(data))
	}
}

func (Terminal) Read(addr uint16) byte {
	if addr == 0 {
		return 0x05
	}
	return 0
}
