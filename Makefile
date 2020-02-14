export BIN_DIR=bin
export PORT=8000
export MACHINE_ID=0
export USE_CACHE=1
export CACHE_EXPIRE=10
export TARGETS=http://localhost:8001;http://localhost:8002
export STALE_TIMEOUT=1

export IMAGE_NAME=freundallein/schooner:latest

init:
	git config core.hooksPath .githooks
run:
	go run schooner.go
test:
	go test -cover ./...
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $$BIN_DIR/schooner
build-healthchecker:
	cd healthchecker && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o ../$$BIN_DIR/healthchecker
dockerbuild:
	make test
	docker build -t $$IMAGE_NAME -f Dockerfile .
distribute:
	make test
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker build -t $$IMAGE_NAME .
	docker push $$IMAGE_NAME