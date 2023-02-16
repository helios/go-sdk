package helioslogrus

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Hook struct {
	levels []logrus.Level
}

var _ logrus.Hook = (*Hook)(nil)

const HS_API_ENDPOINT = "https://app.gethelios.dev"

// Option applies a configuration to the given config.
type Option func(h *Hook)

// WithLevels sets the logrus logging levels on which the hook is fired.
//
// The default is all levels between logrus.PanicLevel and logrus.WarnLevel inclusive.
func WithLevels(levels ...logrus.Level) Option {
	return func(h *Hook) {
		h.levels = levels
	}
}

// NewHook returns a logrus hook.
func NewHook(opts ...Option) *Hook {
	hook := &Hook{
		levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}}

	for _, fn := range opts {
		fn(hook)
	}

	return hook
}

// Fire is a logrus hook that is fired on a new log entry.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	entry.Data["go_to_helios"] = fmt.Sprintf("%s?actionTraceId=%s&spanId=%s", HS_API_ENDPOINT, span.SpanContext().TraceID(), span.SpanContext().SpanID())
	return nil
}

// Levels returns logrus levels on which this hook is fired.
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}
