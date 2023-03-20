package heliosmacaron

import (
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron"
	"gopkg.in/ini.v1"
	origin_macaron "gopkg.in/macaron.v1"
)

func PathUnescape(s string) (string, error) {
	return origin_macaron.PathUnescape(s)
}

type StaticOptions = origin_macaron.StaticOptions

func GenerateETag(fileSize, fileName, modTime string) string {
	return origin_macaron.GenerateETag(fileSize, fileName, modTime)
}

func Static(directory string, staticOpt ...StaticOptions) Handler {
	return origin_macaron.Static(directory, staticOpt...)
}

func Statics(opt StaticOptions, dirs ...string) Handler {
	return origin_macaron.Statics(opt, dirs...)
}

func Recovery() Handler {
	return origin_macaron.Recovery()
}

func NewRouteMap() interface{} {
	return origin_macaron.NewRouteMap()
}

type Router = origin_macaron.Router

func NewRouter() *Router {
	return origin_macaron.NewRouter()
}

type Params = origin_macaron.Params

type Handle = origin_macaron.Handle

type Route = origin_macaron.Route

type ComboRouter = origin_macaron.ComboRouter

type TemplateFile = origin_macaron.TemplateFile

type TemplateFileSystem = origin_macaron.TemplateFileSystem

type Delims = origin_macaron.Delims

type RenderOptions = origin_macaron.RenderOptions

type HTMLOptions = origin_macaron.HTMLOptions

type Render = origin_macaron.Render

type TplFile = origin_macaron.TplFile

func NewTplFile(name string, data []byte, ext string) *TplFile {
	return origin_macaron.NewTplFile(name, data, ext)
}

type TplFileSystem = origin_macaron.TplFileSystem

func NewTemplateFileSystem(opt RenderOptions, omitData bool) TplFileSystem {
	return origin_macaron.NewTemplateFileSystem(opt, omitData)
}

func PrepareCharset(charset string) string {
	return origin_macaron.PrepareCharset(charset)
}

func GetExt(s string) string {
	return origin_macaron.GetExt(s)
}

const DEFAULT_TPL_SET_NAME = origin_macaron.DEFAULT_TPL_SET_NAME

type TemplateSet = origin_macaron.TemplateSet

func NewTemplateSet() *TemplateSet {
	return origin_macaron.NewTemplateSet()
}

func ParseTplSet(tplSet string) (tplName string, tplDir string) {
	return origin_macaron.ParseTplSet(tplSet)
}

func Renderer(options ...RenderOptions) Handler {
	return origin_macaron.Renderer(options...)
}

func Renderers(options RenderOptions, tplSets ...string) Handler {
	return origin_macaron.Renderers(options, tplSets...)
}

type TplRender = origin_macaron.TplRender

type DummyRender = origin_macaron.DummyRender

type ResponseWriter = origin_macaron.ResponseWriter

type BeforeFunc = origin_macaron.BeforeFunc

func NewResponseWriter(method string, rw http.ResponseWriter) ResponseWriter {
	return origin_macaron.NewResponseWriter(method, rw)
}

type ReturnHandler = origin_macaron.ReturnHandler

type Locale = origin_macaron.Locale

type RequestBody = origin_macaron.RequestBody

type Request = origin_macaron.Request

type ContextInvoker = origin_macaron.ContextInvoker

type Context = origin_macaron.Context

var MaxMemory = origin_macaron.MaxMemory

type Leaf = origin_macaron.Leaf

func NewLeaf(parent *Tree, pattern string, handle Handle) *Leaf {
	return origin_macaron.NewLeaf(parent, pattern, handle)
}

type Tree = origin_macaron.Tree

func NewSubtree(parent *Tree, pattern string) *Tree {
	return origin_macaron.NewSubtree(parent, pattern)
}

func NewTree() *Tree {
	return origin_macaron.NewTree()
}

func MatchTest(pattern, url string) bool {
	return origin_macaron.MatchTest(pattern, url)
}

func Version() string {
	return origin_macaron.Version()
}

type Handler = origin_macaron.Handler

type Macaron = origin_macaron.Macaron

func addOtelMiddleware(mac *Macaron) {
	if os.Getenv("HS_DISABLED") != "true" {
		mac.Use(otelmacaron.Middleware("opentelemetry-middleware"))
	}
}

func NewWithLogger(out io.Writer) *Macaron {
	mac := origin_macaron.NewWithLogger(out)
	addOtelMiddleware(mac)
	return mac
}

func New() *Macaron {
	mac := origin_macaron.New()
	addOtelMiddleware(mac)
	return mac
}

func Classic() *Macaron {
	mac := origin_macaron.Classic()
	addOtelMiddleware(mac)
	return mac
}

type BeforeHandler = origin_macaron.BeforeHandler

func GetDefaultListenInfo() (string, int) {
	return origin_macaron.GetDefaultListenInfo()
}

const DEV = origin_macaron.DEV

const PROD = origin_macaron.PROD

const TEST = origin_macaron.TEST

var Env = origin_macaron.Env

var Root = origin_macaron.Root

var FlashNow = origin_macaron.FlashNow

func SetConfig(source interface{}, others ...interface{}) (_ *ini.File, err error) {
	return origin_macaron.SetConfig(source, others)
}

func Config() *ini.File {
	return origin_macaron.Config()
}

var ColorLog = origin_macaron.ColorLog

var LogTimeFormat = origin_macaron.LogTimeFormat

type LoggerInvoker = origin_macaron.LoggerInvoker

func Logger() Handler {
	return origin_macaron.Logger()
}
