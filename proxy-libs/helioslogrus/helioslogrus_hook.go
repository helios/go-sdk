package helioslogrus

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type heliosHook struct {
	levels []logrus.Level
}

// var _ logrus.Hook = (*Hook)(nil)

const hs_api_endpoint = "https://app.gethelios.dev"

// Option applies a configuration to the given config.
type option func(h *heliosHook)

// WithLevels sets the logrus logging levels on which the hook is fired.
//
// The default is all levels between logrus.PanicLevel and logrus.WarnLevel inclusive.
func withLevels(levels ...logrus.Level) option {
	return func(h *heliosHook) {
		h.levels = levels
	}
}

// NewHook returns a logrus hook.
func NewHook(opts ...option) *heliosHook {
	hook := &heliosHook{
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
func (hook *heliosHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	entry.Data["go_to_helios"] = fmt.Sprintf("%s?actionTraceId=%s&spanId=%s&source=logrus&timestamp=%s}", hs_api_endpoint, span.SpanContext().TraceID(), span.SpanContext().SpanID(), fmt.Sprint(time.Now().UnixNano()))
	return nil
}

// Levels returns logrus levels on which this hook is fired.
func (hook *heliosHook) Levels() []logrus.Level {
	return hook.levels
}
