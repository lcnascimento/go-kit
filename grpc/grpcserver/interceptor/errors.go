package interceptor

import "github.com/lcnascimento/go-kit/errors"

var ErrPanic = errors.New("panic").WithCode("PANIC").WithKind(errors.KindCritical)
