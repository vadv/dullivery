SOURCEDIR=./src
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)

GOPATH ?= /usr/local/go
GOPATH := ${CURDIR}:${GOPATH}
export GOPATH

all: ./bin/dullivery-web ./bin/dullivery

./bin/dullivery-web: $(SOURCES)
	go build -o ./bin/dullivery-web -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/cmd/dullivery-web/main.go

./bin/dullivery: $(SOURCES)
	go build -o ./bin/dullivery -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/cmd/dullivery/main.go

generate_proto: clean
	$(MAKE) -C src/api

tar: clean
	mkdir -p rpm/SOURCES
	tar --transform='s,^\.,dullivery-$(VERSION),'\
		-czf rpm/SOURCES/dullivery-$(VERSION).tar.gz .\
		--exclude=rpm/SOURCES

docker: submodule_check tar
	cp -a $(CURDIR)/rpm /build
	cp -a $(CURDIR)/rpm/SPECS/dullivery.spec /build/SPECS/dullivery-$(VERSION).spec
	sed -i 's|%define version unknown|%define version $(VERSION)|g' /build/SPECS/dullivery-$(VERSION).spec
	chown -R root:root /build
	rpmbuild -ba --define '_topdir /build'\
		/build/SPECS/dullivery-$(VERSION).spec

clean:
	rm -f rpm-tmp.*
	rm -rf pkg

test: generate_proto
	go test -x -v server/posix dsl

.DEFAULT_GOAL: all

include Makefile.git
