package logtracer_test

import (
	"context"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/theplant/appkit/log"
	"github.com/theplant/logtracer"
)

func ExampleNew_01() {
	var DatabaseOp2 = func(ctx context.Context) {
		var span opentracing.Span
		span, ctx = opentracing.StartSpanFromContext(ctx, "DatabaseOp2")
		defer span.Finish()
		time.Sleep(100 * time.Millisecond)
	}

	var LongFunc1 = func(ctx context.Context) {
		var span opentracing.Span
		span, ctx = opentracing.StartSpanFromContext(ctx, "LongFunc1")
		defer span.Finish()

		DatabaseOp2(ctx)
		time.Sleep(300 * time.Millisecond)
	}

	tracer := logtracer.New()
	opentracing.SetGlobalTracer(tracer)

	ctx := context.TODO()
	ctx = log.Context(ctx, log.Default())

	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "TestLogTracer")
	defer span.Finish()

	time.Sleep(200 * time.Millisecond)
	LongFunc1(ctx)

}

func TestTrace(t *testing.T) {
	ExampleNew_01()
}
