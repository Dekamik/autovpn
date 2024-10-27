BIN=autovpn
OUT=bin/${BIN}

all: deps build

deps:
	go mod tidy

build:
	go build -o ${OUT} -ldflags "-X main.version=$(shell cat VERSION) -X main.build=$(shell date +%y%m%d%H%M)" cmd/main.go

install:
	mkdir -p /usr/local/bin
	cp ${OUT} /usr/local/bin/${BIN}

uninstall:
	rm /usr/local/bin/${BIN}

clean:
	go clean
	rm ${OUT}
