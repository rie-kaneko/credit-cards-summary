lint:
	golangci-lint run -v

fmt:
	go fmt ./...

test:
	go clean -testcache && go test -v -cover ./...