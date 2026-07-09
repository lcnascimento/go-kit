module github.com/lcnascimento/go-kit/http

go 1.26.4

replace github.com/lcnascimento/go-kit/errors => ../errors

replace github.com/lcnascimento/go-kit/o11y => ../o11y

replace github.com/lcnascimento/go-kit/env => ../env

require (
	github.com/felixge/httpsnoop v1.1.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/lcnascimento/go-kit/env v0.0.0-00010101000000-000000000000
	github.com/lcnascimento/go-kit/errors v0.0.0-00010101000000-000000000000
	github.com/lcnascimento/go-kit/o11y v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.44.1-0.20260626205805-41ff5ed18bec
	go.opentelemetry.io/otel/metric v1.44.1-0.20260625150014-c84013202f01
	go.opentelemetry.io/otel/trace v1.44.1-0.20260625150014-c84013202f01
)

require (
	github.com/caarlos0/env/v10 v10.0.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/bridges/otelslog v0.19.0 // indirect
	go.opentelemetry.io/contrib/processors/minsev v0.16.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.20.0 // indirect
	go.opentelemetry.io/otel/log v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.1-0.20260625150014-c84013202f01 // indirect
	go.opentelemetry.io/otel/sdk/log v0.20.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.39.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/grpc v1.82.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
