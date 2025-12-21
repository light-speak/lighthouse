package extensions

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/light-speak/lighthouse/metrics"
)

type MetricsExtension struct{}

func (e *MetricsExtension) ExtensionName() string {
	return "PrometheusMetrics"
}

func (e *MetricsExtension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// 每个 GraphQL Operation（query / mutation）
func (e *MetricsExtension) InterceptOperation(
	ctx context.Context,
	next graphql.OperationHandler,
) graphql.ResponseHandler {

	opCtx := graphql.GetOperationContext(ctx)
	if opCtx != nil {
		metrics.GQLOperationTotal.WithLabelValues(
			opCtx.OperationName,
			string(opCtx.Operation.Operation),
		).Inc()
	}

	return next(ctx)
}

// 每个 resolver（重点）
func (e *MetricsExtension) InterceptField(
	ctx context.Context,
	next graphql.Resolver,
	info *graphql.FieldContext,
) (res interface{}, err error) {

	start := time.Now()
	res, err = next(ctx)

	metrics.GQLResolverDuration.
		WithLabelValues(info.Object, info.Field.Name).
		Observe(time.Since(start).Seconds())

	return res, err
}
