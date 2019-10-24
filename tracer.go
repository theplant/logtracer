/*
# Use opentracing-go api to print pretty structural trace log

```
15:56:15.45 logtracer operation=TestLogTracer:LongFunc1:DatabaseOp2 duration=105ms start=2019-10-24T15:56:15.34693+08:00
15:56:15.75 logtracer operation=TestLogTracer:LongFunc1 duration=405ms start=2019-10-24T15:56:15.346908+08:00
15:56:15.75 logtracer operation=TestLogTracer duration=608ms start=2019-10-24T15:56:15.1444+08:00
```

## Example

*/
package logtracer

import (
	"context"
	"fmt"
	"time"

	"github.com/theplant/appkit/log"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

type LogTracer struct {
	mocktracer.MockTracer
}

func New() opentracing.Tracer {
	mt := mocktracer.New()
	return &LogTracer{MockTracer: *mt}
}

func (t *LogTracer) ContextWithSpanHook(ctx context.Context, span opentracing.Span) context.Context {
	span.(*LogSpan).ctx = ctx
	return ctx
}

// StartSpan belongs to the Tracer interface.
func (t *LogTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return newMockSpan(t, operationName, opts...)
}

const operationBaggageKey = "log_span_operation"

func newMockSpan(t *LogTracer, name string, opts ...opentracing.StartSpanOption) *LogSpan {
	var parentOperationName string
	for _, op := range opts {
		if sr, ok := op.(opentracing.SpanReference); ok {
			if sr.Type == opentracing.ChildOfRef {
				parentOperationName = sr.ReferencedContext.(mocktracer.MockSpanContext).Baggage[operationBaggageKey]
			}
		}
	}
	if len(parentOperationName) > 0 {
		parentOperationName = parentOperationName + ":"
	}

	spanName := fmt.Sprintf("%s%s", parentOperationName, name)
	mock := t.MockTracer.StartSpan(spanName, opts...)
	mtspan := mock.(*mocktracer.MockSpan)
	mtspan.SpanContext = mtspan.SpanContext.WithBaggageItem(operationBaggageKey, spanName)
	r := &LogSpan{
		MockSpan: *mtspan,
	}

	return r
}

type LogSpan struct {
	ctx context.Context
	mocktracer.MockSpan
}

func (span LogSpan) Finish() {
	span.MockSpan.Finish()
	printlog(&span)
}

func printlog(span *LogSpan) {
	l := log.ForceContext(span.ctx)
	_ = l.Info().Log(
		"msg", "logtracer",
		"operation", span.MockSpan.OperationName,
		"duration", fmt.Sprintf("%dms", span.MockSpan.FinishTime.Sub(span.MockSpan.StartTime)*1.0/1000000),
		"start", span.MockSpan.StartTime.Format(time.RFC3339Nano),
	)
}

func (span LogSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	span.MockSpan.FinishWithOptions(opts)
	printlog(&span)
}
