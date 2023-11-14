package trace

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"im/common/config"
	"im/logger"
	"sync"
)

var (
	tp   *tracesdk.TracerProvider
	once sync.Once
)

/*TODO:该模块负责实现PRPC的链路追踪功能*/

func StartAgent() {
	once.Do(func() {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.GetPrpcTraceUrl())))
		if err != nil {
			logger.Errorf("PRPC 链路追踪代理模块启动失败: error:%v", err.Error())
			return
		}
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exp),
			tracesdk.WithSampler(tracesdk.TraceIDRatioBased(config.GetPrpcTraceSampler())),
			tracesdk.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.GetPrpcTraceServiceName()),
			)),
		)

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		))
	})
}

func StopAgent() {
	/*关闭链路追踪*/
	_ = tp.Shutdown(context.TODO())
}
