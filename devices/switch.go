package devices

import (
	"errors"
	"fmt"

	"github.com/alisdairrankine/frienvironment/vm"
)

const SwitchDeviceType = 0x00

/**
Switch Device

Reg  Addr   Name              R/W   Description
---  ----   ----              ---   -----------
0    0x00   device_type       R     0x01 = switch
1    0x01   vm_id             R     this VM's assigned ID
2    0x02   dest_id           W     destination VM ID
3    0x03   sender_id         R     VM ID of last received message
4    0x04   send_signal       W     signal ID to send
5    0x05   recv_signal       R     signal ID of received message
6    0x06   msg_len           W     outgoing payload length
7    0x07   msg_addr_high     W     high byte of outgoing payload address
8    0x08   msg_addr_low      W     low byte of outgoing payload address
9    0x09   recv_addr_high    W     high byte of receive buffer address
10   0x0A   recv_addr_low     W     low byte of receive buffer address
11   0x0B   callback_high     W     high byte of receive callback
12   0x0C   callback_low      W     low byte of receive callback
13   0x0D   send_trigger      W     write any value to trigger send
14   0x0E   recv_length		  R		incoming payload length
15   0x0F   reserved
**/

type Switch struct {
	ports []*Port //todo: map[byte]*port
}

func NewSwitch() *Switch {
	return &Switch{}
}

func (s *Switch) Attach(deviceNum int, vm *vm.VM) {
	id := len(s.ports)
	port := &Port{
		deviceID: deviceNum,
		s:        s,
		vm:       vm,
		port:     byte(id),
	}
	vm.RegisterDevice(deviceNum, port)
	s.ports = append(s.ports, port)
}

func (s *Switch) Send(sourceAddr, destinationAddr, signalID byte, data []byte) error {
	port := s.ports[destinationAddr]
	if port == nil {
		return errors.New("host not found")
	}
	port.Receive(sourceAddr, signalID, data)
	return nil
}

type Port struct {
	s  *Switch
	vm *vm.VM

	deviceID int

	port            byte
	destinationPort byte
	senderPort      byte
	sendSignal      byte
	receiveSignal   byte
	messageLength   byte
	recvLength      byte
	messageAddr     uint16
	recvAddr        uint16
	callbackAddr    uint16
}

func (p *Port) Receive(sourceAddr, signalID byte, data []byte) {
	p.senderPort = sourceAddr
	p.receiveSignal = signalID
	p.vm.MMIO.WriteData(p.recvAddr, data)
	p.recvLength = byte(len(data))
	p.vm.Interrupt(uint16((byte(p.deviceID)<<4)+0x0B) | 0x0300)
}

func (p *Port) Write(addr uint16, data byte) {
	switch addr & 0x000F {
	case 0x02:
		p.destinationPort = data
	case 0x04:
		p.sendSignal = data
	case 0x06:
		p.messageLength = data
	case 0x07:
		p.messageAddr = (uint16(data) << 8) | (p.messageAddr & 0x00FF)
	case 0x08:
		p.messageAddr = (uint16(data)) | (p.messageAddr & 0xFF00)
	case 0x09:
		p.recvAddr = (uint16(data) << 8) | (p.recvAddr & 0x00FF)
	case 0x0A:
		p.recvAddr = (uint16(data)) | (p.recvAddr & 0xFF00)
	case 0x0B:
		p.callbackAddr = (uint16(data) << 8) | (p.callbackAddr & 0x00FF)
		fmt.Printf("callback: %X\n", p.callbackAddr)
	case 0x0C:
		p.callbackAddr = (uint16(data)) | (p.callbackAddr & 0xFF00)
		fmt.Printf("callback: %X\n", p.callbackAddr)
	case 0x0D:
		data := make([]byte, p.messageLength)
		p.vm.MMIO.ReadData(p.messageAddr, data)
		p.s.Send(p.port, p.destinationPort, p.sendSignal, data)
	}
}

func (p *Port) Read(addr uint16) byte {
	switch addr & 0x000F {
	case 0x00:
		return SwitchDeviceType
	case 0x01:
		return p.port
	case 0x02:
		return p.destinationPort
	case 0x03:
		return p.senderPort
	case 0x04:
		return p.sendSignal
	case 0x05:
		return p.receiveSignal
	case 0x06:
		return p.messageLength
	case 0x07:
		return byte(p.messageAddr >> 8)
	case 0x08:
		return byte(p.messageAddr & 0x00FF)
	case 0x09:
		return byte(p.recvAddr >> 8)
	case 0x0A:
		return byte(p.recvAddr & 0x00FF)
	case 0x0B:
		return byte(p.callbackAddr >> 8)
	case 0x0C:
		return byte(p.callbackAddr & 0x00FF)
	case 0x0E:
		return byte(p.recvLength)
	}
	return 0
}
