package global

import (
	"github.com/lcnascimento/go-kit/o11y/internal/config"
)

var cfg *config.Configuration

// Config returns the global configuration.
func Config() *config.Configuration {
	if cfg != nil {
		return cfg
	}

	cfg = &config.Configuration{}
	cfg.Load()

	return cfg
}
