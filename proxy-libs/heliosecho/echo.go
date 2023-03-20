package heliosecho

import (
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/helios/opentelemetry-go-contrib/instrumentation/github.com/labstack/echo/otelecho"
	origin_echo "github.com/labstack/echo/v4"
)

type Context = origin_echo.Context

const ContextKeyHeaderAllow = origin_echo.ContextKeyHeaderAllow

type Echo = origin_echo.Echo

type Route = origin_echo.Route

type HTTPError = origin_echo.HTTPError

type MiddlewareFunc = origin_echo.MiddlewareFunc

type HandlerFunc = origin_echo.HandlerFunc

type HTTPErrorHandler = origin_echo.HTTPErrorHandler

type Validator = origin_echo.Validator

type JSONSerializer = origin_echo.JSONSerializer

type Renderer = origin_echo.Renderer

type Map = origin_echo.Map

const CONNECT = origin_echo.CONNECT

const DELETE = origin_echo.DELETE

const GET = origin_echo.GET

const HEAD = origin_echo.HEAD

const OPTIONS = origin_echo.OPTIONS

const PATCH = origin_echo.PATCH

const POST = origin_echo.POST

const PUT = origin_echo.PUT

const TRACE = origin_echo.TRACE

const MIMEApplicationJSON = origin_echo.MIMEApplicationJSON

const MIMEApplicationJSONCharsetUTF8 = origin_echo.MIMEApplicationJSONCharsetUTF8

const MIMEApplicationJavaScript = origin_echo.MIMEApplicationJavaScript

const MIMEApplicationJavaScriptCharsetUTF8 = origin_echo.MIMEApplicationJavaScriptCharsetUTF8

const MIMEApplicationXML = origin_echo.MIMEApplicationXML

const MIMEApplicationXMLCharsetUTF8 = origin_echo.MIMEApplicationXMLCharsetUTF8

const MIMETextXML = origin_echo.MIMETextXML

const MIMETextXMLCharsetUTF8 = origin_echo.MIMETextXMLCharsetUTF8

const MIMEApplicationForm = origin_echo.MIMEApplicationForm

const MIMEApplicationProtobuf = origin_echo.MIMEApplicationProtobuf

const MIMEApplicationMsgpack = origin_echo.MIMEApplicationMsgpack

const MIMETextHTML = origin_echo.MIMETextHTML

const MIMETextHTMLCharsetUTF8 = origin_echo.MIMETextHTMLCharsetUTF8

const MIMETextPlain = origin_echo.MIMETextPlain

const MIMETextPlainCharsetUTF8 = origin_echo.MIMETextPlainCharsetUTF8

const MIMEMultipartForm = origin_echo.MIMEMultipartForm

const MIMEOctetStream = origin_echo.MIMEOctetStream

const PROPFIND = origin_echo.PROPFIND

const REPORT = origin_echo.REPORT

const RouteNotFound = origin_echo.RouteNotFound

const HeaderAccept = origin_echo.HeaderAccept

const HeaderAcceptEncoding = origin_echo.HeaderAcceptEncoding

const HeaderAllow = origin_echo.HeaderAllow

const HeaderAuthorization = origin_echo.HeaderAuthorization

const HeaderContentDisposition = origin_echo.HeaderContentDisposition

const HeaderContentEncoding = origin_echo.HeaderContentEncoding

const HeaderContentLength = origin_echo.HeaderContentLength

const HeaderContentType = origin_echo.HeaderContentType

const HeaderCookie = origin_echo.HeaderCookie

const HeaderSetCookie = origin_echo.HeaderSetCookie

const HeaderIfModifiedSince = origin_echo.HeaderIfModifiedSince

const HeaderLastModified = origin_echo.HeaderLastModified

const HeaderLocation = origin_echo.HeaderLocation

const HeaderRetryAfter = origin_echo.HeaderRetryAfter

const HeaderUpgrade = origin_echo.HeaderUpgrade

const HeaderVary = origin_echo.HeaderVary

const HeaderWWWAuthenticate = origin_echo.HeaderWWWAuthenticate

const HeaderXForwardedFor = origin_echo.HeaderXForwardedFor

const HeaderXForwardedProto = origin_echo.HeaderXForwardedProto

const HeaderXForwardedProtocol = origin_echo.HeaderXForwardedProtocol

const HeaderXForwardedSsl = origin_echo.HeaderXForwardedSsl

const HeaderXUrlScheme = origin_echo.HeaderXUrlScheme

const HeaderXHTTPMethodOverride = origin_echo.HeaderXHTTPMethodOverride

const HeaderXRealIP = origin_echo.HeaderXRealIP

const HeaderXRequestID = origin_echo.HeaderXRequestID

const HeaderXCorrelationID = origin_echo.HeaderXCorrelationID

const HeaderXRequestedWith = origin_echo.HeaderXRequestedWith

const HeaderServer = origin_echo.HeaderServer

const HeaderOrigin = origin_echo.HeaderOrigin

const HeaderCacheControl = origin_echo.HeaderCacheControl

const HeaderConnection = origin_echo.HeaderConnection

const HeaderAccessControlRequestMethod = origin_echo.HeaderAccessControlRequestMethod

const HeaderAccessControlRequestHeaders = origin_echo.HeaderAccessControlRequestHeaders

const HeaderAccessControlAllowOrigin = origin_echo.HeaderAccessControlAllowOrigin

const HeaderAccessControlAllowMethods = origin_echo.HeaderAccessControlAllowMethods

const HeaderAccessControlAllowHeaders = origin_echo.HeaderAccessControlAllowHeaders

const HeaderAccessControlAllowCredentials = origin_echo.HeaderAccessControlAllowCredentials

const HeaderAccessControlExposeHeaders = origin_echo.HeaderAccessControlExposeHeaders

