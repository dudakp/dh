install: build
	go install

build: test
	go build -o bin/dh.exe main.go

test:
	go test -v ./...
