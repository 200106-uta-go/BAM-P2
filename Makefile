tests:
	go test ./...

builds:
	go build -o build/bin/rproxy cmd/reverseproxy/rproxy.go
	go build -o build/bin/authServer cmd/jServer/jServer.go

run-local:
	go run cmd/reverseproxy/rproxy.go &
	go run cmd/jServer/jServer.go

run-local-docker:
	CGO_ENABLED=0 GOOS=linux go build -o cmd/jServer/jServer-scratch cmd/jServer/jServer.go
	docker build -t jserver cmd/jServer/.
	docker run --rm -p 8080:8080 jserver
	rm cmd/jServer/jServer-scratch
