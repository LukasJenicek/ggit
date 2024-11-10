build:
	go build -o /home/lj/.local/bin/ggit cmd/main.go

test:
	go test -race -shuffle=on -v ./...

vendor:
	go mod tidy && go mod vendor && go mod tidy