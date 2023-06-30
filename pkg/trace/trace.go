package trace

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"net/http"
)

type Config struct {
	ServiceName string  `json:"service_name" yaml:"service_name"`
	Endpoint    string  `json:"endpoint" yaml:"endpoint"`
	SampleRate  float64 `json:"sample_rate" yaml:"sample_rate"`
}

func InitTrace(c Config) {
	var (
		exporter *jaeger.Exporter
		err      error
	)
	if c.Endpoint == "" {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint())
	} else {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Endpoint)))
	}
	if err != nil {
		panic(err)
	}
	opts := []trace.TracerProviderOption{
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(c.ServiceName),
		)),
	}
	if c.SampleRate >= 1 {
		opts = append(opts, trace.WithSampler(trace.AlwaysSample()))
	} else {
		if c.SampleRate <= 0 {
			c.SampleRate = 0.1
		}
		opts = append(opts, trace.WithSampler(trace.TraceIDRatioBased(c.SampleRate)))
	}
	provider := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	http.DefaultClient.Transport = otelhttp.NewTransport(http.DefaultTransport)
}
