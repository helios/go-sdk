package helioslogrus

import (
	"context"
	"io"
	"time"

	origin_logrus "github.com/sirupsen/logrus"
)

func Exit(code int) {
	origin_logrus.Exit(code)
}

func RegisterExitHandler(handler func()) {
	origin_logrus.RegisterExitHandler(handler)
}

func DeferExitHandler(handler func()) {
	origin_logrus.DeferExitHandler(handler)
}

func StandardLogger() *Logger {
	return origin_logrus.StandardLogger()
}

func SetOutput(out io.Writer) {
	origin_logrus.SetOutput(out)
}

func SetFormatter(formatter Formatter) {
	origin_logrus.SetFormatter(formatter)
}

func SetReportCaller(include bool) {
	origin_logrus.SetReportCaller(include)
}

func SetLevel(level Level) {
	origin_logrus.SetLevel(level)
}

func GetLevel() Level {
	return origin_logrus.GetLevel()
}

func IsLevelEnabled(level Level) bool {
	return origin_logrus.IsLevelEnabled(level)
}

func AddHook(hook Hook) {
	origin_logrus.AddHook(hook)
}

func WithError(err error) *Entry {
	return origin_logrus.WithError(err)
}

func WithContext(ctx context.Context) *Entry {
	return origin_logrus.WithContext(ctx)
}

func WithField(key string, value interface{}) *Entry {
	return origin_logrus.WithField(key, value)
}

func WithFields(fields Fields) *Entry {
	return origin_logrus.WithFields(fields)
}

func WithTime(t time.Time) *Entry {
	return origin_logrus.WithTime(t)
}

func Trace(args ...interface{}) {
	origin_logrus.Trace(args...)
}

func Debug(args ...interface{}) {
	origin_logrus.Debug(args...)
}

func Print(args ...interface{}) {
	origin_logrus.Print(args...)
}

func Info(args ...interface{}) {
	origin_logrus.Info(args...)
}

func Warn(args ...interface{}) {
	origin_logrus.Warn(args...)
}

func Warning(args ...interface{}) {
	origin_logrus.Warning(args...)
}

func Error(args ...interface{}) {
	origin_logrus.Error(args...)
}

func Panic(args ...interface{}) {
	origin_logrus.Panic(args...)
}

func Fatal(args ...interface{}) {
	origin_logrus.Fatal(args...)
}

func TraceFn(fn LogFunction) {
	origin_logrus.TraceFn(fn)
}

func DebugFn(fn LogFunction) {
	origin_logrus.DebugFn(fn)
}

func PrintFn(fn LogFunction) {
	origin_logrus.PrintFn(fn)
}

func InfoFn(fn LogFunction) {
	origin_logrus.InfoFn(fn)
}

func WarnFn(fn LogFunction) {
	origin_logrus.WarnFn(fn)
}

func WarningFn(fn LogFunction) {
	origin_logrus.WarningFn(fn)
}

func ErrorFn(fn LogFunction) {
	origin_logrus.ErrorFn(fn)
}

func PanicFn(fn LogFunction) {
	origin_logrus.PanicFn(fn)
}

func FatalFn(fn LogFunction) {
	origin_logrus.FatalFn(fn)
}

func Tracef(format string, args ...interface{}) {
	origin_logrus.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	origin_logrus.Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	origin_logrus.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	origin_logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	origin_logrus.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	origin_logrus.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	origin_logrus.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	origin_logrus.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	origin_logrus.Fatalf(format, args...)
}

func Traceln(args ...interface{}) {
	origin_logrus.Traceln(args...)
}

func Debugln(args ...interface{}) {
	origin_logrus.Debugln(args...)
}

func Println(args ...interface{}) {
	origin_logrus.Println(args...)
}

func Infoln(args ...interface{}) {
	origin_logrus.Infoln(args...)
}

func Warnln(args ...interface{}) {
	origin_logrus.Warnln(args...)
}

func Warningln(args ...interface{}) {
	origin_logrus.Warningln(args...)
}

func Errorln(args ...interface{}) {
	origin_logrus.Errorln(args...)
}

func Panicln(args ...interface{}) {
	origin_logrus.Panicln(args...)
}

func Fatalln(args ...interface{}) {
	origin_logrus.Fatalln(args...)
}

type TextFormatter = origin_logrus.TextFormatter

type BufferPool = origin_logrus.BufferPool

func SetBufferPool(bp BufferPool) {
	origin_logrus.SetBufferPool(bp)
}

type FieldMap = origin_logrus.FieldMap

type JSONFormatter = origin_logrus.JSONFormatter

type Fields = origin_logrus.Fields

type Level = origin_logrus.Level

func ParseLevel(lvl string) (Level, error) {
	return origin_logrus.ParseLevel(lvl)
}

var AllLevels = origin_logrus.AllLevels

const PanicLevel = origin_logrus.PanicLevel

const FatalLevel = origin_logrus.FatalLevel

const ErrorLevel = origin_logrus.ErrorLevel

const WarnLevel = origin_logrus.WarnLevel

const InfoLevel = origin_logrus.InfoLevel

const DebugLevel = origin_logrus.DebugLevel

const TraceLevel = origin_logrus.TraceLevel

type StdLogger = origin_logrus.StdLogger

type FieldLogger = origin_logrus.FieldLogger

type Ext1FieldLogger = origin_logrus.Ext1FieldLogger

var ErrorKey = origin_logrus.ErrorKey

type Entry = origin_logrus.Entry

func NewEntry(logger *Logger) *Entry {
	return origin_logrus.NewEntry(logger)
}

const FieldKeyMsg = origin_logrus.FieldKeyMsg

const FieldKeyLevel = origin_logrus.FieldKeyLevel

const FieldKeyTime = origin_logrus.FieldKeyTime

const FieldKeyLogrusError = origin_logrus.FieldKeyLogrusError

const FieldKeyFunc = origin_logrus.FieldKeyFunc

const FieldKeyFile = origin_logrus.FieldKeyFile

type Formatter = origin_logrus.Formatter

type LogFunction = origin_logrus.LogFunction

type Logger = origin_logrus.Logger

type MutexWrap = origin_logrus.MutexWrap

func New() *Logger {
	l := origin_logrus.New()
	l.AddHook(NewHook(WithLevels(origin_logrus.AllLevels...)))
	return l
}

type Hook = origin_logrus.Hook

type LevelHooks = origin_logrus.LevelHooks
