all: clean test build

.PHONY: clean
clean:
	rm -f chirp_*

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --enable-all

.PHONY: test
test: api/v1/clipboard.pb.go
	go test -cover ./...

.PHONY: build
build: api/v1/clipboard.pb.go
	gox -arch="amd64" -os="darwin linux" ./...

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:. --go_opt=module=github.com/oclaussen/chirp $<
