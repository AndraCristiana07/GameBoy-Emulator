package main

import (
	"fmt"
	"os"
)

type Cartridge struct {
	bootROM         []byte   //responsible for the boot-up animation played before control is handed over to the cartridge’s ROM
	entryPoint      []byte   // 0100-0103
	nintendoLogo    []byte   // 0104-0133
	title           string   // 0134-013E
	manufacturer    []byte   // 013F-0142
	CCBFlag         uint8    // 0143
	newLicenseeCode []byte   // 0144-0145
	SGBFlag         uint8    // 0146
	cartridgeType   uint8    // 0147
	ROMSize         uint8    // 0148
	RAMSize         uint8    // 0149
	destCode        uint8    // 014A
	oldLicenseeCode uint8    // 014B -> if 33 => new lic
	maskROMVer      uint8    // 014C
	checksum        uint8    //014D
	globalChecksum  [2]uint8 //014E - 014F
	ROMdata         []byte   //all data
}

func LoadCartridge(filename string) (*Cartridge, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading cartridge file: %s\n", err)

	}
	cartridge := Cartridge{
		bootROM:         data[0x0000:0x0100],
		entryPoint:      data[0x0100:0x0104],
		nintendoLogo:    data[0x0104:0x0134],
		title:           string(data[0x134:0x013F]),
		manufacturer:    data[0x013F:0x0143],
		CCBFlag:         data[0x0143],
		newLicenseeCode: data[0x0144:0x0146],
		SGBFlag:         data[0x0146],
		cartridgeType:   data[0x0147],
		ROMSize:         data[0x0148],
		RAMSize:         data[0x0149],
		destCode:        data[0x014A],
		oldLicenseeCode: data[0x014B],
		maskROMVer:      data[0x014C],
		checksum:        data[0x014D],
		globalChecksum:  [2]uint8{data[0x014E], data[0x014F]},
		ROMdata:         data,
	}
	return &cartridge, nil
}

func getROMSize(size uint8) int {
	return 32 * 1024 << size
}

func getRAMSize(size uint8) int {
	sizes := map[uint8]int{
		0x00: 0,
		0x01: 2 * 1024, //unused
		0x02: 8 * 1024,
		0x03: 32 * 1024,
		0x04: 128 * 1024,
		0x05: 64 * 1024,
	}
	return sizes[size]
}
func (cartridge *Cartridge) printInfo() {
	fmt.Printf("Cartridge info:\n")
	fmt.Printf("Title: %s\n", cartridge.title)
	fmt.Printf("CCB Flag: %d\n", cartridge.CCBFlag)
	fmt.Printf("Cartridge type: %d\n", cartridge.cartridgeType)
	fmt.Printf("ROM Size: %d KB\n", getROMSize(cartridge.RAMSize)/1024)
	fmt.Printf("RAM size: %d KB\n", getRAMSize(cartridge.RAMSize)/1024)
	fmt.Printf("Destination code: 0x%d\n", cartridge.destCode)
	fmt.Printf("New licensee code: 0x%d\n", cartridge.newLicenseeCode)
	fmt.Printf("SGB Flag: %d\n", cartridge.SGBFlag)
	fmt.Printf("Checksum: 0x%d\n", cartridge.checksum)
	fmt.Printf("Global Checksum: 0x%d\n", cartridge.globalChecksum)
}

//func main() {
//	cartridge, err := loadCartridge("roms/The Legend of Zelda - Links Awakening (US - EU).gb")
//	if err != nil {
//		fmt.Printf("Error loading cartridge: %s\n", err)
//		return
//	}
//	cartridge.printInfo()

//}
