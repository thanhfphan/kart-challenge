package tracing

import (
	"context"
	"net/http"

	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName    = "bitbucket.org/hasaki-tech/tracing"
	instrumentationVersion = "v0.1.0"
)

var tracer trace.Tracer

func init() {
	provider, err := apmotel.NewTracerProvider()
	if err != nil {
		panic(err)
	}

	otel.SetTracerProvider(provider)
	tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		// trace.WithSchemaURL(semconv.SchemaURL),
	)
}

func Version() string {
	return instrumentationVersion
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}

func WrapTransport(t http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(t)
}

func NewTransport() http.RoundTripper {
	return WrapTransport(http.DefaultTransport)
}
