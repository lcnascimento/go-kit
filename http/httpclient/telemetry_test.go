package httpclient_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/lcnascimento/go-kit/http/httpclient"
)

var metrics = sdkmetric.NewManualReader()

func TestMain(m *testing.M) {
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(metrics)))

	os.Exit(m.Run())
}

func TestClient_KeepsTheConcretePathOutOfTheMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	t.Cleanup(server.Close)

	_, err := httpclient.New().Get(t.Context(), &httpclient.Request{Host: server.URL, Path: "/notes/42"})
	require.NoError(t, err)

	var collected metricdata.ResourceMetrics
	require.NoError(t, metrics.Collect(t.Context(), &collected))

	recorded := map[string]bool{}

	for _, m := range collected.ScopeMetrics[0].Metrics {
		for _, set := range attributeSets(m.Data) {
			recorded[m.Name] = true
			require.False(t, set.HasValue("http.route"), m.Name)
		}
	}

	require.Equal(t, map[string]bool{
		"http.client.active_requests":  true,
		"http.client.request.duration": true,
		"http.client.request.total":    true,
	}, recorded)
}

func attributeSets(data metricdata.Aggregation) []attribute.Set {
	var sets []attribute.Set

	switch d := data.(type) {
	case metricdata.Sum[int64]:
		for _, p := range d.DataPoints {
			sets = append(sets, p.Attributes)
		}
	case metricdata.Histogram[float64]:
		for _, p := range d.DataPoints {
			sets = append(sets, p.Attributes)
		}
	}

	return sets
}
