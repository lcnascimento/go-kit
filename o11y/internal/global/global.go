package global

import (
	"sync/atomic"

	"github.com/lcnascimento/go-kit/o11y/internal/config"
)

var cfg atomic.Pointer[config.Configuration]

func Config() *config.Configuration {
	if c := cfg.Load(); c != nil {
		return c
	}

	return Load()
}

func Load() *config.Configuration {
	c := &config.Configuration{}
	c.Load()
	cfg.Store(c)

	return c
}
