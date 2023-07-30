server:
	go run cmd/main.go

test:
	go test -v -cover -short ./...

.PHONY: server test
