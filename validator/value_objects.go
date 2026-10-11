package validator

import (
	"fmt"
	"reflect"

	validator "github.com/go-playground/validator/v10"

	"github.com/lcnascimento/go-kit/util"
)

type Tags map[string]Tag

type Tag interface {
	Alias() string
	Validate() func(fl validator.FieldLevel) bool
	Translate() string
}

type UnimplementedTag struct{}

func (t UnimplementedTag) Alias() string {
	return ""
}

func (t UnimplementedTag) Validate() func(fl validator.FieldLevel) bool {
	return nil
}

func (t UnimplementedTag) Translate() string {
	return ""
}

type TagRequired struct {
	UnimplementedTag
}

func (t TagRequired) Translate() string {
	return "[{0}] is a required field"
}

type TagEmail struct {
	UnimplementedTag
}

func (t TagEmail) Translate() string {
	return "[{0}] must be a valid email address"
}

type TagPhone struct {
	UnimplementedTag
}

func (t TagPhone) Translate() string {
	return "[{0}] must be a valid E.164 formatted phone number"
}

func (t TagPhone) Alias() string {
	return "e164"
}

// TagCanonicalID validates that a field holds a util.CanonicalID with the prefix
// given as the tag param, e.g. `validate:"canonical_id=note"`. The field may be a
// string or any type whose String method renders the id, such as util.CanonicalID
// itself or a struct embedding it. Its message is registered in
// paramBracketedTemplates, because the tag param must be the {1} placeholder.
type TagCanonicalID struct {
	UnimplementedTag
}

func (t TagCanonicalID) Validate() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		var raw string

		field := fl.Field()
		switch {
		case field.Kind() == reflect.String:
			raw = field.String()
		case field.CanInterface():
			s, ok := field.Interface().(fmt.Stringer)
			if !ok {
				return false
			}

			raw = s.String()
		default:
			return false
		}

		id, err := util.ParseCanonicalID(fl.Param(), raw)

		return err == nil && !id.IsZero()
	}
}
