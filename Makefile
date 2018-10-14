ifdef update
  u=-u
endif

deps:
	go get -d -v -t ./...

test-deps:
	go get -d -v ./...

devel-deps: test-deps
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

test: test-deps
	go test ./...

lint: devel-deps
	go vet ./...
	golint -set_exit_status ./...

cover: devel-deps
	goveralls

.PHONY: deps test-deps devel-deps test lint cover
