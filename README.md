

# Use opentracing-go api to print pretty structural trace log

```
15:56:15.45 logtracer operation=TestLogTracer:LongFunc1:DatabaseOp2 duration=105ms start=2019-10-24T15:56:15.34693+08:00
15:56:15.75 logtracer operation=TestLogTracer:LongFunc1 duration=405ms start=2019-10-24T15:56:15.346908+08:00
15:56:15.75 logtracer operation=TestLogTracer duration=608ms start=2019-10-24T15:56:15.1444+08:00
```

## Example



```go
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
```



