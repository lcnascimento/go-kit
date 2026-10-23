module github.com/lcnascimento/go-kit/grpc

go 1.27.3

replace github.com/lcnascimento/go-kit/errors => ../errors

replace github.com/lcnascimento/go-kit/o11y => ../o11y

replace github.com/lcnascimento/go-kit/env => ../env

require (
	github.com/lcnascimento/go-kit/errors v0.0.0-00010101000000-000000000000
	github.com/lcnascimento/go-kit/o11y v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.57.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d
	google.golang.org/grpc v1.84.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/caarlos0/env/v10 v10.0.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.19.0 // indirect
	github.com/lcnascimento/go-kit/env v0.0.0-00010101000000-000000000000 // indirect
	github.com/posthog/posthog-go v1.36.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.20.0 // indirect
	go.opentelemetry.io/contrib/processors/minsev v0.16.2 // indirect
	go.opentelemetry.io/otel v1.47.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.21.0 // indirect
	go.opentelemetry.io/otel/log v1.47.0 // indirect
	go.opentelemetry.io/otel/metric v1.47.0 // indirect
	go.opentelemetry.io/otel/sdk v1.45.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.47.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
