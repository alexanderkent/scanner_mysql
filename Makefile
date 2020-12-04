
build:
	go build -o bin/scanner .

build-all:
	echo "Compiling for other OS/Platforms"
	GOOS=freebsd GOARCH=386 go build -o bin/scanner-freebsd-386 .
	GOOS=linux GOARCH=386 go build -o bin/scanner-linux-386 .
	GOOS=windows GOARCH=386 go build -o bin/scanner-windows-386 .

run:
	go run .

clean:
	rm bin/*

all: build build-all