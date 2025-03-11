
run:
	go run main.go

unittest:
	go test -v

testcoverage:
	go test -coverprofile cover.out && cat cover.out
