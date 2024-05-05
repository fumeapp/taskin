ray:
	@go get github.com/octoper/go-ray

update:
	go get -u ./... && go mod tidy

lint:
	@golangci-lint run


test:
	@go test -v ./...
