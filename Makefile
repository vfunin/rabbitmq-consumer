FLAGS:= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
CMD:= ./cmd/rabbitmq-consumer

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all
all: up build create-queues

.PHONY: build
build:
	$(info Start build...)
	go build -a -tags netgo -ldflags="-w -extldflags '-static'" -o ./rabbitmq-consumer $(CMD)

.PHONY: build-flags
build-flags:
	$(info Start build...)
	$(FLAGS) go build -a -tags netgo -ldflags="-w -extldflags '-static'" -o ./rabbitmq-consumer $(CMD)

.PHONY: run
run:
	$(info Run...)
	go run cmd/rabbitmq-consumer/main.go -config config.yaml

.PHONY: up
up:
	$(info Start containers...)
	docker compose up -d

.PHONY: down
down:
	$(info Stop containers)
	docker compose down

.PHONY: create-queues
create-queues:
	$(info Create rabbitmq queues...)
	docker exec rabbitmq-consumer-rabbitmq /usr/local/bin/rabbitmqadmin declare queue --vhost=$(RABBITMQ_VHOST) name=$(RABBITMQ_QUEUE_NAME) durable=true -u $(RABBITMQ_USER) -p $(RABBITMQ_PASSWORD)
