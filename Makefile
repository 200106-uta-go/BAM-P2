tests:
	go test ./...

builds:
	go build -o build/bin/rproxy cmd/reverseproxy/rproxy.go
	go build -o build/bin/authServer cmd/auth/authServer.go

run-local:
	go run cmd/reverseproxy/rproxy.go
	go run cmd/auth/authServer.go

run-local-docker:
	CGO_ENABLED=0 GOOS=linux go build -o cmd/auth/authServer-scratch cmd/auth/authServer.go
	docker build -t authserver cmd/auth/.
	docker run --rm -p 8080:8080 authserver
	rm cmd/auth/authServer-scratch
