package main

type Timer struct {
	div             uint16 //divider register - FF04
	oldStoredDivBit byte
	TIMA            byte //timer counter - FF05
	TMA             byte // timer modulo FF06
	TAC             byte // timer control - FF07
}

func (t *Timer) Read(address uint16) byte {
	switch address {
	case 0xFF04:
		return byte(t.div >> 8) //upper 8 bits
	case 0xFF05:
		return t.TIMA
	case 0xFF06:
		return t.TMA
	case 0xFF07:
		return t.TAC
	}
	return 0xFF

}

func (t *Timer) Write(address uint16, value byte) {
	switch address {
	case 0xFF04:
		t.div = 0 //Writing any value to this register resets it to $00
	case 0xFF05:
		t.TIMA = value
	case 0xFF06: // TODO: if set to FF
		t.TMA = value
	case 0xFF07:
		t.TAC = value & 0111
	}
}

func (t *Timer) Update(tCycles int, cpu *CPU) {

	t.div += uint16(tCycles)
	cpu.Memory[0xFF04] = byte(t.div >> 8)
	//if t.div >= 256 {
	//	t.div -= 256
	//}
	if t.TAC&(1<<2) == 0 { //timer disabled
		return
	}
	var rate uint16
	// Bits 1-0 : Clock Select
	switch t.TAC & 0b0011 {
	case 0b00:
		rate = 1024
	case 0b01:
		rate = 16
	case 0b10:
		rate = 64
	case 0b11:
		rate = 256
	}

	// bit that controls the timer inc
	bitPos := map[uint16]uint8{
		1024: 9,
		16:   3,
		64:   6,
		256:  8,
	}[rate]

	if t.TAC&(1<<2) != 0 { //timer enabled
		currDIVBit := (t.div >> bitPos) & 1
		prevBit := (t.oldStoredDivBit >> bitPos) & 1
		if currDIVBit == 1 && prevBit == 0 { //"rising edge" (meaning going from 0 to 1)
			t.TIMA++
			if t.TIMA == 0xFF { //overflow
				t.TIMA = t.TMA               //reload
				cpu.Memory[0xFF0F] |= 1 << 2 //bit 2 of the Interrupt Flag Register at $FF0F is set
			}
		}
		//t.oldStoredDivBit = t.div
		t.oldStoredDivBit = byte(currDIVBit)
	}

	//div -ncremented every 256 tCycles
	//t.div += uint16(tCycles)
	//if t.div >= 256 {
	//	t.div -= 256
	//}
	//if t.TAC&(1<<2) == 0 { //timer disabled
	//	return
	//}
	//var rate int
	//// Bits 1-0 : Clock Select
	//switch t.TAC & 0b0011 {
	//case 0b00:
	//	rate = 1024
	//case 0b01:
	//	rate = 16
	//case 0b10:
	//	rate = 64
	//case 0b11:
	//	rate = 256
	//}
	//if t.TAC&(1<<2) != 0 { //timer enabled
	//	currDIVBit := t.div & uint16(rate)
	//
	//	if currDIVBit == 1 && t.oldStoredDivBit == 0 { //"rising edge" (meaning going from 0 to 1)
	//		t.TIMA++
	//		if t.TIMA == 0 { //overflow
	//			t.TIMA = t.TMA               //reload
	//			cpu.Memory[0xFF0F] |= 1 << 2 //bit 2 of the Interrupt Flag Register at $FF0F is set
	//		}
	//	}
	//	t.oldStoredDivBit = byte(currDIVBit)
	//}

}
