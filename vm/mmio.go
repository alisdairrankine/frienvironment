package vm

import "fmt"

type MMIO struct {
	data [65536]byte

	devices [16]Device
}

func (m *MMIO) WriteByte(addr uint16, data byte) {
	if (addr & 0xFF00) == 0x0300 {
		m.writeToDevice(addr, data)
		return
	}
	m.data[addr] = data
}

func (m *MMIO) ReadByte(addr uint16) byte {
	if (addr & 0xFF00) == 0x0300 {
		return m.readFromDevice(addr)
	}
	return m.data[addr]
}

func (m *MMIO) ReadData(addr uint16, length int) []byte {
	return m.data[addr : int(addr)+length]
}

func (m *MMIO) WriteData(addr uint16, data []byte) {
	fmt.Printf("WRITE DATA %X\n", addr)
	copy(m.data[addr:int(addr)+len(data)], data[:])
}

func (m *MMIO) writeToDevice(addr uint16, data byte) {
	deviceNumber := (addr & 0x00F0) >> 4
	if device := m.devices[deviceNumber]; device != nil {
		deviceAddr := uint16(addr & 0x000F)
		device.Write(deviceAddr, data)
	}
}

func (m *MMIO) readFromDevice(addr uint16) byte {
	deviceNumber := (addr & 0x00F0) >> 4
	if device := m.devices[deviceNumber]; device != nil {
		deviceAddr := uint16(addr & 0x000F)
		return device.Read(deviceAddr)
	}
	return 0
}
