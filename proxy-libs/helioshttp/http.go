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

// Errors used by the HTTP server.
var (
	ErrNotSupported         = realHttp.ErrNotSupported
	ErrUnexpectedTrailer    = realHttp.ErrUnexpectedTrailer
	ErrMissingBoundary      = realHttp.ErrMissingBoundary
	ErrNotMultipart         = realHttp.ErrNotMultipart
	ErrHeaderTooLong        = realHttp.ErrHeaderTooLong
	ErrShortBody            = realHttp.ErrShortBody
	ErrMissingContentLength = realHttp.ErrMissingContentLength
)

var (
	ErrBodyNotAllowed  = realHttp.ErrBodyNotAllowed
	ErrHijacked        = realHttp.ErrHijacked
	ErrContentLength   = realHttp.ErrContentLength
	ErrWriteAfterFlush = realHttp.ErrWriteAfterFlush
)

var (
	ServerContextKey    = realHttp.ServerContextKey
	LocalAddrContextKey = realHttp.LocalAddrContextKey
)

type Header = realHttp.Header
type Response = realHttp.Response
type Transport = realHttp.Transport
type Handler = realHttp.Handler
type ResponseWriter = realHttp.ResponseWriter
type Flusher = realHttp.Flusher
type Hijacker = realHttp.Hijacker
type CloseNotifier = realHttp.CloseNotifier
type Request = realHttp.Request

const TrailerPrefix = realHttp.TrailerPrefix
const DefaultMaxHeaderBytes = realHttp.DefaultMaxHeaderBytes
const TimeFormat = realHttp.TimeFormat
const DefaultMaxIdleConnsPerHost = realHttp.DefaultMaxIdleConnsPerHost

var ErrAbortHandler = realHttp.ErrAbortHandler
var ErrBodyReadAfterClose = realHttp.ErrBodyReadAfterClose
var ErrHandlerTimeout = realHttp.ErrHandlerTimeout
var ErrLineTooLong = realHttp.ErrLineTooLong
var ErrMissingFile = realHttp.ErrMissingFile
var ErrNoCookie = realHttp.ErrNoCookie
var ErrNoLocation = realHttp.ErrNoLocation
var ErrServerClosed = realHttp.ErrServerClosed
var ErrSkipAltProtocol = realHttp.ErrSkipAltProtocol
var ErrUseLastResponse = realHttp.ErrUseLastResponse
var NoBody = realHttp.NoBody
var DefaultTransport = realHttp.DefaultTransport

type HandlerFunc = realHttp.HandlerFunc

func FileServer(root FileSystem) Handler {
	return realHttp.FileServer(root)
}

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
type File = realHttp.File

func NewServeMux() *ServeMux {
	return realHttp.NewServeMux()
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	return realHttp.NewRequest(method, url, body)
}

func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return realHttp.NewRequestWithContext(ctx, method, url, body)
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

func FS(fsys fs.FS) FileSystem {
	return realHttp.FS(fsys)
}

var DefaultClient = &Client{}

func Get(url string) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func Head(url string) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

func Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, "POST", url, body)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}

	return DefaultClient.Do(req)
}

type Server = realHttp.Server

type ConnState = realHttp.ConnState

const (
	StateNew      ConnState = realHttp.StateNew
	StateActive             = realHttp.StateActive
	StateIdle               = realHttp.StateIdle
	StateHijacked           = realHttp.StateHijacked
	StateClosed             = realHttp.StateClosed
)

func AllowQuerySemicolons(h Handler) Handler {
	return realHttp.AllowQuerySemicolons(h)
}

func ListenAndServe(addr string, handler Handler) error {
	return realHttp.ListenAndServe(addr, handler)
}

func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	return realHttp.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
	return realHttp.TimeoutHandler(h, dt, msg)
}

