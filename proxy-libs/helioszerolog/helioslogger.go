package helioszerolog

import (
	"context"
	"io"

	"github.com/rs/zerolog"

	origin_zerolog "github.com/rs/zerolog"
)

type Logger struct {
	realLogger origin_zerolog.Logger
}

func (l *Logger) UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	l.realLogger.UpdateContext(update)
}

func (l *Logger) Trace(ctx context.Context) *Event {
	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("traceparent", ctx.Value("traceparent").(string))
	})
	return l.realLogger.wTrace()
}

func (l *Logger) Debug() *Event {
	return l.realLogger.Debug()
}

func (l *Logger) Info() *Event {
	return l.realLogger.Debug()
}

func (l *Logger) Warn() *Event {
	return l.realLogger.Warn()
}

func (l *Logger) Error() *Event {
	return l.realLogger.Error()
}

func (l *Logger) Err(err error) *Event {
	return l.realLogger.Err(err)
}

func (l *Logger) Fatal() *Event {
	return l.realLogger.Fatal()
}

func (l *Logger) Panic() *Event {
	return l.realLogger.Panic()
}

func (l *Logger) WithLevel(level Level) *Event {
	return l.realLogger.WithLevel(level)
}

func (l *Logger) Log() *Event {
	return l.realLogger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.realLogger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.realLogger.Printf(format, v...)
}

func (l Logger) Output(w io.Writer) Logger {
	return l.realLogger.Output(w)
}

func (l Logger) With() Context {
	return l.With()
}

func (l Logger) Level(lvl Level) Logger {
	return l.realLogger.Level(lvl)
}

func (l Logger) GetLevel() Level {
	return l.realLogger.GetLevel()
}

func (l Logger) Sample(s Sampler) Logger {
	return l.realLogger.Sample(s)
}

func (l Logger) Hook(h Hook) Logger {
	return l.realLogger.Hook(h)
}

func (l Logger) Write(p []byte) (n int, err error) {
	return l.realLogger.Write(p)
}
