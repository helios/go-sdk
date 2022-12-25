package helioschi

import (
	"context"
	"net/http"

	origin_chi "github.com/go-chi/chi/v5"
)

func RegisterMethod(method string) {
	origin_chi.RegisterMethod(method)
}

type Route = origin_chi.Route

type WalkFunc = origin_chi.WalkFunc

func Walk(r Routes, walkFn WalkFunc) error {
	return origin_chi.Walk(r, walkFn)
}

func Chain(middlewares ...func(http.Handler) http.Handler) Middlewares {
	return origin_chi.Chain(middlewares...)
}

type ChainHandler = origin_chi.ChainHandler

func addOpentelemetryMiddleware(mux *Mux) {
	mux.Use(Middleware("opentelemetry-middleware", WithChiRoutes(mux)))
}

func NewRouter() *Mux {
	router := origin_chi.NewRouter()
	addOpentelemetryMiddleware(router)
	return router
}

type Router = origin_chi.Router

type Routes = origin_chi.Routes

type Middlewares = origin_chi.Middlewares

func URLParam(r *http.Request, key string) string {
	return origin_chi.URLParam(r, key)
}

func URLParamFromCtx(ctx context.Context, key string) string {
	return origin_chi.URLParamFromCtx(ctx, key)
}

func RouteContext(ctx context.Context) *Context {
	return origin_chi.RouteContext(ctx)
}

func NewRouteContext() *Context {
	return origin_chi.NewRouteContext()
}

var RouteCtxKey = origin_chi.RouteCtxKey

type Context = origin_chi.Context

type RouteParams = origin_chi.RouteParams

type Mux = origin_chi.Mux

func NewMux() *Mux {
	mux := origin_chi.NewMux()
	addOpentelemetryMiddleware(mux)
	return mux
}
