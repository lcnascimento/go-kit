package env

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// builtInParsers is a map of built-in parsers for the basic types.
var builtInParsers = map[reflect.Kind]func(string) (any, error){
	reflect.Bool: func(v string) (any, error) {
		return strconv.ParseBool(v)
	},
	reflect.String: func(v string) (any, error) {
		return v, nil
	},
	reflect.Int: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	},
	reflect.Int16: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	},
	reflect.Int32: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	},
	reflect.Int64: func(v string) (any, error) {
		if duration, err := time.ParseDuration(v); err == nil {
			return duration, nil
		}

		return strconv.ParseInt(v, 10, 64)
	},
	reflect.Int8: func(v string) (any, error) {
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	},
	reflect.Uint: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	},
	reflect.Uint16: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	},
	reflect.Uint32: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	},
	reflect.Uint64: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	},
	reflect.Uint8: func(v string) (any, error) {
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	},
	reflect.Float64: func(v string) (any, error) {
		return strconv.ParseFloat(v, 64)
	},
	reflect.Float32: func(v string) (any, error) {
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	},
}

// getterConfig is the configuration for the Get function.
// defaultValue is the default value to return if the environment variable is not defined. It is optional.
// customParser is a custom parser to convert the environment variable to the type T. It is optional.
// genericType is the generic type of the type T.
type getterConfig[T any] struct {
	defaultValue     *T
	defaultListValue []T
	listSeparator    string
	customParser     func(string) (T, error)
	genericType      reflect.Type
	enumValues       []T
}

// getParser returns the parser function for the type T. It returns nil if the type T is a pointer or if there is no
// parser for the type T.
func (c *getterConfig[T]) getParser() func(string) (T, error) {
	if c.customParser != nil {
		return c.customParser
	}

	kind := c.genericType.Kind()
	if kind == reflect.Ptr {
		return nil
	}

	parser, ok := builtInParsers[kind]
	if !ok {
		return nil
	}

	// adapt parser func
	return func(v string) (T, error) {
		parsedValue, err := parser(v)

		return parsedValue.(T), err
	}
}

// getDefaultOrZero returns the default value or zero value of the type T.
func (c *getterConfig[T]) getDefaultOrZero() T {
	if c.defaultValue != nil {
		return *c.defaultValue
	}

	var zero T
	return zero
}

func (c *getterConfig[T]) getDefaultListValue() []T {
	if len(c.defaultListValue) > 0 {
		return c.defaultListValue
	}

	return []T{}
}

func (c *getterConfig[T]) getListSeparator() string {
	if c.listSeparator != "" {
		return c.listSeparator
	}

	return ","
}

// validateEnum validates if the value is in the enum values.
func (c *getterConfig[T]) validateEnum(value T) bool {
	if len(c.enumValues) == 0 {
		return true
	}

	for _, v := range c.enumValues {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}

	return false
}

type GetterOption[T any] func(cfg *getterConfig[T])

// WithDefaultValue sets the default value to return if the environment variable is not defined or if its value is zero.
func WithDefaultValue[T any](value T) GetterOption[T] {
	return func(cfg *getterConfig[T]) {
		cfg.defaultValue = &value
	}
}

// WithDefaultListValue sets the default value to return if the environment variable is not defined or if its value is an empty list.
func WithDefaultListValue[T any](values []T) GetterOption[T] {
	return func(cfg *getterConfig[T]) {
		cfg.defaultListValue = values
	}
}

// WithListSeparator sets the separator for the list of values.
func WithListSeparator[T any](separator string) GetterOption[T] {
	return func(cfg *getterConfig[T]) {
		cfg.listSeparator = separator
	}
}

// WithEnum sets the enum values to validate the environment variable against.
func WithEnum[T any](values []T) GetterOption[T] {
	return func(cfg *getterConfig[T]) {
		cfg.enumValues = values
	}
}

// WithCustomParser sets a custom parser to convert the environment variable to the type T.
func WithCustomParser[T any](parser func(string) (T, error)) GetterOption[T] {
	return func(cfg *getterConfig[T]) {
		cfg.customParser = parser
	}
}

// buildGetterConfig builds the getterConfig from the options.
func buildGetterConfig[T any](opts []GetterOption[T]) getterConfig[T] {
	var zero [0]T

	cfg := getterConfig[T]{
		genericType: reflect.TypeOf(zero).Elem(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// Get retrieves the value of the environment variable with the given name. It converts the value to the type T. If the
// environment variable is not defined, it returns the default value or zero. If the environment variable is defined but
// cannot be converted to the type T, it returns the default value or zero.
//
// POINTERS ARE NOT SUPPORTED. If T is a pointer, it ALWAYS returns nil.
func Get[T any](envName string, opts ...GetterOption[T]) T {
	cfg := buildGetterConfig(opts)

	parser := cfg.getParser()
	if parser == nil {
		return cfg.getDefaultOrZero()
	}

	envValueStr, ok := os.LookupEnv(envName)
	if !ok {
		return cfg.getDefaultOrZero()
	}

	parsedValue, err := parser(envValueStr)
	if err != nil {
		return cfg.getDefaultOrZero()
	}

	if !cfg.validateEnum(parsedValue) {
		return cfg.getDefaultOrZero()
	}

	return parsedValue
}

// GetGetList retrieves a list of values from the environment variable with the given name. The values are expected to be
// separated by commas. It converts each value to the type T. If the environment variable is not defined, it returns an empty list.
// If the environment variable is defined but cannot be converted to the type T, it returns an empty list.
//
// POINTERS ARE NOT SUPPORTED. If T is a pointer, it ALWAYS returns an empty list.
func GetList[T any](envName string, opts ...GetterOption[T]) []T {
	cfg := buildGetterConfig(opts)

	parser := cfg.getParser()
	if parser == nil {
		return cfg.getDefaultListValue()
	}

	envValueStr, ok := os.LookupEnv(envName)
	if !ok {
		return cfg.getDefaultListValue()
	}

	output := make([]T, 0)

	values := strings.Split(envValueStr, cfg.getListSeparator())
	for _, value := range values {
		parsed, err := parser(value)
		if err != nil {
			continue
		}

		if !cfg.validateEnum(parsed) {
			continue
		}

		output = append(output, parsed)
	}

	return output
}
