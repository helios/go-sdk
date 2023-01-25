package helioshttp

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"net"
	realHttp "net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrBodyNotAllowed = realHttp.ErrBodyNotAllowed

var ErrHijacked = realHttp.ErrHijacked

var ErrContentLength = realHttp.ErrContentLength

var ErrWriteAfterFlush = realHttp.ErrWriteAfterFlush

type Handler = realHttp.Handler

type ResponseWriter = realHttp.ResponseWriter

type Flusher = realHttp.Flusher

type Hijacker = realHttp.Hijacker

type CloseNotifier = realHttp.CloseNotifier

var ServerContextKey = realHttp.ServerContextKey

var LocalAddrContextKey = realHttp.LocalAddrContextKey

const TrailerPrefix = realHttp.TrailerPrefix

const DefaultMaxHeaderBytes = realHttp.DefaultMaxHeaderBytes

const TimeFormat = realHttp.TimeFormat

var ErrAbortHandler = realHttp.ErrAbortHandler

type HandlerFunc = realHttp.HandlerFunc

func Error(w ResponseWriter, error string, code int) {
	realHttp.Error(w, error, code)
}

func NotFound(w ResponseWriter, r *Request) {
	realHttp.NotFound(w, r)
}

func NotFoundHandler() Handler {
	return realHttp.NotFoundHandler()
}

func StripPrefix(prefix string, h Handler) Handler {
	return realHttp.StripPrefix(prefix, h)
}

func Redirect(w ResponseWriter, r *Request, url string, code int) {
	realHttp.Redirect(w, r, url, code)
}

func RedirectHandler(url string, code int) Handler {
	return realHttp.RedirectHandler(url, code)
}

type ServeMux = realHttp.ServeMux

func NewServeMux() *ServeMux {
	return realHttp.NewServeMux()
}

var DefaultServeMux = realHttp.DefaultServeMux

func Handle(pattern string, handler Handler) {
	handler = otelhttp.NewHandler(handler, pattern)
	realHttp.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	realHttp.HandleFunc(pattern, handler)
}

func Serve(l net.Listener, handler Handler) error {
	return realHttp.Serve(l, handler)
}

func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error {
	return realHttp.ServeTLS(l, handler, certFile, keyFile)
}

type Server = realHttp.Server

type ConnState = realHttp.ConnState

const StateNew = realHttp.StateNew

const StateActive = realHttp.StateActive

const StateIdle = realHttp.StateIdle

const StateHijacked = realHttp.StateHijacked

const StateClosed = realHttp.StateClosed

func AllowQuerySemicolons(h Handler) Handler {
	return realHttp.AllowQuerySemicolons(h)
}

var ErrServerClosed = realHttp.ErrServerClosed

func ListenAndServe(addr string, handler Handler) error {
	return realHttp.ListenAndServe(addr, handler)
}

func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	return realHttp.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
	return realHttp.TimeoutHandler(h, dt, msg)
}

var ErrHandlerTimeout = realHttp.ErrHandlerTimeout

func MaxBytesHandler(h Handler, n int64) Handler {
	return realHttp.MaxBytesHandler(h, n)
}

var ErrMissingFile = realHttp.ErrMissingFile

type ProtocolError = realHttp.ProtocolError

var ErrNotSupported = realHttp.ErrNotSupported

var ErrUnexpectedTrailer = realHttp.ErrUnexpectedTrailer

var ErrMissingBoundary = realHttp.ErrMissingBoundary

var ErrNotMultipart = realHttp.ErrNotMultipart

var ErrHeaderTooLong = realHttp.ErrHeaderTooLong

var ErrShortBody = realHttp.ErrShortBody

var ErrMissingContentLength = realHttp.ErrMissingContentLength

type Request = realHttp.Request

var ErrNoCookie = realHttp.ErrNoCookie

func ParseHTTPVersion(vers string) (major, minor int, ok bool) {
	return realHttp.ParseHTTPVersion(vers)
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	return realHttp.NewRequest(method, url, body)
}

