package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const ServerName = "kratos-shop"

// https://github.com/go-kratos/examples/blob/main/traces/app/user/main.go
func InitTracer(url string) (*tracesdk.TracerProvider, error) {
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))), //
		// 始终确保再生成中批量处理
		tracesdk.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("kratos-trace"),
			attribute.String("exporter", "jaeger"),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp, err
}
