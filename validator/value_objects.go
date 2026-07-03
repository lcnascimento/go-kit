package validator

import validator "github.com/go-playground/validator/v10"

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

func (p TagPhone) Alias() string {
	return "e164"
}
