package frienvironment

import "sync"

type Buffer struct {
	bytes []byte
	sync.Mutex
}

func (b *Buffer) Read() []byte {
	b.Lock()
	defer b.Unlock()
	cpy := []byte{}
	copy(cpy, b.bytes)
	return cpy
}

func (b *Buffer) AddString(s string) {
	b.Lock()
	defer b.Unlock()
	b.bytes = append(b.bytes, []byte(s)...)
}

func (b *Buffer) Clear() {
	b.Lock()
	defer b.Unlock()
	b.bytes = []byte{}
}

func (b *Buffer) Send() {
	b.Lock()
	defer b.Unlock()
	b.bytes = []byte{}
}
