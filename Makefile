tests:
	go test ./...

builds:
	go build -o build/bin/rproxy cmd/reverseproxy/rproxy.go
	go build -o build/bin/authServer cmd/jServer/*.go
	go build -o build/bin/loggerd cmd/loggerd/loggerd.go

run-local:
	go run cmd/reverseproxy/rproxy.go &
	go run cmd/jServer/jServer.go

docker-build:
	docker build -t jserver cmd/jServer/.
	docker build -t rproxy cmd/reverseproxy/.
	docker build -t loggerd cmd/loggerd/.

docker-build-no-cache:
	docker build --no-cache -t jserver cmd/jServer/.
	docker build --no-cache -t rproxy cmd/reverseproxy/.
	docker build --no-cache -t loggerd cmd/loggerd/.

docker-run:
	docker run --rm --env-file .env --net dev-bridge --name logger -p 6666:9090 loggerd &
	docker run --rm --env-file .env --net dev-bridge --name jserver -p 4444:8080 jserver &
	docker run --rm --env-file .env --net dev-bridge --name rproxy -p 5555:8080 rproxy

docker-bridge:
	docker network create dev-bridge
	
	
