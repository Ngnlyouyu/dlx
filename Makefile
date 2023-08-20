all: clean build

build:
	go build -o ./build/bin/dlx ./dlx.go

clean:
	rm -rf ./build

check:
	golangci-lint run -v
