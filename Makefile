.PHONY: all build

default: all

all: build plugins

build:
	go build -o bin/gokamux cmd/gokamux/main.go

clean:
	rm -f bin/gokamux
	rm -fR bin/plugins

PLUGINSDIR = modules/plugins
RESULTSDIR = bin/plugins

plugins:
	@for p in $(shell ls ${PLUGINSDIR}); do echo "Building plugin $${p}"; go build -buildmode=plugin -o ${RESULTSDIR}/$${p}.so ${PLUGINSDIR}/$${p}/*.go $${f}; done