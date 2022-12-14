package heliosbeego

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	origin_beego "github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
)

type FlashData = origin_beego.FlashData

func NewFlash() *FlashData {
	return origin_beego.NewFlash()
}

func ReadFromRequest(c *Controller) *FlashData {
	return origin_beego.ReadFromRequest(c)
}

const LevelEmergency = origin_beego.LevelEmergency

const LevelAlert = origin_beego.LevelAlert

const LevelCritical = origin_beego.LevelCritical

const LevelError = origin_beego.LevelError

const LevelWarning = origin_beego.LevelWarning

const LevelNotice = origin_beego.LevelNotice

const LevelInformational = origin_beego.LevelInformational

const LevelDebug = origin_beego.LevelDebug

var BeeLogger = origin_beego.BeeLogger

func SetLevel(l int) {
	origin_beego.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	origin_beego.SetLogFuncCall(b)
}

func SetLogger(adaptername string, config string) error {
	return origin_beego.SetLogger(adaptername, config)
}

func Emergency(v ...interface{}) {
	origin_beego.Emergency(v)
}

func Alert(v ...interface{}) {
	origin_beego.Alert(v)
}

func Critical(v ...interface{}) {
	origin_beego.Critical(v)
}

func Error(v ...interface{}) {
	origin_beego.Error(v)
}

func Warning(v ...interface{}) {
	origin_beego.Warning(v)
}

func Warn(v ...interface{}) {
	origin_beego.Warn(v)
}

func Notice(v ...interface{}) {
	origin_beego.Notice(v)
}

func Informational(v ...interface{}) {
	origin_beego.Informational(v)
}

func Info(v ...interface{}) {
	origin_beego.Info(v)
}

func Debug(v ...interface{}) {
	origin_beego.Debug(v)
}

func Trace(v ...interface{}) {
	origin_beego.Trace(v)
}

const VERSION = origin_beego.VERSION

const DEV = origin_beego.DEV

const PROD = origin_beego.PROD

type M = origin_beego.M

// TODO: readd and solve
// func AddAPPStartHook(hf ...[]func() error) {
// origin_beego.AddAPPStartHook(hf...)
// }

func Run(params ...string) {
	origin_beego.Run(params...)
}

func RunWithMiddleWares(addr string, mws ...MiddleWare) {
	otel_middleware := otelbeego.NewOTelBeegoMiddleWare("opentelemetry-middleware")

	new_mws := make([]MiddleWare, len(mws)+1)
	new_mws = append(new_mws, otel_middleware)
	new_mws = append(new_mws, mws...)

	origin_beego.RunWithMiddleWares(addr, new_mws...)
}

func TestBeegoInit(ap string) {
	origin_beego.TestBeegoInit(ap)
}

func InitBeegoBeforeTest(appConfigPath string) {
	origin_beego.InitBeegoBeforeTest(appConfigPath)
}

var ErrAbort = origin_beego.ErrAbort

var GlobalControllerRouter = origin_beego.GlobalControllerRouter

type ControllerFilter = origin_beego.ControllerFilter

type ControllerFilterComments = origin_beego.ControllerFilterComments

type ControllerImportComments = origin_beego.ControllerImportComments

type ControllerComments = origin_beego.ControllerComments

type ControllerCommentsSlice = origin_beego.ControllerCommentsSlice

type Controller = origin_beego.Controller

type ControllerInterface = origin_beego.ControllerInterface

var ErrorMaps = origin_beego.ErrorMaps

func ErrorHandler(code string, h http.HandlerFunc) *App {
	return origin_beego.ErrorHandler(code, h)
}

func ErrorController(c ControllerInterface) *App {
	return origin_beego.ErrorController(c)
}

func Exception(errCode uint64, ctx *context.Context) {
	origin_beego.Exception(errCode, ctx)
}

type PolicyFunc = origin_beego.PolicyFunc

func Policy(pattern, method string, policy ...PolicyFunc) {
	origin_beego.Policy(pattern, method, policy...)
}

const BeforeStatic = origin_beego.BeforeStatic

const BeforeRouter = origin_beego.BeforeRouter

const BeforeExec = origin_beego.BeforeExec

const AfterExec = origin_beego.AfterExec

const FinishRouter = origin_beego.FinishRouter

var HTTPMETHOD = origin_beego.HTTPMETHOD

var DefaultAccessLogFilter = origin_beego.DefaultAccessLogFilter

type FilterHandler = origin_beego.FilterHandler

func ExceptMethodAppend(action string) {
	origin_beego.ExceptMethodAppend(action)
}

