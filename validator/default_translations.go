package validator

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

// registerBracketedDefaults overrides go-playground's default English
// translations so every message is prefixed with the full field namespace
// wrapped in brackets, e.g. "[Personal.Birthdate] is a required field".
//
// Whenever a new default tag (https://pkg.go.dev/github.com/go-playground/validator/v10)
// starts being used in any service, add an entry here so its message follows
// the same convention.
func registerBracketedDefaults(v *Validator) error {
	if err := registerSimpleBracketed(v); err != nil {
		return err
	}

	return registerKindBracketed(v)
}

func namespaceOf(fe validator.FieldError) string {
	parts := strings.Split(fe.Namespace(), ".")

	namespace := strings.Join(parts[1:], ".")
	if namespace == "" {
		namespace = fe.Field()
	}

	return namespace
}

var simpleBracketedTemplates = map[string]string{
	"required_with":        "[{0}] is a required field",
	"required_with_all":    "[{0}] is a required field",
	"required_without":     "[{0}] is a required field",
	"required_without_all": "[{0}] is a required field",
	"required_unless":      "[{0}] is a required field",
	"number":               "[{0}] must be a valid number",
	"numeric":              "[{0}] must be a valid numeric value",
}

func registerSimpleBracketed(v *Validator) error {
	for tag, template := range simpleBracketedTemplates {
		name, tmpl := tag, template

		register := func(ut ut.Translator) error {
			return ut.Add(name, tmpl, true)
		}

		transFn := func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(name, namespaceOf(fe))

			return t
		}

		if err := v.RegisterTranslation(name, v.trans, register, transFn); err != nil {
			return err
		}
	}

	return nil
}

type kindBracketedTag struct {
	name    string
	stringT string
	numberT string
	itemsT  string
}

var kindBracketedTags = []kindBracketedTag{
	{
		name:    "len",
		stringT: "[{0}] must be {1} in length",
		numberT: "[{0}] must be equal to {1}",
		itemsT:  "[{0}] must contain {1}",
	},
	{
		name:    "min",
		stringT: "[{0}] must be at least {1} in length",
		numberT: "[{0}] must be {1} or greater",
		itemsT:  "[{0}] must contain at least {1}",
	},
	{
		name:    "max",
		stringT: "[{0}] must be a maximum of {1} in length",
		numberT: "[{0}] must be {1} or less",
		itemsT:  "[{0}] must contain at maximum {1}",
	},
}

func registerKindBracketed(v *Validator) error {
	for _, kt := range kindBracketedTags {
		tag := kt

		register := func(ut ut.Translator) error {
			if err := ut.Add(tag.name+"-string", tag.stringT, true); err != nil {
				return err
			}

			if err := ut.Add(tag.name+"-number", tag.numberT, true); err != nil {
				return err
			}

			return ut.Add(tag.name+"-items", tag.itemsT, true)
		}

		transFn := func(ut ut.Translator, fe validator.FieldError) string {
			return translateKindAware(ut, fe, tag)
		}

		if err := v.RegisterTranslation(tag.name, v.trans, register, transFn); err != nil {
			return err
		}
	}

	return nil
}

func translateKindAware(ut ut.Translator, fe validator.FieldError, kt kindBracketedTag) string {
	namespace := namespaceOf(fe)

	kind := fe.Kind()
	if kind == reflect.Ptr {
		kind = fe.Type().Elem().Kind()
	}

	var digits uint64

	if idx := strings.Index(fe.Param(), "."); idx != -1 {
		digits = uint64(len(fe.Param()[idx+1:]))
	}

	switch kind {
	case reflect.String:
		f64, err := strconv.ParseFloat(fe.Param(), 64)
		if err != nil {
			return namespace + " " + kt.name + " " + fe.Param()
		}

		c, _ := ut.C(kt.name+"-string-character", f64, digits, ut.FmtNumber(f64, digits))
		t, _ := ut.T(kt.name+"-string", namespace, c)

		return t

	case reflect.Slice, reflect.Map, reflect.Array:
		f64, err := strconv.ParseFloat(fe.Param(), 64)
		if err != nil {
			return namespace + " " + kt.name + " " + fe.Param()
		}

		c, _ := ut.C(kt.name+"-items-item", f64, digits, ut.FmtNumber(f64, digits))
		t, _ := ut.T(kt.name+"-items", namespace, c)

		return t

	default:
		if fe.Type() == reflect.TypeOf(time.Duration(0)) {
			t, _ := ut.T(kt.name+"-number", namespace, fe.Param())

			return t
		}

		f64, err := strconv.ParseFloat(fe.Param(), 64)
		if err != nil {
			return namespace + " " + kt.name + " " + fe.Param()
		}

		t, _ := ut.T(kt.name+"-number", namespace, ut.FmtNumber(f64, digits))

		return t
	}
}
