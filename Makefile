
ifeq ($(ROM),)
$(error No ROM given)
endif

run:
	LOGXI=* go run main.go gpu.go cpu.go parseRom.go timer.go cpu_common.go cpu_ops.go joypad.go $(ROM)

release:
	go run main.go gpu.go cpu.go parseRom.go timer.go cpu_common.go cpu_ops.go joypad.go $(ROM)


unittest:
	go test -v

testcoverage:
	go test -coverprofile cover.out &> /dev/null
	cat cover.out | grep .0$$ > report.cover && rm cover.out
	cat report.cover
