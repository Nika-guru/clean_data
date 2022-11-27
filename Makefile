run:
	go run cmd/main/*.go

build:
	go build cmd/main/main.go

linux:
	env GOOS=linux GOARCH=amd64 go build cmd/main/main.go

clean:
	rm -rf vendor