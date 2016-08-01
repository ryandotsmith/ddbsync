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
	go build -v -o bin/ddbsync
fmt:
	go fmt ./...
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./ -race -cover -coverprofile=$(COVERAGEDIR)/ddbsync.coverprofile
cover:
	go tool cover -html=$(COVERAGEDIR)/ddbsync.coverprofile -o $(COVERAGEDIR)/ddbsync.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	go clean
	rm -f bin/ddbsync
	rm -rf coverage/
gen-mocks:
	mockery -name AWSDynamoer
	mockery -name DBer
	mockery -name LockServicer