const HeaderAccessControlMaxAge = origin_echo.HeaderAccessControlMaxAge

const HeaderStrictTransportSecurity = origin_echo.HeaderStrictTransportSecurity

const HeaderXContentTypeOptions = origin_echo.HeaderXContentTypeOptions

const HeaderXXSSProtection = origin_echo.HeaderXXSSProtection

const HeaderXFrameOptions = origin_echo.HeaderXFrameOptions

const HeaderContentSecurityPolicy = origin_echo.HeaderContentSecurityPolicy

const HeaderContentSecurityPolicyReportOnly = origin_echo.HeaderContentSecurityPolicyReportOnly

const HeaderXCSRFToken = origin_echo.HeaderXCSRFToken

const HeaderReferrerPolicy = origin_echo.HeaderReferrerPolicy

const Version = origin_echo.Version

var ErrUnsupportedMediaType = origin_echo.ErrUnsupportedMediaType

var ErrNotFound = origin_echo.ErrNotFound

var ErrUnauthorized = origin_echo.ErrUnauthorized

var ErrForbidden = origin_echo.ErrForbidden

var ErrMethodNotAllowed = origin_echo.ErrMethodNotAllowed

var ErrStatusRequestEntityTooLarge = origin_echo.ErrStatusRequestEntityTooLarge

var ErrTooManyRequests = origin_echo.ErrTooManyRequests

var ErrBadRequest = origin_echo.ErrBadRequest

var ErrBadGateway = origin_echo.ErrBadGateway

var ErrInternalServerError = origin_echo.ErrInternalServerError

var ErrRequestTimeout = origin_echo.ErrRequestTimeout

var ErrServiceUnavailable = origin_echo.ErrServiceUnavailable

var ErrValidatorNotRegistered = origin_echo.ErrValidatorNotRegistered

var ErrRendererNotRegistered = origin_echo.ErrRendererNotRegistered

var ErrInvalidRedirectCode = origin_echo.ErrInvalidRedirectCode

var ErrCookieNotFound = origin_echo.ErrCookieNotFound

var ErrInvalidCertOrKeyType = origin_echo.ErrInvalidCertOrKeyType

var ErrInvalidListenerNetwork = origin_echo.ErrInvalidListenerNetwork

var NotFoundHandler = origin_echo.NotFoundHandler

var MethodNotAllowedHandler = origin_echo.MethodNotAllowedHandler

func New() (e *Echo) {
	echo := origin_echo.New()
	if os.Getenv("HS_DISABLED") != "true" {
		echo.Use(otelecho.Middleware("opentelemetry-middleware"))
	}
	return echo
}

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	return origin_echo.NewHTTPError(code, message...)
}

func WrapHandler(h http.Handler) HandlerFunc {
	return origin_echo.WrapHandler(h)
}

func WrapMiddleware(m func(http.Handler) http.Handler) MiddlewareFunc {
	return origin_echo.WrapMiddleware(m)
}

func GetPath(r *http.Request) string {
	return origin_echo.GetPath(r)
}

type TrustOption = origin_echo.TrustOption

func TrustLoopback(v bool) TrustOption {
	return origin_echo.TrustLoopback(v)
}

func TrustLinkLocal(v bool) TrustOption {
	return origin_echo.TrustLinkLocal(v)
}

func TrustPrivateNet(v bool) TrustOption {
	return origin_echo.TrustPrivateNet(v)
}

func TrustIPRange(ipRange *net.IPNet) TrustOption {
	return origin_echo.TrustIPRange(ipRange)
}

type IPExtractor = origin_echo.IPExtractor

func ExtractIPDirect() IPExtractor {
	return origin_echo.ExtractIPDirect()
}

func ExtractIPFromRealIPHeader(options ...TrustOption) IPExtractor {
	return origin_echo.ExtractIPFromRealIPHeader(options...)
}

func ExtractIPFromXFFHeader(options ...TrustOption) IPExtractor {
	return origin_echo.ExtractIPFromXFFHeader(options...)
}

type Binder = origin_echo.Binder

type DefaultBinder = origin_echo.DefaultBinder

type BindUnmarshaler = origin_echo.BindUnmarshaler

type Group = origin_echo.Group

type DefaultJSONSerializer = origin_echo.DefaultJSONSerializer

type Logger = origin_echo.Logger

type Response = origin_echo.Response

func NewResponse(w http.ResponseWriter, e *Echo) (r *Response) {
	return origin_echo.NewResponse(w, e)
}

type BindingError = origin_echo.BindingError

func NewBindingError(sourceParam string, values []string, message interface{}, internalError error) error {
	return origin_echo.NewBindingError(sourceParam, values, message, internalError)
}

type ValueBinder = origin_echo.ValueBinder

func QueryParamsBinder(c Context) *ValueBinder {
	return origin_echo.QueryParamsBinder(c)
}

func PathParamsBinder(c Context) *ValueBinder {
	return origin_echo.PathParamsBinder(c)
}

func FormFieldBinder(c Context) *ValueBinder {
	return origin_echo.FormFieldBinder(c)
}

type Router = origin_echo.Router

func NewRouter(e *Echo) *Router {
	return origin_echo.NewRouter(e)
}

func StaticDirectoryHandler(fileSystem fs.FS, disablePathUnescaping bool) HandlerFunc {
	return origin_echo.StaticDirectoryHandler(fileSystem, disablePathUnescaping)
}

func StaticFileHandler(file string, filesystem fs.FS) HandlerFunc {
	return origin_echo.StaticFileHandler(file, filesystem)
}

func MustSubFS(currentFs fs.FS, fsRoot string) fs.FS {
	return origin_echo.MustSubFS(currentFs, fsRoot)
}
