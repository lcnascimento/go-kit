package handlers

import (
	"log/slog"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"

	"github.com/lcnascimento/go-kit/o11y/internal/global"
)

// Zap returns a slog.Handler that uses zap as the core logger.
func Zap(name string) (slog.Handler, error) {
	cfg := global.Config()

	var level zapcore.Level
	switch cfg.LogLevel.ToSlogLevel() {
	case slog.LevelDebug:
		level = zap.DebugLevel
	case slog.LevelWarn:
		level = zap.WarnLevel
	case slog.LevelError:
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	if cfg.PrettyPrint {
		zapConfig.EncoderConfig = localEncoderConfig
		zapConfig.Encoding = "console"
	} else {
		zapConfig.EncoderConfig = cloudEncoderConfig
	}

	z, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(z)
	return zapslog.NewHandler(z.Core(), zapslog.WithName(name)), nil
}

var (
	localEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "severity",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     prettyTimestampEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// https://cloud.google.com/logging/docs/structured-logging#structured_logging_special_fields
	cloudEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "severity",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "logging.googleapis.com/sourceLocation",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    severityEncoder,
		EncodeTime:     timestampEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   sourceLocationEncoder,
	}
)

type timestamp struct {
	Seconds int64
	Nanos   int
}

func (t timestamp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("seconds", t.Seconds)
	enc.AddInt("nanos", t.Nanos)
	return nil
}

// timestampEncoder is a encoder for timestamp following Google Cloud patterns.
//
// https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
func timestampEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	if arrayEnc, ok := enc.(zapcore.ArrayEncoder); ok {
		_ = arrayEnc.AppendObject(timestamp{
			Seconds: t.Unix(),
			Nanos:   t.Nanosecond(),
		})
	} else {
		zapcore.RFC3339NanoTimeEncoder(t, enc)
	}
}

// prettyTimestampEncoder is a encoder for output pretty timestamp, normally used locally.
func prettyTimestampEncoder(t time.Time, e zapcore.PrimitiveArrayEncoder) {
	e.AppendString(t.Format("15:04:05.000"))
}

type sourceLocation struct {
	File     string
	Line     int
	Function string
}

func (l sourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.File)
	enc.AddString("line", strconv.Itoa(l.Line))
	enc.AddString("function", l.Function)
	return nil
}

// sourceLocationEncoder is a encoder for SourceLocation following Google Cloud patterns.
//
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logentrysourcelocation
func sourceLocationEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if arrayEnc, ok := enc.(zapcore.ArrayEncoder); ok {
		_ = arrayEnc.AppendObject(sourceLocation{
			File:     caller.File,
			Line:     caller.Line,
			Function: caller.Function,
		})
	} else {
		enc.AppendString(caller.TrimmedPath())
	}
}

// severityEncoder is an encoder for severity following Google Cloud patterns.
//
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
func severityEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[l])
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}
