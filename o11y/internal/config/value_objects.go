package config

import (
	"log/slog"
	"strings"

	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel/attribute"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/lcnascimento/go-kit/env"
)

// TracesExporter is the type for the traces exporter.
type TracesExporter string

const (
	TracesExporterNone   TracesExporter = "none"
	TracesExporterOTLP   TracesExporter = "otlp"
	TracesExporterStdout TracesExporter = "stdout"
)

var availableTracesExporters = []string{
	string(TracesExporterNone),
	string(TracesExporterOTLP),
	string(TracesExporterStdout),
}

// MetricsExporter is the type for the metrics exporter.
type MetricsExporter string

const (
	MetricsExporterNone   MetricsExporter = "none"
	MetricsExporterOTLP   MetricsExporter = "otlp"
	MetricsExporterStdout MetricsExporter = "stdout"
)

var availableMetricsExporters = []string{
	string(MetricsExporterNone),
	string(MetricsExporterOTLP),
	string(MetricsExporterStdout),
}

// MetricsExemplarFilter is a filter for which measurements can become Exemplars.
type MetricsExemplarFilter string

const (
	MetricsExemplarFilterAlwaysOn   MetricsExemplarFilter = "always_on"
	MetricsExemplarFilterAlwaysOff  MetricsExemplarFilter = "always_off"
	MetricsExemplarFilterTraceBased MetricsExemplarFilter = "trace_based"
)

var availableMetricsExemplarFilters = []string{
	string(MetricsExemplarFilterAlwaysOn),
	string(MetricsExemplarFilterAlwaysOff),
	string(MetricsExemplarFilterTraceBased),
}

// LogsExporter is the type for the logs exporter.
type LogsExporter string

const (
	LogsExporterNone   LogsExporter = "none"
	LogsExporterOTLP   LogsExporter = "otlp"
	LogsExporterStdout LogsExporter = "stdout"
)

var availableLogsExporters = []string{
	string(LogsExporterNone),
	string(LogsExporterOTLP),
	string(LogsExporterStdout),
}

// ProfilesExporter is the type for the profiles exporter.
type ProfilesExporter string

const (
	ProfilesExporterNone      ProfilesExporter = "none"
	ProfilesExporterPyroscope ProfilesExporter = "pyroscope"
)

var availableProfilesExporters = []string{
	string(ProfilesExporterNone),
	string(ProfilesExporterPyroscope),
}

// ResourceAttributes is a list of resource attributes.
// More details on: https://opentelemetry.io/docs/specs/semconv/resource/#semantic-attributes-with-dedicated-environment-variable
type ResourceAttributes map[attribute.Key]string

// Load loads the ResourceAttributes from the environment variables.
func (t *ResourceAttributes) Load() ResourceAttributes {
	output := make(ResourceAttributes)

	values := env.GetList[string]("OTEL_RESOURCE_ATTRIBUTES")
	if len(values) == 0 {
		return output
	}

	for _, value := range values {
		parts := strings.Split(value, "=")
		if len(parts) != 2 { //nolint:mnd // OK
			continue
		}

		key := attribute.Key(parts[0])
		value := parts[1]

		if _, ok := availableResourceAttributes[key]; !ok {
			continue
		}

		output[key] = value
	}

	*t = output
	return output
}

var availableResourceAttributes = map[attribute.Key]func(string) attribute.KeyValue{
	semconv.ServiceNameKey:           semconv.ServiceName,
	semconv.ServiceNamespaceKey:      semconv.ServiceNamespace,
	semconv.ServiceVersionKey:        semconv.ServiceVersion,
	semconv.ServiceInstanceIDKey:     semconv.ServiceInstanceID,
	semconv.DeploymentEnvironmentKey: semconv.DeploymentEnvironment,
}

// ToList converts the ResourceAttributes to a list of attribute.KeyValue.
func (t *ResourceAttributes) ToList() []attribute.KeyValue {
	output := make([]attribute.KeyValue, 0)

	for key, value := range *t {
		output = append(output, attribute.String(string(key), value))
	}

	for _, attr := range output {
		if attr.Key == semconv.ServiceNameKey && attr.Value.AsString() != "" {
			return output
		}
	}

	serviceName := env.Get[string]("OTEL_SERVICE_NAME")
	if serviceName != "" {
		output = append(output, attribute.String(string(semconv.ServiceNameKey), serviceName))

		return output
	}

	serviceName = env.Get[string]("OTEL_SERVICE_NAME")
	if serviceName != "" {
		output = append(output, attribute.String(string(semconv.ServiceNameKey), serviceName))
	}

	return output
}

// LogHandler is the type for the log handler.
type LogHandler string

const (
	LogHandlerZap      LogHandler = "zap"
	LogHandlerOtelSlog LogHandler = "otelslog"
)

