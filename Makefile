GO15VENDOREXPERIMENT := 1
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: build test cover
install-deps:
	glide install
build:
	if [ ! -d bin ]; then mkdir bin; fi
	go build -v -o bin/go-fastly-purge
fmt:
	go fmt ./...
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./ -race -cover -coverprofile=$(COVERAGEDIR)/go-fastly-purge.coverprofile
cover:
	go tool cover -html=$(COVERAGEDIR)/go-fastly-purge.coverprofile -o $(COVERAGEDIR)/go-fastly-purge.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	go clean
	rm -f bin/go-fastly-purge
	rm -rf coverage/
