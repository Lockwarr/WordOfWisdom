install:
	go mod download

test-with-component:
	go test --tags=component ./...

test:
	go test ./...

start-server:
	go run server/cmd/main.go

start-client:
	go run client/cmd/main.go

build-docker:
	docker-compose build

run-on-docker:
	docker-compose up

