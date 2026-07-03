package interceptor

import "github.com/lcnascimento/go-kit/errors"

var ErrPanic = errors.New("panic").WithCode("ERR_PANIC").WithKind(errors.KindCritical)
