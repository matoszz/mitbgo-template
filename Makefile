default: all

all: fmt test build

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: 
	$(info ******************** running tests ********************)
	go test -v ./...

generate:
	$(info ******************** generating ent schema ********************)
	go mod tidy
	go generate ./..
	