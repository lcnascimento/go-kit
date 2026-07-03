//nolint:mnd // OK
package config

import (
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/lcnascimento/go-kit/env"
)

// Configuration is a general struct that Otel configurations extracted from the environment variables.
// More details on: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/#general-sdk-configuration
type Configuration struct {
	Disabled                           bool
	ServiceName                        string
	ResourceAttributes                 ResourceAttributes
	LogLevel                           LogLevel
	Propagators                        Propagators
	TraceSampler                       TraceSampler
	TraceSamplerArg                    float64
	BspScheduleDelay                   time.Duration
	BspExportTimeout                   time.Duration
	BspMaxQueueSize                    uint
	BspMaxExportBatchSize              uint
	BlrpScheduleDelay                  time.Duration
	BlrpExportTimeout                  time.Duration
	BlrpMaxQueueSize                   uint
	BlrpMaxExportBatchSize             uint
	AttributeValueLengthLimit          uint
	AttributeCountLimit                uint
	SpanAttributeValueLengthLimit      uint
	SpanAttributeCountLimit            uint
	SpanEventCountLimit                uint
	SpanEventAttributeCountLimit       uint
	SpanLinkCountLimit                 uint
	SpanLinkAttributeCountLimit        uint
	LogRecordAttributeValueLengthLimit uint
	LogRecordAttributeCountLimit       uint
	TracesExporter                     TracesExporter
	MetricsExporter                    MetricsExporter
	LogsExporter                       LogsExporter
	MetricsExemplarFilter              MetricsExemplarFilter
	MetricExportInterval               time.Duration
	MetricExportTimeout                time.Duration

	// OTLP
	OtlpExporterProtocol        OtlpExporterProtocol
	OtlpExporterEndpoint        string
	OtlpTracesExporterEndpoint  string
	OtlpMetricsExporterEndpoint string
	OtlpLogsExporterEndpoint    string
	// implement security later!
	OtlpInsecure        bool
	OtlpTracesInsecure  bool
	OtlpMetricsInsecure bool
	OtlpLogsInsecure    bool

	// Non-Official
	PrettyPrint           bool
	LogHandler            LogHandler
	ProfilesExporter      ProfilesExporter
	Profiles              Profiles
	ProfileExportInterval time.Duration
}

