.PHONY: test test-complete

all: test-complete

test:
	go test -short

test-complete:
	go test
