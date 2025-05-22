
run:
	LOGXI=* go run main.go gpu.go cpu.go parseRom.go timer.go

release:
	go run main.go gpu.go cpu.go parseRom.go timer.go


unittest:
	go test -v

testcoverage:
	go test -coverprofile cover.out &> /dev/null
	cat cover.out | grep .0$$ > report.cover && rm cover.out
	cat report.cover
