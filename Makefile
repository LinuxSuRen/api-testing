build:
	mkdir -p bin
	rm -rf bin/atest
	go build -o bin/atest main.go

copy: build
	sudo cp bin/atest /usr/local/bin/
test:
	go test ./...
