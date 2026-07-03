//nolint:revive // OK
package log

import "log/slog"

// wrappers around slog functions to avoid importing slog and common's log package when logging data.
var (
	String     = slog.String
	Int        = slog.Int
	Float      = slog.Float64
	Bool       = slog.Bool
	Any        = slog.Any
	Group      = slog.Group
	GroupValue = slog.GroupValue

	// Object wraps slog any, to provide an zap like interface.
	Object = slog.Any
)

// wrappers around slog type to avoid importing slog and common's log package when logging data.
type (
	LogValuer = slog.LogValuer //nolint:revive // OK
	Value     = slog.Value
	Attr      = slog.Attr
)

// Custom keys to be used with slog.
const (
	ErrorKey  = "error"
	LoggerKey = "logger"
)

// Level is a wrapper around slog.Leveler.
// It adds some extra Levels, like Critical.
const (
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelWarn     = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelCritical = slog.Level(12)
	LevelFatal    = slog.Level(14)
)
