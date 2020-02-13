export BIN_DIR=bin
export PORT=8000

export IMAGE_NAME=freundallein/schooner:latest

init:
	git config core.hooksPath .githooks
run:
	go run main.go
test:
	go test -cover ./...
build:
	make test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $$BIN_DIR/schooner
dockerbuild:
	make test
	docker build -t $$IMAGE_NAME -f Dockerfile .
distribute:
	make test
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker build -t $$IMAGE_NAME .
	docker push $$IMAGE_NAME