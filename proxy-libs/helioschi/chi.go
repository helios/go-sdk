package helioschi

import (
	origin_chi "github.com/go-chi/chi/v5"
	"github.com/helios/otelchi"
)

var InstrumentedSymbols = [...]string{"NewRouter"}

func addOpentelemetryMiddleware(mux *origin_chi.Mux) {
	mux.Use(otelchi.Middleware("opentelemetry-middleware", otelchi.WithChiRoutes(mux)))
}

func NewRouter() *origin_chi.Mux {
	router := origin_chi.NewRouter()
	addOpentelemetryMiddleware(router)
	return router
}

func NewMux() *origin_chi.Mux {
	mux := origin_chi.NewMux()
	addOpentelemetryMiddleware(mux)
	return mux
}
