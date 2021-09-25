.PHONY: all build

default: all

all: build

build:
	go build -o bin/gokamux cmd/gokamux/main.go

clean:
	rm bin/gokamux