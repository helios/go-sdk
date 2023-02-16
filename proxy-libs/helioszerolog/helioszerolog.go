package helioszerolog

import (
	"context"
	"io"

	origin_zerolog "github.com/rs/zerolog"
)
type Array = origin_zerolog.Array

func Arr() (*Array) {
	return origin_zerolog.Arr()
 }

type Level = origin_zerolog.Level

const DebugLevel = origin_zerolog.DebugLevel

const InfoLevel = origin_zerolog.InfoLevel

const WarnLevel = origin_zerolog.WarnLevel

const ErrorLevel = origin_zerolog.ErrorLevel

const FatalLevel = origin_zerolog.FatalLevel

const PanicLevel = origin_zerolog.PanicLevel

const NoLevel = origin_zerolog.NoLevel

const Disabled = origin_zerolog.Disabled

const TraceLevel = origin_zerolog.TraceLevel

func ParseLevel(levelStr string) (Level,error) {
	return origin_zerolog.ParseLevel(levelStr)
 }

func New(w io.Writer) Logger {
	return Logger{origin_zerolog.New(w)}
 }

func Nop() Logger {
	return Logger{origin_zerolog.Nop()}
 }

type LevelWriter = origin_zerolog.LevelWriter

func SyncWriter(w io.Writer) (io.Writer) {
	return origin_zerolog.SyncWriter(w)
 }

func MultiLevelWriter(writers ...io.Writer) (LevelWriter) {
	return origin_zerolog.MultiLevelWriter(writers...)
 }

type TestingLog = origin_zerolog.TestingLog

type TestWriter = origin_zerolog.TestWriter

func NewTestWriter(t TestingLog) (TestWriter) {
	return origin_zerolog.NewTestWriter(t)
 }

func ConsoleTestWriter(t TestingLog) (func(*ConsoleWriter) ) {
	return origin_zerolog.ConsoleTestWriter(t)
 }

func Ctx(ctx context.Context) (*Logger) {
	return &Logger{*origin_zerolog.Ctx(ctx)}
 }

type Hook = origin_zerolog.Hook

type HookFunc = origin_zerolog.HookFunc

type LevelHook = origin_zerolog.LevelHook

func NewLevelHook() (LevelHook) {
	return origin_zerolog.NewLevelHook()
 }

var Often = origin_zerolog.Often

var Sometimes = origin_zerolog.Sometimes

var Rarely = origin_zerolog.Rarely

type Sampler = origin_zerolog.Sampler

type RandomSampler = origin_zerolog.RandomSampler

type BasicSampler = origin_zerolog.BasicSampler

type BurstSampler = origin_zerolog.BurstSampler

type LevelSampler = origin_zerolog.LevelSampler

type Event = origin_zerolog.Event

type LogObjectMarshaler = origin_zerolog.LogObjectMarshaler

type LogArrayMarshaler = origin_zerolog.LogArrayMarshaler

func Dict() (*Event) {
	return origin_zerolog.Dict()
 }

const TimeFormatUnix = origin_zerolog.TimeFormatUnix

const TimeFormatUnixMs = origin_zerolog.TimeFormatUnixMs

const TimeFormatUnixMicro = origin_zerolog.TimeFormatUnixMicro

const TimeFormatUnixNano = origin_zerolog.TimeFormatUnixNano

var TimestampFieldName = origin_zerolog.TimestampFieldName

var LevelFieldName = origin_zerolog.LevelFieldName

var LevelTraceValue = origin_zerolog.LevelTraceValue

var LevelDebugValue = origin_zerolog.LevelDebugValue

var LevelInfoValue = origin_zerolog.LevelInfoValue

var LevelWarnValue = origin_zerolog.LevelWarnValue

var LevelErrorValue = origin_zerolog.LevelErrorValue

var LevelFatalValue = origin_zerolog.LevelFatalValue

var LevelPanicValue = origin_zerolog.LevelPanicValue

var LevelFieldMarshalFunc = origin_zerolog.LevelFieldMarshalFunc

var MessageFieldName = origin_zerolog.MessageFieldName

var ErrorFieldName = origin_zerolog.ErrorFieldName

var CallerFieldName = origin_zerolog.CallerFieldName

var CallerSkipFrameCount = origin_zerolog.CallerSkipFrameCount

var CallerMarshalFunc = origin_zerolog.CallerMarshalFunc

var ErrorStackFieldName = origin_zerolog.ErrorStackFieldName

var ErrorStackMarshaler = origin_zerolog.ErrorStackMarshaler

var ErrorMarshalFunc = origin_zerolog.ErrorMarshalFunc

var InterfaceMarshalFunc = origin_zerolog.InterfaceMarshalFunc

var TimeFieldFormat = origin_zerolog.TimeFieldFormat

var TimestampFunc = origin_zerolog.TimestampFunc

var DurationFieldUnit = origin_zerolog.DurationFieldUnit

var DurationFieldInteger = origin_zerolog.DurationFieldInteger

var ErrorHandler = origin_zerolog.ErrorHandler

var DefaultContextLogger = origin_zerolog.DefaultContextLogger

func SetGlobalLevel(l Level) {
	origin_zerolog.SetGlobalLevel(l)
 }

func GlobalLevel() (Level) {
	return origin_zerolog.GlobalLevel()
 }

func DisableSampling(v bool) {
	origin_zerolog.DisableSampling(v)
 }

type SyslogWriter = origin_zerolog.SyslogWriter

func SyslogLevelWriter(w SyslogWriter) (LevelWriter) {
	return origin_zerolog.SyslogLevelWriter(w)
 }

func SyslogCEEWriter(w SyslogWriter) (LevelWriter) {
	return origin_zerolog.SyslogCEEWriter(w)
 }

type Formatter = origin_zerolog.Formatter

type ConsoleWriter = origin_zerolog.ConsoleWriter

func NewConsoleWriter(options ...func(*ConsoleWriter) ) (ConsoleWriter) {
	return origin_zerolog.NewConsoleWriter(options...)
 }

type Context = origin_zerolog.Context

