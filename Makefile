GO ?= godep go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: build test cover
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/ddbsync
fmt:
	$(GO) fmt ./...
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	$(GO) test -v ./ -race -cover -coverprofile=$(COVERAGEDIR)/ddbsync.coverprofile
cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/ddbsync.coverprofile -o $(COVERAGEDIR)/ddbsync.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	$(GO) clean
	rm -f bin/ddbsync
	rm -rf coverage/
gen-mocks:
	mockery -name AWSDynamoer
	mockery -name DBer
godep-save:
	godep save ./...