type ControllerInfo = origin_beego.ControllerInfo

type ControllerRegister = origin_beego.ControllerRegister

func NewControllerRegister() *ControllerRegister {
	return origin_beego.NewControllerRegister()
}

func LogAccess(ctx *context.Context, startTime *time.Time, statusCode int) {
	origin_beego.LogAccess(ctx, startTime, statusCode)
}

type LinkNamespace = origin_beego.LinkNamespace

type Namespace = origin_beego.Namespace
type NamespaceCond = func(*context.Context) bool

func NewNamespace(prefix string, params ...LinkNamespace) *Namespace {
	return origin_beego.NewNamespace(prefix, params...)
}

func AddNamespace(nl ...*Namespace) {
	origin_beego.AddNamespace(nl...)
}

func NSCond(cond NamespaceCond) LinkNamespace {
	return origin_beego.NSCond(cond)
}

func NSBefore(filterList ...FilterFunc) LinkNamespace {
	return origin_beego.NSBefore(filterList...)
}

func NSAfter(filterList ...FilterFunc) LinkNamespace {
	return origin_beego.NSAfter(filterList...)
}

func NSInclude(cList ...ControllerInterface) LinkNamespace {
	return origin_beego.NSInclude(cList...)
}

func NSRouter(rootpath string, c ControllerInterface, mappingMethods ...string) LinkNamespace {
	return origin_beego.NSRouter(rootpath, c, mappingMethods...)
}

func NSGet(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSGet(rootpath, f)
}

func NSPost(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSPost(rootpath, f)
}

func NSHead(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSHead(rootpath, f)
}

func NSPut(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSPut(rootpath, f)
}

func NSDelete(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSDelete(rootpath, f)
}

func NSAny(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSAny(rootpath, f)
}

func NSOptions(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSOptions(rootpath, f)
}

func NSPatch(rootpath string, f FilterFunc) LinkNamespace {
	return origin_beego.NSPatch(rootpath, f)
}

func NSAutoRouter(c ControllerInterface) LinkNamespace {
	return origin_beego.NSAutoRouter(c)
}

func NSAutoPrefix(prefix string, c ControllerInterface) LinkNamespace {
	return origin_beego.NSAutoPrefix(prefix, c)
}

func NSNamespace(prefix string, params ...LinkNamespace) LinkNamespace {
	return origin_beego.NSNamespace(prefix, params...)
}

func NSHandler(rootpath string, h http.Handler) LinkNamespace {
	return origin_beego.NSHandler(rootpath, h)
}

func Substr(s string, start int, length int) string {
	return origin_beego.Substr(s, start, length)
}

func HTML2str(html string) string {
	return origin_beego.HTML2str(html)
}

func DateFormat(t time.Time, layout string) (datestring string) {
	return origin_beego.DateFormat(t, layout)
}

func DateParse(dateString string, format string) (time.Time, error) {
	return origin_beego.DateParse(dateString, format)
}

func Date(t time.Time, format string) string {
	return origin_beego.Date(t, format)
}

func Compare(a, b interface{}) (equal bool) {
	return origin_beego.Compare(a, b)
}

func CompareNot(a, b interface{}) (equal bool) {
	return origin_beego.CompareNot(a, b)
}

func NotNil(a interface{}) (isNil bool) {
	return origin_beego.NotNil(a)
}

func GetConfig(returnType, key string, defaultVal interface{}) (value interface{}, err error) {
	return origin_beego.GetConfig(returnType, key, defaultVal)
}

func Str2html(raw string) template.HTML {
	return origin_beego.Str2html(raw)
}

func Htmlquote(text string) string {
	return origin_beego.Htmlquote(text)
}

func Htmlunquote(text string) string {
	return origin_beego.Htmlunquote(text)
}

func URLFor(endpoint string, values ...interface{}) string {
	return origin_beego.URLFor(endpoint, values)
}

func AssetsJs(text string) template.HTML {
	return origin_beego.AssetsJs(text)
}

func AssetsCSS(text string) template.HTML {
	return origin_beego.AssetsCSS(text)
}

func ParseForm(form url.Values, obj interface{}) error {
	return origin_beego.ParseForm(form, obj)
}

func RenderForm(obj interface{}) template.HTML {
	return origin_beego.RenderForm(obj)
}

func MapGet(arg1 interface{}, arg2 ...interface{}) (interface{}, error) {
	return origin_beego.MapGet(arg1, arg2)
}

var FilterMonitorFunc = origin_beego.FilterMonitorFunc

func PrintTree() M {
	return origin_beego.PrintTree()
}

type Config = origin_beego.Config

