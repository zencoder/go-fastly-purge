GO ?= godep go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: build test cover
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/go-fastly-purge
fmt:
	$(GO) fmt ./...
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	$(GO) test -v ./ -race -cover -coverprofile=$(COVERAGEDIR)/go-fastly-purge.coverprofile
cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/go-fastly-purge.coverprofile -o $(COVERAGEDIR)/go-fastly-purge.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	$(GO) clean
	rm -f bin/go-fastly-purge
	rm -rf coverage/
godep-save:
	godep save ./...
