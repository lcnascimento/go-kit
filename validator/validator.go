package validator

import (
	"encoding/json"
	"reflect"
	"time"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/lcnascimento/go-kit/errors"
)

type (
	StructLevel     = validator.StructLevel
	StructLevelFunc = validator.StructLevelFunc
)

type Validator struct {
	*validator.Validate
	trans ut.Translator

	tags Tags
}

func New(opts ...Option) (*Validator, error) {
	v := &Validator{
		Validate: validator.New(),
		tags: map[string]Tag{
			"required":     new(TagRequired),
			"required_if":  new(TagRequired),
			"email":        new(TagEmail),
			"phone":        new(TagPhone),
			"canonical_id": new(TagCanonicalID),
		},
	}

	english := en.New()
	uni := ut.New(english, english)
	v.trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v.Validate, v.trans)

	if err := registerBracketedDefaults(v); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(v)
	}

	for name, tag := range v.tags {
		if err := v.registerTag(name, tag); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (v *Validator) registerTag(name string, tag Tag) error {
	if alias := tag.Alias(); alias != "" {
		v.RegisterAlias(name, alias)
	}

	if validate := tag.Validate(); validate != nil {
		if err := v.RegisterValidation(name, validate); err != nil {
			return err
		}
	}

	translate := tag.Translate()
	if translate == "" {
		return nil
	}

	register := func(ut ut.Translator) error {
		return ut.Add(name, translate, true)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		var strValue string

		if t, ok := fe.Value().(time.Time); ok {
			strValue = t.Format("2006-01-02")
		} else {
			strValue = reflect.ValueOf(fe.Value()).String()
		}

		t, _ := ut.T(name, namespaceOf(fe), strValue, fe.Param())

		return t
	}

	return v.RegisterTranslation(name, v.trans, register, transFn)
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func (v *Validator) Struct(s any) error {
	err := v.Validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return ErrUnexpectedValidationError.WithCause(err)
	}

	output := errors.ErrInvalidInput
	for _, e := range errs {
		output = output.WithCause(errors.New("%s", e.Translate(v.trans)))
	}

	payload, err := json.Marshal(s)
	if err != nil {
		return errors.ErrCastPayload.WithCause(err)
	}

	return output.WithAttribute("payload", string(payload))
}
