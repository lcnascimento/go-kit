# Standardized tool versions
include .env
export

GOPATH=$(shell go env GOPATH)

.PHONY: producer
producer:
	@OTEL_SERVICE_NAME=kafka-example-producer go run example/producer/main.go

.PHONY: consumer
consumer:
	@OTEL_SERVICE_NAME=kafka-example-consumer go run example/consumer/main.go
