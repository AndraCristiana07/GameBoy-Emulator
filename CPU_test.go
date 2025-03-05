package main

import (
	"fmt"
	"testing"
)

func TestSetBC(t *testing.T) {
	register := Registers{}
	register.setBC(0x0001)
	fmt.Println(register.getBC())
	if register.getBC() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getBC())
	}
}

func TestSetDE(t *testing.T) {
	register := Registers{}
	register.setDE(0x0001)
	fmt.Println(register.getDE())
	if register.getDE() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getDE())
	}
}

func TestSetHL(t *testing.T) {
	register := Registers{}
	register.setHL(0x0001)
	fmt.Println(register.getHL())
	if register.getHL() != 0x0001 {
		t.Error("Expected 0x0001, got ", register.getDE())
	}
}
