package env

import (
	"net/url"
	"reflect"

	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v10"
)

type (
	// ParseOptions is an option for the Parse.
	ParseOptions func(pc *ParseConfig)

	// ParseConfig is a configuration for the Parser.
	ParseConfig struct {
		customTypes map[reflect.Type]env.ParserFunc
		prefix      string
	}

	// CustomTypeParserFunc is a generic function that can parse a string into a custom type.
	CustomTypeParserFunc[T any] func(string) (T, error)
)

// WithCustomTypeParser adds a custom type parser to the config. Same as ParseConfig.AddCustomTypeParser.
func WithCustomTypeParser[T any](p CustomTypeParserFunc[T]) ParseOptions {
	return func(pc *ParseConfig) {
		pc.AddCustomTypeParser(CustomTypeParser(p))
	}
}

// FromFiles loads the environment variables from the given files.
// Same as LoadFromFile.
func FromFiles(filenames ...string) ParseOptions {
	return func(_ *ParseConfig) {
		_ = LoadFromFile(filenames...)
	}
}

// WithPrefix sets a prefix for the environment variable keys.
func WithPrefix(prefix string) ParseOptions {
	return func(pc *ParseConfig) {
		pc.AddPrefix(prefix)
	}
}

// AddPrefix adds a prefix to the config.
func (c *ParseConfig) AddPrefix(prefix string) *ParseConfig {
	c.prefix = prefix
	return c
}

// AddCustomTypeParser adds a custom type parser to the config.
func (c *ParseConfig) AddCustomTypeParser(r reflect.Type, pf func(v string) (interface{}, error)) *ParseConfig {
	if c.customTypes == nil {
		c.customTypes = map[reflect.Type]env.ParserFunc{}
	}

	c.customTypes[r] = pf
	return c
}

// envOptions returns the env.Options for the config.
func (c *ParseConfig) envOptions() env.Options {
	return env.Options{
		FuncMap: c.customTypes,
		Prefix:  c.prefix,
	}
}

// CustomTypeParser is a helper function to create a custom parser.
func CustomTypeParser[T any](p CustomTypeParserFunc[T]) (reflect.Type, func(v string) (interface{}, error)) {
	var zero [0]T
	parserType := reflect.TypeOf(zero).Elem()
	parserFunc := func(v string) (interface{}, error) {
		return p(v)
	}

	return parserType, parserFunc
}

func configFromOptions(opts []ParseOptions) *ParseConfig {
	pc := &ParseConfig{
		customTypes: defaultTypeParsers(),
	}

	for _, opt := range opts {
		opt(pc)
	}

	return pc
}

// defaultTypeParsers returns the default type parsers.
func defaultTypeParsers() map[reflect.Type]env.ParserFunc {
	return map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(url.URL{}): func(v string) (interface{}, error) {
			u, err := url.Parse(v)
			if err != nil {
				return nil, err
			}
			return *u, nil
		},
	}
}

// Parse parses the environment variables into the given struct.
func Parse(v any, opts ...ParseOptions) error {
	pc := configFromOptions(opts)

	return ParseWithConfig(v, pc)
}

// ParseWithConfig parses the environment variables into the given struct. It uses a config for the parsing.
func ParseWithConfig(v any, pc *ParseConfig) error {
	return env.ParseWithOptions(v, pc.envOptions())
}

// LoadFromFile loads the environment variables from the given files.
// By default, it loads from the .env file.
func LoadFromFile(filenames ...string) error {
	return godotenv.Load(filenames...)
}
