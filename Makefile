PKG = github.com/k1LoW/sheer-heart-attack
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	SED = gsed
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	SED = sed
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: test depsdev

test:
	go test ./... -coverprofile=coverage.out -covermode=count

lint:
	golangci-lint run ./...

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

prerelease:
	ghch -w -N ${VER}
	gocredits -w .
	git add CHANGELOG.md CREDITS
	git commit -m'Bump up version number'
	git tag ${VER}

prerelease_for_tagpr:
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	goreleaser --rm-dist

.PHONY: default test
