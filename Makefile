.PHONY: all build docker

PACKAGE_BASE=github.com/datoga/gokamux/cmd/gokamux/cmd
PLUGINSDIR = modules/plugins
RESULTSDIR = bin/plugins

default: all

all: core plugins

core:
	go build -ldflags="-s -w -X '${PACKAGE_BASE}.Version=${APP_VERSION}' -X '${PACKAGE_BASE}.Build=${APP_BUILD}'" -o bin/gokamux cmd/gokamux/main.go

plugins:
	@for p in $(shell ls ${PLUGINSDIR}); do echo "Building plugin $${p}"; go build -buildmode=plugin -o ${RESULTSDIR}/$${p}.so ${PLUGINSDIR}/$${p}/*.go $${f}; done

docker:
	docker build -t gokamux . -f docker/Dockerfile

docker-with-plugins:
	docker build --build-arg PLUGINS=1 -t gokamux . -f docker/Dockerfile

clean:
	rm -f bin/gokamux
	rm -fR bin/plugins