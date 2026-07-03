package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/lcnascimento/go-kit/o11y/internal/config"
)

func TestConfig(t *testing.T) {
	_ = os.Setenv("OTEL_SERVICE_NAME", "o11y")
	_ = os.Setenv("OTEL_SERVICE_VERSION", "v0.1.0")
	_ = os.Setenv("OTEL_OTLP_ENDPOINT", "http://fake:4317")
	_ = os.Setenv("OTEL_TRACE_ENABLED", "true")
	_ = os.Setenv("OTEL_TRACES_EXPORTER", "stdout")
	_ = os.Setenv("OTEL_TRACE_SAMPLER", "0.5")
	_ = os.Setenv("OTEL_METRICS_ENABLED", "true")
	_ = os.Setenv("OTEL_METRICS_EXPORTER", "stdout")
	_ = os.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.version=v0.1.0,service.namespace=acquisition,deployment.environment=local")

	cfg := &config.Configuration{}
	cfg.Load()

	assert.Equal(t, "o11y", cfg.ServiceName)
	assert.Equal(t, "v0.1.0", cfg.ResourceAttributes[semconv.ServiceVersionKey])
	assert.Equal(t, "acquisition", cfg.ResourceAttributes[semconv.ServiceNamespaceKey])
	assert.Equal(t, "local", cfg.ResourceAttributes[semconv.DeploymentEnvironmentKey])
	assert.Equal(t, "http://fake:4317", cfg.OtlpExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpTracesExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpMetricsExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpLogsExporterEndpoint)
	assert.Equal(t, config.TracesExporterStdout, cfg.TracesExporter)
	assert.Equal(t, config.TraceSamplerParentBasedTraceIDRatio, cfg.TraceSampler)
	assert.Equal(t, 0.5, cfg.TraceSamplerArg)
	assert.Equal(t, config.MetricsExporterStdout, cfg.MetricsExporter)
	assert.Equal(t, config.MetricsExemplarFilterTraceBased, cfg.MetricsExemplarFilter)
	assert.Equal(t, time.Second*60, cfg.MetricExportInterval)
	assert.Equal(t, time.Second*30, cfg.MetricExportTimeout)
	assert.Equal(t, config.LogsExporterStdout, cfg.LogsExporter)

	assert.False(t, cfg.Disabled)
	assert.Equal(t, config.LogLevelInfo, cfg.LogLevel)
	assert.Equal(t, config.Propagators{
		config.PropagatorTraceContext: true,
		config.PropagatorBaggage:      true,
	}, cfg.Propagators)
	assert.Equal(t, time.Second*5, cfg.BspScheduleDelay)
	assert.Equal(t, time.Second*30, cfg.BspExportTimeout)
	assert.Equal(t, uint(5120), cfg.BspMaxQueueSize)
	assert.Equal(t, uint(5120), cfg.BspMaxExportBatchSize)
	assert.Equal(t, time.Second*1, cfg.BlrpScheduleDelay)
	assert.Equal(t, time.Second*30, cfg.BlrpExportTimeout)
	assert.Equal(t, uint(2048), cfg.BlrpMaxQueueSize)
	assert.Equal(t, uint(2048), cfg.BlrpMaxExportBatchSize)
	// assert.Equal(t, uint(128), cfg.AttributeValueLengthLimit)
	assert.Equal(t, uint(128), cfg.AttributeCountLimit)
	// assert.Equal(t, uint(128), cfg.SpanAttributeValueLengthLimit)
	assert.Equal(t, uint(128), cfg.SpanAttributeCountLimit)
	assert.Equal(t, uint(128), cfg.SpanEventCountLimit)
	assert.Equal(t, uint(128), cfg.SpanEventAttributeCountLimit)
	assert.Equal(t, uint(128), cfg.SpanLinkCountLimit)
	assert.Equal(t, uint(128), cfg.SpanLinkAttributeCountLimit)
	assert.Equal(t, uint(128), cfg.LogRecordAttributeCountLimit)
	// assert.Equal(t, uint(128), cfg.LogRecordAttributeValueLengthLimit)

	assert.Equal(t, config.OtlpExporterProtocolGRPC, cfg.OtlpExporterProtocol)
	assert.Equal(t, "http://fake:4317", cfg.OtlpExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpTracesExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpMetricsExporterEndpoint)
	assert.Equal(t, "http://fake:4317", cfg.OtlpLogsExporterEndpoint)

	assert.True(t, cfg.OtlpInsecure)
	assert.True(t, cfg.OtlpTracesInsecure)
	assert.True(t, cfg.OtlpMetricsInsecure)
	assert.True(t, cfg.OtlpLogsInsecure)

	assert.False(t, cfg.PrettyPrint)
	assert.Equal(t, config.LogHandlerZap, cfg.LogHandler)
	assert.Equal(t, config.ProfilesExporterNone, cfg.ProfilesExporter)
	assert.Equal(t, config.Profiles{
		config.ProfileCPU:  true,
		config.ProfileHeap: true,
	}, cfg.Profiles)
	assert.Equal(t, time.Second*10, cfg.ProfileExportInterval)
}