var availableLogHandlers = []string{
	string(LogHandlerZap),
	string(LogHandlerOtelSlog),
}

// LogLevel is a value object that represents the log level for the Otel SDK.
type LogLevel string

// available log levels.
const (
	LogLevelDebug    LogLevel = "DEBUG"
	LogLevelInfo     LogLevel = "INFO"
	LogLevelWarn     LogLevel = "WARN"
	LogLevelError    LogLevel = "ERROR"
	LogLevelCritical LogLevel = "CRITICAL"
	LogLevelFatal    LogLevel = "FATAL"
)

var availableLogLevels = []string{
	string(LogLevelDebug),
	string(LogLevelInfo),
	string(LogLevelWarn),
	string(LogLevelError),
	string(LogLevelCritical),
	string(LogLevelFatal),
}

// ToSlogLevel converts the LogLevel to a slog.Level.
func (l *LogLevel) ToSlogLevel() slog.Level {
	switch *l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	case LogLevelCritical:
		return slog.Level(12) //nolint:mnd // OK
	case LogLevelFatal:
		return slog.Level(14) //nolint:mnd // OK
	default:
		return slog.LevelInfo
	}
}

// ToSeverity converts the LogLevel to a severity.
func (l *LogLevel) ToSeverity() minsev.Severity {
	switch *l {
	case LogLevelDebug:
		return minsev.SeverityDebug
	case LogLevelInfo:
		return minsev.SeverityInfo
	case LogLevelWarn:
		return minsev.SeverityWarn
	case LogLevelError:
		return minsev.SeverityError
	case LogLevelCritical:
		return minsev.SeverityFatal1
	case LogLevelFatal:
		return minsev.SeverityFatal3
	default:
		return minsev.SeverityInfo
	}
}

// OtlpExporterProtocol is the protocol to be used for the OTLP exporter.
type OtlpExporterProtocol string

const (
	OtlpExporterProtocolGRPC OtlpExporterProtocol = "grpc"
)

var availableOtlpExporterProtocols = []string{
	string(OtlpExporterProtocolGRPC),
}

// Propagator defines how data should be propagated through the system.
type Propagator string

// available propagators.
const (
	PropagatorTraceContext Propagator = "tracecontext"
	PropagatorBaggage      Propagator = "baggage"
)

var availablePropagators = []Propagator{
	PropagatorTraceContext,
	PropagatorBaggage,
}

// Propagators is a list of propagators.
type Propagators map[Propagator]bool

// Load loads the Propagators from the environment variables.
func (p *Propagators) Load() Propagators {
	output := make(Propagators)

	values := env.GetList("OTEL_PROPAGATORS",
		env.WithEnum(availablePropagators),
		env.WithDefaultListValue([]Propagator{PropagatorTraceContext, PropagatorBaggage}),
	)
	if len(values) == 0 {
		return output
	}

	for _, value := range values {
		output[value] = true
	}

	*p = output
	return output
}

// TraceSampler is a value object that represents the trace sampler for the Otel SDK.
type TraceSampler string

// available trace samplers.
const (
	TraceSamplerAlwaysOn                TraceSampler = "always_on"
	TraceSamplerAlwaysOff               TraceSampler = "always_off"
	TraceSamplerParentBasedAlwaysOn     TraceSampler = "parentbased_always_on"
	TraceSamplerParentBasedAlwaysOff    TraceSampler = "parentbased_always_off"
	TraceSamplerParentBasedTraceIDRatio TraceSampler = "parentbased_traceidratio"
)

var availableTraceSamplers = []string{
	string(TraceSamplerAlwaysOn),
	string(TraceSamplerAlwaysOff),
	string(TraceSamplerParentBasedAlwaysOn),
	string(TraceSamplerParentBasedAlwaysOff),
	string(TraceSamplerParentBasedTraceIDRatio),
}

// Profile is the type for the profile.
type Profile string

const (
	ProfileCPU       Profile = "cpu"
	ProfileHeap      Profile = "heap"
	ProfileGoroutine Profile = "goroutine"
	ProfileMutex     Profile = "mutex"
)

var availableProfiles = []Profile{
	ProfileCPU,
	ProfileHeap,
	ProfileGoroutine,
	ProfileMutex,
}

// Profiles is a map of profiles.
type Profiles map[Profile]bool

// Load loads the value of the Profiles.
func (p *Profiles) Load() Profiles {
	output := make(Profiles, 0)

	values := env.GetList("OTEL_PROFILES",
		env.WithEnum(availableProfiles),
		env.WithDefaultListValue([]Profile{ProfileCPU, ProfileHeap}),
	)
	if len(values) == 0 {
		return output
	}

	for _, prof := range values {
		output[prof] = true
	}

	*p = output
	return output
}
