.PHONY: test
test: 
	go clean -testcache && go test . --tags=all -v

.PHONY: build
build:
	go get
	go build -o gowc

.PHONY: install
install: build
	go install