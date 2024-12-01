package validator

// Option is a function that configures a Validator.
type Option func(*val)

// WithCustomValidator adds custom validators to the validator.
func WithCustomValidator(cv CustomValidator) Option {
	return func(v *val) {
		_ = v.validator.RegisterValidation(cv.Tag(), cv.Func())
	}
}
