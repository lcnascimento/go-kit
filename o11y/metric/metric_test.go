package metric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	kitmetric "github.com/lcnascimento/go-kit/o11y/metric"
)

func TestMustDurationHistogram_RecordsSecondsInSemconvBuckets(t *testing.T) {
	t.Parallel()

	reader := metric.NewManualReader()
	meter := metric.NewMeterProvider(metric.WithReader(reader)).Meter("test")

	kitmetric.MustDurationHistogram(meter, "op.duration", "Duration of one op").Record(t.Context(), 0.2)

	var collected metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &collected))

	got := collected.ScopeMetrics[0].Metrics[0]
	require.Equal(t, "op.duration", got.Name)
	require.Equal(t, "s", got.Unit)

	point := got.Data.(metricdata.Histogram[float64]).DataPoints[0]
	require.Equal(t, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10}, point.Bounds)
	require.Equal(t, uint64(1), point.BucketCounts[6])
}
