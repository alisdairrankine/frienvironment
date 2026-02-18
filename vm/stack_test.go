package vm

import "testing"

func TestStackOverflow(t *testing.T) {
	m := New()
	m.Run()
}
