COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

ifdef VERBOSE
V = -v
else
.SILENT:
endif

all: build test cover

install-deps:
	glide install
build:
	mkdir -p bin
	go build $(V) -o bin/go-fastly-purge
fmt:
	go fmt ./...
test:
	mkdir -p coverage
	go test $(V) ./ -race -cover -coverprofile=$(COVERAGEDIR)/go-fastly-purge.coverprofile
cover:
	go tool cover -html=$(COVERAGEDIR)/go-fastly-purge.coverprofile -o $(COVERAGEDIR)/go-fastly-purge.html
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	go clean
	rm -f bin/go-fastly-purge
	rm -rf coverage/
