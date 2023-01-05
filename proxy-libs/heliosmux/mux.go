package heliosmux

import (
	originalMux "github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var InstrumentedSymbols = [...]string{"NewRouter"}

func NewRouter() *originalMux.Router {
	router := originalMux.NewRouter()
	router.Use(otelmux.Middleware("opentelemetry-middleware"))
	return router
}
