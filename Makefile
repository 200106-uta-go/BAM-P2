tests:
	go test ./...

builds:
	go build -o build/bin/rproxy cmd/reverseproxy/rproxy.go

build-scratch:
	CGO_ENABLED=0 GOOS=linux go build -o build/bin/rproxy-scratch cmd/reverseproxy/rproxy.go 

run-local:
	go run cmd/reverseproxy/rproxy.go
