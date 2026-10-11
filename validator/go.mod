module github.com/lcnascimento/go-kit/validator

go 1.27.2

replace github.com/lcnascimento/go-kit/errors => ../errors

replace github.com/lcnascimento/go-kit/util => ../util

replace github.com/lcnascimento/go-kit/o11y => ../o11y

replace github.com/lcnascimento/go-kit/env => ../env

require (
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.23.0
	github.com/lcnascimento/go-kit/errors v0.0.0-00010101000000-000000000000
	github.com/lcnascimento/go-kit/util v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lcnascimento/go-kit/o11y v0.0.0-00010101000000-000000000000 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/otel v1.45.0 // indirect
	go.opentelemetry.io/otel/trace v1.45.0 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
