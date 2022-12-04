package heliosgin

import (
	"io"
	"net/http"

	origin_gin "github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Param = origin_gin.Param

type Params = origin_gin.Params

type RecoveryFunc = origin_gin.RecoveryFunc

func Recovery() HandlerFunc {
	return origin_gin.Recovery()
}

func CustomRecovery(handle RecoveryFunc) HandlerFunc {
	return origin_gin.CustomRecovery(handle)
}

func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	return origin_gin.RecoveryWithWriter(out, recovery...)
}

func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
	return origin_gin.CustomRecoveryWithWriter(out, handle)
}

func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	return origin_gin.CreateTestContext(w)
}

func IsDebugging() bool {
	return origin_gin.IsDebugging()
}

var DebugPrintRouteFunc = origin_gin.DebugPrintRouteFunc

type HandlerFunc = origin_gin.HandlerFunc

type HandlersChain = origin_gin.HandlersChain

type RouteInfo = origin_gin.RouteInfo

type RoutesInfo = origin_gin.RoutesInfo

const PlatformGoogleAppEngine = origin_gin.PlatformGoogleAppEngine

const PlatformCloudflare = origin_gin.PlatformCloudflare

type Engine = origin_gin.Engine

func New() *Engine {
	engine := origin_gin.New()
	engine.Use(otelgin.Middleware("opentelemetry-middleware"))
	return engine
}

func Default() *Engine {
	return origin_gin.Default()
}

const MIMEJSON = origin_gin.MIMEJSON

const MIMEHTML = origin_gin.MIMEHTML

const MIMEXML = origin_gin.MIMEXML

const MIMEXML2 = origin_gin.MIMEXML2

const MIMEPlain = origin_gin.MIMEPlain

const MIMEPOSTForm = origin_gin.MIMEPOSTForm

const MIMEMultipartPOSTForm = origin_gin.MIMEMultipartPOSTForm

const MIMEYAML = origin_gin.MIMEYAML

const MIMETOML = origin_gin.MIMETOML

const BodyBytesKey = origin_gin.BodyBytesKey

const ContextKey = origin_gin.ContextKey

type Context = origin_gin.Context

type Negotiate = origin_gin.Negotiate

type ErrorType = origin_gin.ErrorType

const ErrorTypeBind = origin_gin.ErrorTypeBind

const ErrorTypeRender = origin_gin.ErrorTypeRender

const ErrorTypePrivate = origin_gin.ErrorTypePrivate

const ErrorTypePublic = origin_gin.ErrorTypePublic

const ErrorTypeAny = origin_gin.ErrorTypeAny

const ErrorTypeNu = origin_gin.ErrorTypeNu

type Error = origin_gin.Error

type LoggerConfig = origin_gin.LoggerConfig

type LogFormatter = origin_gin.LogFormatter

type LogFormatterParams = origin_gin.LogFormatterParams

func DisableConsoleColor() {
	origin_gin.DisableConsoleColor()
}

func ForceConsoleColor() {
	origin_gin.ForceConsoleColor()
}

func ErrorLogger() HandlerFunc {
	return origin_gin.ErrorLogger()
}

func ErrorLoggerT(typ ErrorType) HandlerFunc {
	return origin_gin.ErrorLoggerT(typ)
}

func Logger() HandlerFunc {
	return origin_gin.Logger()
}

func LoggerWithFormatter(f LogFormatter) HandlerFunc {
	return origin_gin.LoggerWithFormatter(f)
}

func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	return origin_gin.LoggerWithWriter(out, notlogged...)
}

func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	return origin_gin.LoggerWithConfig(conf)
}

const EnvGinMode = origin_gin.EnvGinMode

const DebugMode = origin_gin.DebugMode

const ReleaseMode = origin_gin.ReleaseMode

const TestMode = origin_gin.TestMode

var DefaultWriter = origin_gin.DefaultWriter

var DefaultErrorWriter = origin_gin.DefaultErrorWriter

func SetMode(value string) {
	origin_gin.SetMode(value)
}

func DisableBindValidation() {
	origin_gin.DisableBindValidation()
}

func EnableJsonDecoderUseNumber() {
	origin_gin.EnableJsonDecoderUseNumber()
}

func EnableJsonDecoderDisallowUnknownFields() {
	origin_gin.EnableJsonDecoderDisallowUnknownFields()
}

func Mode() string {
	return origin_gin.Mode()
}

type ResponseWriter = origin_gin.ResponseWriter

type IRouter = origin_gin.IRouter

type IRoutes = origin_gin.IRoutes

type RouterGroup = origin_gin.RouterGroup

const BindKey = origin_gin.BindKey

func Bind(val any) HandlerFunc {
	return origin_gin.Bind(val)
}

func WrapF(f http.HandlerFunc) HandlerFunc {
	return origin_gin.WrapF(f)
}

func WrapH(h http.Handler) HandlerFunc {
	return origin_gin.WrapH(h)
}

type H = origin_gin.H

const Version = origin_gin.Version

const AuthUserKey = origin_gin.AuthUserKey

type Accounts = origin_gin.Accounts

func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	return origin_gin.BasicAuthForRealm(accounts, realm)
}

func BasicAuth(accounts Accounts) HandlerFunc {
	return origin_gin.BasicAuth(accounts)
}

func Dir(root string, listDirectory bool) http.FileSystem {
	return origin_gin.Dir(root, listDirectory)
}