func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return realHttp.NewRequestWithContext(ctx, method, url, body)
}

func ReadRequest(b *bufio.Reader) (*Request, error) {
	return realHttp.ReadRequest(b)
}

func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	return realHttp.MaxBytesReader(w, r, n)
}

type MaxBytesError = realHttp.MaxBytesError

type Response = realHttp.Response

var ErrNoLocation = realHttp.ErrNoLocation

func ReadResponse(r *bufio.Reader, req *Request) (*Response, error) {
	return realHttp.ReadResponse(r, req)
}

var DefaultTransport = realHttp.DefaultTransport

const DefaultMaxIdleConnsPerHost = realHttp.DefaultMaxIdleConnsPerHost

type Transport = realHttp.Transport

func ProxyFromEnvironment(req *Request) (*url.URL, error) {
	return realHttp.ProxyFromEnvironment(req)
}

func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error) {
	return realHttp.ProxyURL(fixedURL)
}

var ErrSkipAltProtocol = realHttp.ErrSkipAltProtocol

type Header = realHttp.Header

func ParseTime(text string) (t time.Time, err error) {
	return realHttp.ParseTime(text)
}

func CanonicalHeaderKey(s string) string {
	return realHttp.CanonicalHeaderKey(s)
}

type CookieJar = realHttp.CookieJar

func DetectContentType(data []byte) string {
	return realHttp.DetectContentType(data)
}

const MethodGet = realHttp.MethodGet

const MethodHead = realHttp.MethodHead

const MethodPost = realHttp.MethodPost

const MethodPut = realHttp.MethodPut

const MethodPatch = realHttp.MethodPatch

const MethodDelete = realHttp.MethodDelete

const MethodConnect = realHttp.MethodConnect

const MethodOptions = realHttp.MethodOptions

const MethodTrace = realHttp.MethodTrace

var DefaultClient = &Client{Transport: otelhttp.NewTransport(realHttp.DefaultTransport)}

type RoundTripper = realHttp.RoundTripper

func Get(url string) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

var ErrUseLastResponse = realHttp.ErrUseLastResponse

func Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "POST", url, body)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func PostForm(url string, data url.Values) (resp *Response, err error) {
	return realHttp.PostForm(url, data)
}

func Head(url string) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func NewFileTransport(fs FileSystem) RoundTripper {
	return realHttp.NewFileTransport(fs)
}

const StatusContinue = realHttp.StatusContinue

const StatusSwitchingProtocols = realHttp.StatusSwitchingProtocols

const StatusProcessing = realHttp.StatusProcessing

const StatusEarlyHints = realHttp.StatusEarlyHints

const StatusOK = realHttp.StatusOK

const StatusCreated = realHttp.StatusCreated

const StatusAccepted = realHttp.StatusAccepted

const StatusNonAuthoritativeInfo = realHttp.StatusNonAuthoritativeInfo

const StatusNoContent = realHttp.StatusNoContent

const StatusResetContent = realHttp.StatusResetContent

const StatusPartialContent = realHttp.StatusPartialContent

const StatusMultiStatus = realHttp.StatusMultiStatus

const StatusAlreadyReported = realHttp.StatusAlreadyReported

const StatusIMUsed = realHttp.StatusIMUsed

const StatusMultipleChoices = realHttp.StatusMultipleChoices

const StatusMovedPermanently = realHttp.StatusMovedPermanently

const StatusFound = realHttp.StatusFound

const StatusSeeOther = realHttp.StatusSeeOther

const StatusNotModified = realHttp.StatusNotModified

const StatusUseProxy = realHttp.StatusUseProxy

const StatusTemporaryRedirect = realHttp.StatusTemporaryRedirect

const StatusPermanentRedirect = realHttp.StatusPermanentRedirect

const StatusBadRequest = realHttp.StatusBadRequest

const StatusUnauthorized = realHttp.StatusUnauthorized

const StatusPaymentRequired = realHttp.StatusPaymentRequired

const StatusForbidden = realHttp.StatusForbidden

