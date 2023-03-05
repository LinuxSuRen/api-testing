build:
	mkdir -p bin
	go build -o bin/atest cmd/*.go

copy: build
	cp bin/atest /usr/local/bin/
test:
	go test ./...
