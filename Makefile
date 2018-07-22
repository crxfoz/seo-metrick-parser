BINARY_NAME_SERVER=seo_server
BINARY_NAME_WORKER=seo_worker
BINARY_NAME_LINUX_SERVER=seo_server_linux
BINARY_NAME_LINUX_WORKER=seo_worker_linux
DEPLOY_DIR=./deployment

all: test build

run: build-linux
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml -f $(DEPLOY_DIR)/docker-compose.worker.yml stop
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml -f $(DEPLOY_DIR)/docker-compose.worker.yml build
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml -f $(DEPLOY_DIR)/docker-compose.worker.yml up

run-server: build-linux
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml stop
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml build
	docker-compose -f $(DEPLOY_DIR)/docker-compose.yml up

build:
	go build -o $(DEPLOY_DIR)/$(BINARY_NAME_SERVER) ./server/cmd/server
	go build -o $(DEPLOY_DIR)/$(BINARY_NAME_WORKER) ./worker/cmd/worker

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./deployment/$(BINARY_NAME_LINUX_SERVER) ./server/cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./deployment/$(BINARY_NAME_LINUX_WORKER) ./worker/cmd/worker

build-docker: clean
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/crxfoz/seo_metrick_parser golang:latest go build -o $(DEPLOY_DIR)/$(BINARY_NAME_LINUX_SERVER) -v ./server/cmd/server
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/crxfoz/seo_metrick_parser golang:latest go build -o $(DEPLOY_DIR)/$(BINARY_NAME_LINUX_WORKER) -v ./worker/cmd/worker

clean:
	rm -f $(DEPLOY_DIR)/$(BINARY_NAME_SERVER)
	rm -f $(DEPLOY_DIR)/$(BINARY_NAME_WORKER)
	rm -f $(DEPLOY_DIR)/$(BINARY_NAME_LINUX_SERVER)
	rm -f $(DEPLOY_DIR)/$(BINARY_NAME_LINUX_WORKER)

test:
	docker run --rm -p 5432:5432 -d --name seo_parser_postgres_test -e POSTGRES_DB=seo_parser_test postgres
	docker run --rm -p 6379:6379 -d --name seo_parser_redis_test redis
	usql -f ./deployment/scheme.sql "pg://postgres@localhost:5432/seo_parser_test?sslmode=disable"
	GOCACHE=off go test -v ./...
	docker stop seo_parser_postgres_test
	docker stop seo_parser_redis_test