type Listen = origin_beego.Listen

type WebConfig = origin_beego.WebConfig

type SessionConfig = origin_beego.SessionConfig

type LogConfig = origin_beego.LogConfig

var BConfig = origin_beego.BConfig

var AppConfig = origin_beego.AppConfig

var AppPath = origin_beego.AppPath

var GlobalSessions = origin_beego.GlobalSessions

var WorkPath = origin_beego.WorkPath

func LoadAppConfig(adapterName string, configPath string) error {
	return origin_beego.LoadAppConfig(adapterName, configPath)
}

func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return origin_beego.ExecuteTemplate(wr, name, data)
}

func ExecuteViewPathTemplate(wr io.Writer, name string, viewPath string, data interface{}) error {
	return origin_beego.ExecuteViewPathTemplate(wr, name, viewPath, data)
}

func AddFuncMap(key string, fn interface{}) error {
	return origin_beego.AddFuncMap(key, fn)
}

func HasTemplateExt(paths string) bool {
	return origin_beego.HasTemplateExt(paths)
}

func AddTemplateExt(ext string) {
	origin_beego.AddTemplateExt(ext)
}

func AddViewPath(viewPath string) error {
	return origin_beego.AddViewPath(viewPath)
}

func BuildTemplate(dir string, files ...string) error {
	return origin_beego.BuildTemplate(dir, files...)
}

func SetTemplateFSFunc(fnt func() http.FileSystem) {
	origin_beego.SetTemplateFSFunc(fnt)
}

func SetViewsPath(path string) *App {
	return origin_beego.SetViewsPath(path)
}

func SetStaticPath(url string, path string) *App {
	return origin_beego.SetStaticPath(url, path)
}

func DelStaticPath(url string) *App {
	return origin_beego.DelStaticPath(url)
}

func AddTemplateEngine(extension string, fn func(root, path string, funcs template.FuncMap) (*template.Template, error)) *App {
	return origin_beego.AddTemplateEngine(extension, fn)
}

var BeeApp = origin_beego.BeeApp

type App = origin_beego.App

func NewApp() *App {
	return origin_beego.NewApp()
}

type MiddleWare = origin_beego.MiddleWare

func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	return origin_beego.Router(rootpath, c, mappingMethods...)
}

func UnregisterFixedRoute(fixedRoute string, method string) *App {
	return origin_beego.UnregisterFixedRoute(fixedRoute, method)
}

func Include(cList ...ControllerInterface) *App {
	return origin_beego.Include(cList...)
}

func RESTRouter(rootpath string, c ControllerInterface) *App {
	return origin_beego.RESTRouter(rootpath, c)
}

func AutoRouter(c ControllerInterface) *App {
	return origin_beego.AutoRouter(c)
}

func AutoPrefix(prefix string, c ControllerInterface) *App {
	return origin_beego.AutoPrefix(prefix, c)
}

func Get(rootpath string, f FilterFunc) *App {
	return origin_beego.Get(rootpath, f)
}

func Post(rootpath string, f FilterFunc) *App {
	return origin_beego.Post(rootpath, f)
}

func Delete(rootpath string, f FilterFunc) *App {
	return origin_beego.Delete(rootpath, f)
}

func Put(rootpath string, f FilterFunc) *App {
	return origin_beego.Put(rootpath, f)
}

func Head(rootpath string, f FilterFunc) *App {
	return origin_beego.Head(rootpath, f)
}

func Options(rootpath string, f FilterFunc) *App {
	return origin_beego.Options(rootpath, f)
}

func Patch(rootpath string, f FilterFunc) *App {
	return origin_beego.Patch(rootpath, f)
}

func Any(rootpath string, f FilterFunc) *App {
	return origin_beego.Any(rootpath, f)
}

func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	return origin_beego.Handler(rootpath, h, options)
}

func InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) *App {
	return origin_beego.InsertFilter(pattern, pos, filter, params...)
}

type FilterFunc = origin_beego.FilterFunc

type FilterRouter = origin_beego.FilterRouter

type Tree = origin_beego.Tree

func NewTree() *Tree {
	return origin_beego.NewTree()
}

var BuildVersion = origin_beego.BuildVersion

var BuildGitRevision = origin_beego.BuildGitRevision

var BuildStatus = origin_beego.BuildStatus

var BuildTag = origin_beego.BuildTag

var BuildTime = origin_beego.BuildTime

var GoVersion = origin_beego.GoVersion

var GitBranch = origin_beego.GitBranch

type FileSystem = origin_beego.FileSystem

func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	return origin_beego.Walk(fs, root, walkFn)
}
