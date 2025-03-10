GOLANGCI-VERSION = v1.64.6

build:
	go build -o ./bin/ggit cmd/main.go
	go build -o ./bin/inflate cmd/utils/inflate/main.go

install:
	go build -o /home/lj/.local/bin/ggit cmd/main.go

test:
	go test -race -shuffle=on -v ./...

vendor:
	go mod tidy && go mod vendor && go mod tidy

lint:
	golangci-lint run -c .golangci.yml --fix

install-golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(HOME)/.local/bin" $(GOLANGCI-VERSION)
