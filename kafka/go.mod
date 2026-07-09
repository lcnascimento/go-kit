module github.com/lcnascimento/go-kit/kafka

go 1.26.4

replace github.com/lcnascimento/go-kit/env => ../env

require (
	github.com/google/uuid v1.6.0
	github.com/lcnascimento/go-kit/env v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.51
	go.opentelemetry.io/otel v1.44.1-0.20260626205805-41ff5ed18bec
	go.opentelemetry.io/otel/trace v1.44.1-0.20260625150014-c84013202f01
)

require (
	github.com/caarlos0/env/v10 v10.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.19.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.44.1-0.20260625150014-c84013202f01 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/text v0.39.0 // indirect
)
