.PHONY: all build docker

default: all

all: build plugins

build:
	echo "TODO: -X DOES NOT WORK"
	go build -ldflags="-s -w -X 'cmd/gokamux/cmd.Version=${APP_VERSION}' -X 'cmd/gokamux/cmd.Build=${APP_BUILD}'" -o bin/gokamux cmd/gokamux/main.go

clean:
	rm -f bin/gokamux
	rm -fR bin/plugins

PLUGINSDIR = modules/plugins
RESULTSDIR = bin/plugins

plugins:
	@for p in $(shell ls ${PLUGINSDIR}); do echo "Building plugin $${p}"; go build -buildmode=plugin -o ${RESULTSDIR}/$${p}.so ${PLUGINSDIR}/$${p}/*.go $${f}; done

docker:
	docker build -t gokamux . -f docker/Dockerfile