// Load loads the configuration from the environment variables.
//
//nolint:funlen // OK
func (c *Configuration) Load() {
	c.Disabled = env.Get("OTEL_SDK_DISABLED", env.WithDefaultValue(false))
	c.ServiceName = env.Get[string]("OTEL_SERVICE_NAME")
	if c.ServiceName == "" {
		c.ServiceName = "unknown_service_name"
	}

	c.ResourceAttributes = new(ResourceAttributes).Load()
	c.Propagators = new(Propagators).Load()
	c.Profiles = new(Profiles).Load()

	if version := env.Get[string]("OTEL_SERVICE_VERSION"); version != "" {
		c.ResourceAttributes[semconv.ServiceVersionKey] = version
	}

	c.LogLevel = LogLevel(env.Get("OTEL_LOG_LEVEL",
		env.WithDefaultValue(string(LogLevelInfo)),
		env.WithEnum(availableLogLevels),
	))
	c.TraceSampler = TraceSampler(env.Get("OTEL_TRACE_SAMPLER",
		env.WithDefaultValue(string(TraceSamplerParentBasedAlwaysOn)),
		env.WithEnum(availableTraceSamplers),
	))
	c.TraceSamplerArg = env.Get("OTEL_TRACE_SAMPLER_ARG",
		env.WithDefaultValue(1.0),
	)
	c.BspScheduleDelay = env.Get("OTEL_BSP_SCHEDULE_DELAY",
		env.WithDefaultValue(5*time.Second),
	)
	c.BspExportTimeout = env.Get("OTEL_BSP_EXPORT_TIMEOUT",
		env.WithDefaultValue(30*time.Second),
	)
	c.BspMaxQueueSize = env.Get("OTEL_BSP_MAX_QUEUE_SIZE",
		env.WithDefaultValue(uint(5120)),
	)
	c.BspMaxExportBatchSize = env.Get("OTEL_BSP_MAX_EXPORT_BATCH_SIZE",
		env.WithDefaultValue(uint(5120)),
	)
	c.BlrpScheduleDelay = env.Get("OTEL_BLRP_SCHEDULE_DELAY",
		env.WithDefaultValue(1*time.Second),
	)
	c.BlrpExportTimeout = env.Get("OTEL_BLRP_EXPORT_TIMEOUT",
		env.WithDefaultValue(30*time.Second),
	)
	c.BlrpMaxQueueSize = env.Get("OTEL_BLRP_MAX_QUEUE_SIZE",
		env.WithDefaultValue(uint(2048)),
	)
	c.BlrpMaxExportBatchSize = env.Get("OTEL_BLRP_MAX_EXPORT_BATCH_SIZE",
		env.WithDefaultValue(uint(2048)),
	)
	c.AttributeValueLengthLimit = env.Get("OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT",
		env.WithDefaultValue(uint(1024)),
	)
	c.AttributeCountLimit = env.Get("OTEL_ATTRIBUTE_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.SpanAttributeValueLengthLimit = env.Get("OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT",
		env.WithDefaultValue(uint(1024)),
	)
	c.SpanAttributeCountLimit = env.Get("OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.SpanEventCountLimit = env.Get("OTEL_SPAN_EVENT_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.SpanEventAttributeCountLimit = env.Get("OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.SpanLinkCountLimit = env.Get("OTEL_SPAN_LINK_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.SpanLinkAttributeCountLimit = env.Get("OTEL_LINK_ATTRIBUTE_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.LogRecordAttributeValueLengthLimit = env.Get("OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT",
		env.WithDefaultValue(uint(1024)),
	)
	c.LogRecordAttributeCountLimit = env.Get("OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT",
		env.WithDefaultValue(uint(128)),
	)
	c.TracesExporter = TracesExporter(env.Get("OTEL_TRACES_EXPORTER",
		env.WithDefaultValue(string(TracesExporterOTLP)),
		env.WithEnum(availableTracesExporters),
	))
	c.MetricsExporter = MetricsExporter(env.Get("OTEL_METRICS_EXPORTER",
		env.WithDefaultValue(string(MetricsExporterOTLP)),
		env.WithEnum(availableMetricsExporters),
	))
	c.LogsExporter = LogsExporter(env.Get("OTEL_LOGS_EXPORTER",
		env.WithDefaultValue(string(LogsExporterStdout)),
		env.WithEnum(availableLogsExporters),
	))
	c.MetricsExemplarFilter = MetricsExemplarFilter(env.Get("OTEL_METRICS_EXEMPLAR_FILTER",
		env.WithDefaultValue(string(MetricsExemplarFilterTraceBased)),
		env.WithEnum(availableMetricsExemplarFilters),
	))
	c.MetricExportInterval = env.Get("OTEL_METRIC_EXPORT_INTERVAL",
		env.WithDefaultValue(60*time.Second),
	)
	c.MetricExportTimeout = env.Get("OTEL_METRIC_EXPORT_TIMEOUT",
		env.WithDefaultValue(30*time.Second),
	)
	c.ProfilesExporter = ProfilesExporter(env.Get("OTEL_PROFILES_EXPORTER",
		env.WithDefaultValue(string(ProfilesExporterNone)),
		env.WithEnum(availableProfilesExporters),
	))
	c.ProfileExportInterval = env.Get("OTEL_PROFILE_EXPORT_INTERVAL",
		env.WithDefaultValue(10*time.Second),
	)
	c.OtlpExporterProtocol = OtlpExporterProtocol(env.Get("OTEL_EXPORTER_OTLP_PROTOCOL",
		env.WithDefaultValue(string(OtlpExporterProtocolGRPC)),
		env.WithEnum(availableOtlpExporterProtocols),
	))
	c.OtlpExporterEndpoint = env.Get("OTEL_EXPORTER_OTLP_ENDPOINT",
		env.WithDefaultValue("localhost:4317"),
	)
	c.OtlpTracesExporterEndpoint = env.Get("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT",
		env.WithDefaultValue("localhost:4317"),
	)
	c.OtlpMetricsExporterEndpoint = env.Get("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT",
		env.WithDefaultValue("localhost:4317"),
	)
	c.OtlpLogsExporterEndpoint = env.Get("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT",
		env.WithDefaultValue("localhost:4317"),
	)
	c.OtlpInsecure = env.Get("OTEL_EXPORTER_OTLP_INSECURE",
		env.WithDefaultValue(true),
	)
	c.OtlpTracesInsecure = env.Get("OTEL_EXPORTER_OTLP_TRACES_INSECURE",
		env.WithDefaultValue(true),
	)
	c.OtlpMetricsInsecure = env.Get("OTEL_EXPORTER_OTLP_METRICS_INSECURE",
		env.WithDefaultValue(true),
	)
	c.OtlpLogsInsecure = env.Get("OTEL_EXPORTER_OTLP_LOGS_INSECURE",
		env.WithDefaultValue(true),
	)
	c.PrettyPrint = env.Get("OTEL_PRETTY_PRINT",
		env.WithDefaultValue(false),
	)
	c.LogHandler = LogHandler(env.Get("OTEL_LOG_HANDLER",
		env.WithDefaultValue(string(LogHandlerZap)),
		env.WithEnum(availableLogHandlers),
	))

	if endpoint := env.Get[string]("OTEL_OTLP_ENDPOINT"); endpoint != "" {
		c.OtlpExporterEndpoint = endpoint
		c.OtlpTracesExporterEndpoint = endpoint
		c.OtlpMetricsExporterEndpoint = endpoint
		c.OtlpLogsExporterEndpoint = endpoint
	}

	if ratio := env.Get[float64]("OTEL_TRACE_SAMPLER"); ratio != 0 {
		c.TraceSampler = TraceSamplerParentBasedTraceIDRatio
		c.TraceSamplerArg = ratio
	}
}
