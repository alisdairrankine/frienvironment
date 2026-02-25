package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/alisdairrankine/frienvironment/assembler"
	"github.com/alisdairrankine/frienvironment/devices"
	"github.com/alisdairrankine/frienvironment/vm"
)

func main() {

	sender, err := LoadProgram("progs/sender.ca")
	if err != nil {
		log.Fatal(err)
	}

	receiver, err := LoadProgram("progs/receiver.ca")
	if err != nil {
		log.Fatal(err)
	}

	sys := devices.NewSystem()
	netSwitch := devices.NewSwitch()
	sys.Spawn(receiver, func(vm *vm.VM) {
		netSwitch.Attach(0, vm)
		vm.RegisterDevice(1, devices.NewTerminal(vm))
		// vm.Debug = true
	})
	sys.Spawn(sender, func(vm *vm.VM) {
		netSwitch.Attach(0, vm)

	})

	time.Sleep(100 * time.Second)
}

func LoadProgram(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return assembler.Assemble(string(data)), nil
}
