package helioszerolog

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	origin_zerolog "github.com/rs/zerolog"
)
const hsApiEndpoint = "https://app.gethelios.dev"

func extractDataFromContext(otelContext context.Context, event *Event) {
	if otelContext != nil {
		span := trace.SpanFromContext(otelContext)
		if span.IsRecording() {
			traceId := span.SpanContext().TraceID().String()
			spanId := span.SpanContext().SpanID().String()
			event.Str("spanId", spanId)
			event.Str("traceId", traceId)
			event.Str("go_to_helios", fmt.Sprintf("%s?actionTraceId=%s&spanId=%s&source=zerolog&timestamp=%s",hsApiEndpoint, traceId, spanId, fmt.Sprint(time.Now().UnixNano())))
		}
	}
}

type Logger struct {
	WrappedLogger *origin_zerolog.Logger
	otelContext   *context.Context
}

func (l *Logger) UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	l.WrappedLogger.UpdateContext(update)
}

func (l *Logger) Trace() *Event {
	event := l.WrappedLogger.Trace()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Debug() *Event {
	event := l.WrappedLogger.Debug()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Info() *Event {
	event := l.WrappedLogger.Info()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Warn() *Event {
	event := l.WrappedLogger.Warn()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Error() *Event {
	event := l.WrappedLogger.Error()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Err(err error) *Event {
	event := l.WrappedLogger.Err(err)
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Fatal() *Event {
	event := l.WrappedLogger.Fatal()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) Panic() *Event {
	event := l.WrappedLogger.Panic()
	extractDataFromContext(*l.otelContext, event)
	return event
}

func (l *Logger) WithLevel(level Level) *Event {
	return l.WrappedLogger.WithLevel(level)
}

func (l *Logger) Log() *Event {
	event := l.WrappedLogger.Log()
	extractDataFromContext(*l.otelContext, event)
	return event
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.WrappedLogger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.WrappedLogger.Printf(format, v...)
}

func (l Logger) Output(w io.Writer) Logger {
	internalLogger := l.WrappedLogger.Output(w)
	return Logger{&internalLogger, l.otelContext}
}

func (l Logger) With() Context {
	return l.WrappedLogger.With()
}

func (l Logger) Level(lvl Level) Logger {
	internalLogger := l.WrappedLogger.Level(lvl)
	return Logger{&internalLogger, l.otelContext}
}

func (l Logger) GetLevel() Level {
	return l.WrappedLogger.GetLevel()
}

func (l Logger) Sample(s Sampler) Logger {
	internalLogger := l.WrappedLogger.Sample(s)
	return Logger{&internalLogger, l.otelContext}
}

func (l Logger) Hook(h Hook) Logger {
	internalLogger := l.WrappedLogger.Hook(h)
	return Logger{&internalLogger, l.otelContext}
}

func (l Logger) Write(p []byte) (n int, err error) {
	return l.WrappedLogger.Write(p)
}

type ctxKey struct{}

func (l Logger) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKey{}).(*Logger); !ok && l.WrappedLogger.GetLevel() == Disabled {
		// Do not store disabled logger.
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, l.WrappedLogger)
}