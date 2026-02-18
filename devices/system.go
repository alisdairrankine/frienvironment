package devices

import (
	"errors"

	"github.com/alisdairrankine/frienvironment/vm"
)

type System struct {
	vms map[uint8]*vm.VM

	next uint8

	dead []uint8
}

func NewSystem() *System {
	return &System{
		vms: make(map[uint8]*vm.VM),
	}
}

type SpawnOption func(*vm.VM)

func (s *System) Spawn(program []byte, spawnOptions ...SpawnOption) (id uint8, err error) {

	if len(s.dead) > 0 {
		id = s.dead[0]
		s.dead = s.dead[1:]
	} else {
		if s.next < 255 {
			id = s.next
			s.next++
		} else {
			return 0, errors.New("no capacity")
		}
	}
	machine := vm.New()
	machine.LoadProgram(program)
	for _, opt := range spawnOptions {
		opt(machine)
	}

	machine.Run()
	return id, nil
}

func (s *System) Kill(vmID uint8) error {
	if vm, ok := s.vms[vmID]; ok {
		vm.Stop()
		return nil
	}
	return errors.New("vm not found")
}

func (s *System) GetVM(vmID uint8) (*vm.VM, error) {
	if vm, ok := s.vms[vmID]; ok {
		return vm, nil
	}
	return nil, errors.New("vm not found")
}
