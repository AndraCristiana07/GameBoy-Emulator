package main

import rl "github.com/gen2brain/raylib-go/raylib"

// FF00 — P1/JOYP: Joypad
type Joypad struct {
	selectButtons bool // bit 5
	selectDpad    bool // bit 4
	buttons       Buttons
	prevButtons   Buttons
	lastWrite     byte
}

type Buttons struct {
	A, B, Select, Start   bool
	Right, Left, Up, Down bool
}

func (j *Joypad) write(value byte) {
	j.lastWrite = value & 0b00110000
	if value&0b00100000 == 0 {
		j.selectButtons = true
	} else {
		j.selectButtons = false
	}
	if value&0b00010000 == 0 {
		j.selectDpad = true
	} else {
		j.selectDpad = false
	}
}

func (j *Joypad) read() uint8 {
	result := uint8(0x0F)
	if j.selectButtons {
		result &^= 0b00100000

		if j.buttons.A == false {
			result |= 0b00000001
		}

		if j.buttons.B == false {
			result |= 0b00000010
		}

		if j.buttons.Select == false {
			result |= 0b00000100
		}

		if j.buttons.Start == false {
			result |= 0b00001000
		}

		if j.buttons.A {
			result &^= 0b00000001
		}
		if j.buttons.B {
			result &^= 0b00000010
		}
		if j.buttons.Select {
			result &^= 0b00000100
		}
		if j.buttons.Start {
			result &^= 0b00001000
		}
	}
	if j.selectDpad {
		result &^= 0b00010000
		if j.buttons.Right == false {
			result |= 0b00000001
		}
		if j.buttons.Left == false {
			result |= 0b00000010
		}
		if j.buttons.Up == false {
			result |= 0b00000100
		}
		if j.buttons.Down == false {
			result |= 0b00001000
		}

		if j.buttons.Right {
			result &^= 0b00000001
		}
		if j.buttons.Left {
			result &^= 0b00000010
		}
		if j.buttons.Up {
			result &^= 0b00000100
		}
		if j.buttons.Down {
			result &^= 0b00001000
		}
	}
	result = (j.lastWrite & 0b00110000) | (result & 0b11001111)
	return result
}

func (j *Joypad) UpdateJoypad() {
	//println("Update Joypad")
	newButtons := j.inputs()
	j.prevButtons = j.buttons
	j.buttons = newButtons
}

func (j *Joypad) inputs() Buttons {
	return Buttons{
		A:      rl.IsKeyDown(rl.KeyC),
		B:      rl.IsKeyDown(rl.KeyV),
		Select: rl.IsKeyDown(rl.KeySpace),
		Start:  rl.IsKeyDown(rl.KeyEnter),
		Right:  rl.IsKeyDown(rl.KeyD),
		Left:   rl.IsKeyDown(rl.KeyA),
		Up:     rl.IsKeyDown(rl.KeyW),
		Down:   rl.IsKeyDown(rl.KeyS),
	}
}
