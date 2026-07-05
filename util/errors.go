package util

import "github.com/lcnascimento/go-kit/errors"

var (
	ErrParseID         = errors.New("failed to parse id").WithCode("ERR_PARSE_ID")
	ErrInvalidIDPrefix = func(prefix string) error {
		return errors.New("id must have prefix '%s_'", prefix).WithCode("ERR_INVALID_ID_PREFIX").WithKind(errors.KindInvalidInput)
	}
)
