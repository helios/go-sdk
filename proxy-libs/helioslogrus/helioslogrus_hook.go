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

const hsApiEndpoint = "https://app.gethelios.dev"


func AddHeliosHook() *heliosHook {
	hook := &heliosHook{
		levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}}

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
	entry.Data["go_to_helios"] = fmt.Sprintf("%s?actionTraceId=%s&spanId=%s&source=logrus&timestamp=%s", hsApiEndpoint, span.SpanContext().TraceID(), span.SpanContext().SpanID(), fmt.Sprint(time.Now().UnixNano()))
	return nil
}

// Levels returns logrus levels on which this hook is fired.
func (hook *heliosHook) Levels() []logrus.Level {
	return hook.levels
}
