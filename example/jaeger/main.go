package main

import (
	"context"
	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "host.docker.internal:6831",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func TestDemo(req string, ctx context.Context) (reply string) {
	// 1. 创建span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo")
	defer func() {
		// 4. 接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	//2. 模拟耗时
	time.Sleep(time.Second / 2)
	//3. 返回reply
	reply = "TestDemoReply"
	return
}

// TestDemo2, 和上面TestDemo 逻辑代码一样
func TestDemo2(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo2")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	time.Sleep(time.Second / 2)
	reply = "TestDemo2Reply"
	return
}

func TestDemo3(req string, ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo3")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", "11111")
		span.Finish()
		TestDemo2("二次调用", ctx)
		TestDemo("二次调用", ctx)
	}()
	println(req)
	time.Sleep(time.Second * 10)
	return
}

func main() {
	tracer, closer := initJaeger("jager-test-demo")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)


	span := tracer.StartSpan("span_root")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r1 := TestDemo("Hello TestDemo", ctx)
	r2 := TestDemo2("Hello TestDemo2", ctx)
	go TestDemo3("Hello TestDemo3", ctx)
	fmt.Println(r1, r2)
	span.Finish()
	for {
		select {}
	}
}
