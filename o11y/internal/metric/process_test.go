package metric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	kitmetric "github.com/lcnascimento/go-kit/o11y/internal/metric"
)

func TestStartProcess_EmitsCPUTimeAndMemoryUsage(t *testing.T) {
	t.Parallel()

	reader := metric.NewManualReader()
	require.NoError(t, kitmetric.StartProcess(metric.NewMeterProvider(metric.WithReader(reader))))

	var collected metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &collected))

	got := map[string]metricdata.Metrics{}
	for _, scope := range collected.ScopeMetrics {
		for _, m := range scope.Metrics {
			got[m.Name] = m
		}
	}
	require.Len(t, got, 2)

	cpu := got["process.cpu.time"]
	require.Equal(t, "s", cpu.Unit)
	cpuPoints := cpu.Data.(metricdata.Sum[float64]).DataPoints
	cpuModes := make([]string, 0, len(cpuPoints))
	for _, point := range cpuPoints {
		mode, _ := point.Attributes.Value(attribute.Key("cpu.mode"))
		cpuModes = append(cpuModes, mode.AsString())
	}
	require.ElementsMatch(t, []string{"user", "system"}, cpuModes)

	mem := got["process.memory.usage"]
	require.Equal(t, "By", mem.Unit)
	points := mem.Data.(metricdata.Sum[int64]).DataPoints
	require.Len(t, points, 1)
	require.Positive(t, points[0].Value)
}
