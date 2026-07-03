package temporalclient

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"google.golang.org/protobuf/types/known/durationpb"

	enumspb "go.temporal.io/api/enums/v1"

	"github.com/lcnascimento/go-kit/env"

	"github.com/lcnascimento/go-kit/temporal/internal"
)

var (
	module = "github.com/lcnascimento/go-kit/temporal"
	logger = internal.NewLogger()
	meter  = otel.Meter(module)
)

const retentionPeriod = 14 * 24 * time.Hour // 14 days

func New(opts ...Option) (client.Client, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	host := env.Get("TEMPORAL_HOST", env.WithDefaultValue("127.0.0.1:7233"))
	namespace := env.Get("TEMPORAL_NAMESPACE", env.WithDefaultValue("default"))

	nsClient, err := client.NewNamespaceClient(client.Options{HostPort: host})
	if err != nil {
		return nil, err
	}

	var alreadyExists *serviceerror.NamespaceAlreadyExists

	err = nsClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: durationpb.New(retentionPeriod),
	})
	if err != nil && !errors.As(err, &alreadyExists) {
		return nil, err
	}

	var metricsHandler client.MetricsHandler
	if o.noopMetrics {
		metricsHandler = client.MetricsNopHandler
	} else {
		metricsHandler = opentelemetry.NewMetricsHandler(opentelemetry.MetricsHandlerOptions{
			Meter: meter,
		})
	}

	clientOpts := client.Options{
		HostPort:           host,
		Namespace:          namespace,
		Logger:             logger,
		MetricsHandler:     metricsHandler,
		Interceptors:       Interceptors(opts...),
		ContextPropagators: Propagators(),
		FailureConverter:   NewFailureConverter(),
	}

	cli, err := client.Dial(clientOpts)
	if err != nil {
		return nil, err
	}

	ensureSearchAttributes(context.Background(), cli, namespace)

	return cli, nil
}

func ensureSearchAttributes(ctx context.Context, cli client.Client, namespace string) {
	_, err := cli.OperatorService().AddSearchAttributes(ctx, &operatorservice.AddSearchAttributesRequest{
		Namespace: namespace,
		SearchAttributes: map[string]enumspb.IndexedValueType{
			"user_id":      enumspb.INDEXED_VALUE_TYPE_KEYWORD,
			"organization": enumspb.INDEXED_VALUE_TYPE_KEYWORD,
			"account_id":   enumspb.INDEXED_VALUE_TYPE_KEYWORD,
		},
	})
	if err != nil {
		onWarnSearchAttributes(ctx, namespace, err)
	}
}

func MustNew() client.Client {
	cli, err := New()
	if err != nil {
		panic(err)
	}

	return cli
}
