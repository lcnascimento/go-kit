package validator

import (
	"reflect"
	"strings"
)

type Option func(*Validator)

func WithTags(tags Tags) Option {
	return func(v *Validator) {
		for name, tag := range tags {
			v.tags[name] = tag
		}
	}
}

func WithJSONFieldNames() Option {
	return func(v *Validator) {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
			if name == "-" {
				return ""
			}

			return name
		})
	}
}
