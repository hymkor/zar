ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    NUL=nul
else
    SET=export
    NUL=/dev/null
endif

NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT:=-ldflags "-s -w -X main.version=$(VERSION)"
EXE:=$(shell go env GOEXE)

ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
else
    SET=export
endif

all:
	go fmt ./...
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_dist:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXE)

dist:
	go fmt ./...
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

manifest:
	make-scoop-manifest *-windows-*.zip > $(NAME).json

release:
	pwsh -Command "latest-notes.ps1" | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

.PHONY: all dist _dist release manifest