const (
	StatusContinue                      = realHttp.StatusContinue
	StatusSwitchingProtocols            = realHttp.StatusSwitchingProtocols
	StatusProcessing                    = realHttp.StatusProcessing
	StatusEarlyHints                    = realHttp.StatusEarlyHints
	StatusOK                            = realHttp.StatusOK
	StatusCreated                       = realHttp.StatusCreated
	StatusAccepted                      = realHttp.StatusAccepted
	StatusNonAuthoritativeInfo          = realHttp.StatusNonAuthoritativeInfo
	StatusNoContent                     = realHttp.StatusNoContent
	StatusResetContent                  = realHttp.StatusResetContent
	StatusPartialContent                = realHttp.StatusPartialContent
	StatusMultiStatus                   = realHttp.StatusMultiStatus
	StatusAlreadyReported               = realHttp.StatusAlreadyReported
	StatusIMUsed                        = realHttp.StatusIMUsed
	StatusMultipleChoices               = realHttp.StatusMultipleChoices
	StatusMethodNotAllowed              = realHttp.StatusMethodNotAllowed
	StatusMovedPermanently              = realHttp.StatusMovedPermanently
	StatusFound                         = realHttp.StatusFound
	StatusNotAcceptable                 = realHttp.StatusNotAcceptable
	StatusSeeOther                      = realHttp.StatusSeeOther
	StatusNotModified                   = realHttp.StatusNotModified
	StatusUseProxy                      = realHttp.StatusUseProxy
	StatusTemporaryRedirect             = realHttp.StatusTemporaryRedirect
	StatusPermanentRedirect             = realHttp.StatusPermanentRedirect
	StatusBadRequest                    = realHttp.StatusBadRequest
	StatusUnauthorized                  = realHttp.StatusUnauthorized
	StatusPaymentRequired               = realHttp.StatusPaymentRequired
	StatusForbidden                     = realHttp.StatusForbidden
	StatusNotFound                      = realHttp.StatusNotFound
	StatusProxyAuthRequired             = realHttp.StatusProxyAuthRequired
	StatusRequestTimeout                = realHttp.StatusRequestTimeout
	StatusConflict                      = realHttp.StatusConflict
	StatusGone                          = realHttp.StatusGone
	StatusLengthRequired                = realHttp.StatusLengthRequired
	StatusPreconditionFailed            = realHttp.StatusPreconditionFailed
	StatusRequestEntityTooLarge         = realHttp.StatusRequestEntityTooLarge
	StatusRequestURITooLong             = realHttp.StatusRequestURITooLong
	StatusUnsupportedMediaType          = realHttp.StatusUnsupportedMediaType
	StatusRequestedRangeNotSatisfiable  = realHttp.StatusRequestedRangeNotSatisfiable
	StatusExpectationFailed             = realHttp.StatusExpectationFailed
	StatusTeapot                        = realHttp.StatusTeapot
	StatusMisdirectedRequest            = realHttp.StatusMisdirectedRequest
	StatusUnprocessableEntity           = realHttp.StatusUnprocessableEntity
	StatusLocked                        = realHttp.StatusLocked
	StatusFailedDependency              = realHttp.StatusFailedDependency
	StatusTooEarly                      = realHttp.StatusTooEarly
	StatusUpgradeRequired               = realHttp.StatusUpgradeRequired
	StatusPreconditionRequired          = realHttp.StatusPreconditionRequired
	StatusTooManyRequests               = realHttp.StatusTooManyRequests
	StatusRequestHeaderFieldsTooLarge   = realHttp.StatusRequestHeaderFieldsTooLarge
	StatusUnavailableForLegalReasons    = realHttp.StatusUnavailableForLegalReasons
	StatusInternalServerError           = realHttp.StatusInternalServerError
	StatusNotImplemented                = realHttp.StatusNotImplemented
	StatusBadGateway                    = realHttp.StatusBadGateway
	StatusServiceUnavailable            = realHttp.StatusServiceUnavailable
	StatusGatewayTimeout                = realHttp.StatusGatewayTimeout
	StatusHTTPVersionNotSupported       = realHttp.StatusHTTPVersionNotSupported
	StatusVariantAlsoNegotiates         = realHttp.StatusVariantAlsoNegotiates
	StatusInsufficientStorage           = realHttp.StatusInsufficientStorage
	StatusLoopDetected                  = realHttp.StatusLoopDetected
	StatusNotExtended                   = realHttp.StatusNotExtended
	StatusNetworkAuthenticationRequired = realHttp.StatusNetworkAuthenticationRequired
)

type ProtocolError = realHttp.ProtocolError
type SameSite = realHttp.SameSite

const (
	SameSiteDefaultMode SameSite = realHttp.SameSiteDefaultMode
	SameSiteLaxMode              = realHttp.SameSiteLaxMode
	SameSiteStrictMode           = realHttp.SameSiteStrictMode
	SameSiteNoneMode             = realHttp.SameSiteNoneMode
)

func MaxBytesHandler(h Handler, n int64) Handler {
	return realHttp.MaxBytesHandler(h, n)
}

func ServeFile(w ResponseWriter, r *Request, name string) {
	realHttp.ServeFile(w, r, name)
}

func ReadResponse(r *bufio.Reader, req *Request) (*Response, error) {
	return realHttp.ReadResponse(r, req)
}

func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error) {
	return realHttp.ProxyURL(fixedURL)
}

func ProxyFromEnvironment(req *Request) (*url.URL, error) {
	return realHttp.ProxyFromEnvironment(req)
}

func PostForm(url string, data url.Values) (resp *Response, err error) {
	return realHttp.PostForm(url, data)
}

func NewFileTransport(fs FileSystem) RoundTripper {
	return realHttp.NewFileTransport(fs)
}

func ParseHTTPVersion(vers string) (major, minor int, ok bool) {
	return realHttp.ParseHTTPVersion(vers)
}

func ParseTime(text string) (t time.Time, err error) {
	return realHttp.ParseTime(text)
}

const (
	MethodGet     = realHttp.MethodGet
	MethodHead    = realHttp.MethodHead
	MethodPost    = realHttp.MethodPost
	MethodPut     = realHttp.MethodPut
	MethodPatch   = realHttp.MethodPatch
	MethodDelete  = realHttp.MethodDelete
	MethodConnect = realHttp.MethodConnect
	MethodOptions = realHttp.MethodOptions
	MethodTrace   = realHttp.MethodTrace
)

type PushOptions = realHttp.PushOptions
type Pusher = realHttp.Pusher

func ReadRequest(b *bufio.Reader) (*Request, error) {
	return realHttp.ReadRequest(b)
}

func StatusText(code int) string {
	return realHttp.StatusText(code)
}

func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	return realHttp.MaxBytesReader(w, r, n)
}

func DetectContentType(data []byte) string {
	return realHttp.DetectContentType(data)
}

func CanonicalHeaderKey(s string) string {
	return realHttp.CanonicalHeaderKey(s)
}

type FileSystem = realHttp.FileSystem
type Dir = realHttp.Dir
type Cookie = realHttp.Cookie

func SetCookie(w ResponseWriter, cookie *Cookie) {
	realHttp.SetCookie(w, cookie)
}

func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker) {
	realHttp.ServeContent(w, req, name, modtime, content)
}

type RoundTripper = realHttp.RoundTripper
type CookieJar = realHttp.CookieJar
