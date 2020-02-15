tests:
	go test ./...

builds:
	go build -o build/bin/rproxy cmd/reverseproxy/rproxy.go
	go build -o build/bin/authServer cmd/jServer/jServer.go

run-local:
	go run cmd/reverseproxy/rproxy.go &
	go run cmd/jServer/jServer.go

docker-build:
	docker build -t jserver cmd/jServer/.

docker-run:
	docker run --rm --env-file .env -p 4444:8080 jserver
