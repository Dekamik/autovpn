BIN=autovpn
OUT=bin/${BIN}

all: deps build

deps:
	go mod tidy

build:
	go build -o ${OUT} -ldflags "-X main.version=$(shell cat VERSION) -X main.build=$(shell date +%y%m%d%H%M)" cmd/main.go

install:
	cp ${OUT} /usr/local/bin

uninstall:
	rm /usr/local/bin/${BIN}

clean:
	go clean
	rm ${OUT}
