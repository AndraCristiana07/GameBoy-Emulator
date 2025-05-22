package main

// FF00 — P1/JOYP: Joypad
type Joypad struct {
	selectButtons bool // bit 5
	selectDpad    bool // bit 4
	buttons       []bool
}

// func (j * Joypad) write(value byte){

// }
//func (j *Joypad) read()
