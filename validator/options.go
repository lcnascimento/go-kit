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

			switch {
			case name == "-":
				return ""
			case name == "" && fld.Anonymous && isStruct(fld.Type):
				return embeddedSegment
			}

			return name
		})
	}
}

// embeddedSegment names an embedded struct without a json name in the
// validation namespace, so namespaceOf can drop it the way encoding/json
// flattens the struct's fields into its parent.
const embeddedSegment = "\x00embedded"

func isStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}