const StatusNotFound = realHttp.StatusNotFound

const StatusMethodNotAllowed = realHttp.StatusMethodNotAllowed

const StatusNotAcceptable = realHttp.StatusNotAcceptable

const StatusProxyAuthRequired = realHttp.StatusProxyAuthRequired

const StatusRequestTimeout = realHttp.StatusRequestTimeout

const StatusConflict = realHttp.StatusConflict

const StatusGone = realHttp.StatusGone

const StatusLengthRequired = realHttp.StatusLengthRequired

const StatusPreconditionFailed = realHttp.StatusPreconditionFailed

const StatusRequestEntityTooLarge = realHttp.StatusRequestEntityTooLarge

const StatusRequestURITooLong = realHttp.StatusRequestURITooLong

const StatusUnsupportedMediaType = realHttp.StatusUnsupportedMediaType

const StatusRequestedRangeNotSatisfiable = realHttp.StatusRequestedRangeNotSatisfiable

const StatusExpectationFailed = realHttp.StatusExpectationFailed

const StatusTeapot = realHttp.StatusTeapot

const StatusMisdirectedRequest = realHttp.StatusMisdirectedRequest

const StatusUnprocessableEntity = realHttp.StatusUnprocessableEntity

const StatusLocked = realHttp.StatusLocked

const StatusFailedDependency = realHttp.StatusFailedDependency

const StatusTooEarly = realHttp.StatusTooEarly

const StatusUpgradeRequired = realHttp.StatusUpgradeRequired

const StatusPreconditionRequired = realHttp.StatusPreconditionRequired

const StatusTooManyRequests = realHttp.StatusTooManyRequests

const StatusRequestHeaderFieldsTooLarge = realHttp.StatusRequestHeaderFieldsTooLarge

const StatusUnavailableForLegalReasons = realHttp.StatusUnavailableForLegalReasons

const StatusInternalServerError = realHttp.StatusInternalServerError

const StatusNotImplemented = realHttp.StatusNotImplemented

const StatusBadGateway = realHttp.StatusBadGateway

const StatusServiceUnavailable = realHttp.StatusServiceUnavailable

const StatusGatewayTimeout = realHttp.StatusGatewayTimeout

const StatusHTTPVersionNotSupported = realHttp.StatusHTTPVersionNotSupported

const StatusVariantAlsoNegotiates = realHttp.StatusVariantAlsoNegotiates

const StatusInsufficientStorage = realHttp.StatusInsufficientStorage

const StatusLoopDetected = realHttp.StatusLoopDetected

const StatusNotExtended = realHttp.StatusNotExtended

const StatusNetworkAuthenticationRequired = realHttp.StatusNetworkAuthenticationRequired

func StatusText(code int) string {
	return realHttp.StatusText(code)
}

var NoBody = realHttp.NoBody

type PushOptions = realHttp.PushOptions

type Pusher = realHttp.Pusher

var ErrLineTooLong = realHttp.ErrLineTooLong

var ErrBodyReadAfterClose = realHttp.ErrBodyReadAfterClose

type Cookie = realHttp.Cookie

type SameSite = realHttp.SameSite

const SameSiteDefaultMode = realHttp.SameSiteDefaultMode

const SameSiteLaxMode = realHttp.SameSiteLaxMode

const SameSiteStrictMode = realHttp.SameSiteStrictMode

const SameSiteNoneMode = realHttp.SameSiteNoneMode

func SetCookie(w ResponseWriter, cookie *Cookie) {
	realHttp.SetCookie(w, cookie)
}

type Dir = realHttp.Dir

type FileSystem = realHttp.FileSystem

type File = realHttp.File

func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker) {
	realHttp.ServeContent(w, req, name, modtime, content)
}

func ServeFile(w ResponseWriter, r *Request, name string) {
	realHttp.ServeFile(w, r, name)
}

func FS(fsys fs.FS) FileSystem {
	return realHttp.FS(fsys)
}

func FileServer(root FileSystem) Handler {
	return realHttp.FileServer(root)
}
