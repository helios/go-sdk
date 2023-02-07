package heliosmux

import (
	"net/http"

	originalMux "github.com/gorilla/mux"
	"github.com/helios/opentelemetry-go-contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type BuildVarsFunc = originalMux.BuildVarsFunc
type MatcherFunc = originalMux.MatcherFunc
type MiddlewareFunc = originalMux.MiddlewareFunc
type Route = originalMux.Route
type RouteMatch = originalMux.RouteMatch
type Router = originalMux.Router
type WalkFunc = originalMux.WalkFunc

var ErrMethodMismatch = originalMux.ErrMethodMismatch
var ErrNotFound = originalMux.ErrNotFound
var SkipRouter = originalMux.SkipRouter

func CORSMethodMiddleware(r *Router) MiddlewareFunc {
	return originalMux.CORSMethodMiddleware(r)
}

func CurrentRoute(r *http.Request) *Route {
	return originalMux.CurrentRoute(r)
}

func NewRouter() *Router {
	router := originalMux.NewRouter()
	router.Use(otelmux.Middleware("opentelemetry-middleware"))
	return router
}

func SetURLVars(r *http.Request, val map[string]string) *http.Request {
	return originalMux.SetURLVars(r, val)
}

func Vars(r *http.Request) map[string]string {
	return originalMux.Vars(r)
}
